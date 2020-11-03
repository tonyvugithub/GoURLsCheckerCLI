package features

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/mb0/glob"
	"github.com/tonyvugithub/GoURLsCheckerCLI/lib/utils"
	"github.com/tonyvugithub/GoURLsCheckerCLI/models"
)

//CheckVersion to display the current version of the app
func CheckVersion(versionFlag *bool) {
	//Check if version flag was provided avsc
	if *versionFlag {
		//If exactly 2 arguments, print the name and version of the app
		if len(os.Args) == 2 {
			utils.PrintVersion()
			os.Exit(0)
		} else {
			fmt.Println("There should be no extra argument after -v/-version")
			os.Exit(1)
		}
	}
}

//CheckWithFileFlag to run link checking with -f flag
func CheckWithFileFlag(globFlag *bool, args []string, channel chan models.LinkStatus, wg *sync.WaitGroup, userAgent *string, summary *models.CheckSummary) {
	if *globFlag {
		fmt.Println("-f flag cannot go with -g flag, maybe you mean -d ?")
	} else {
		for _, file := range args {
			wg.Add(1)
			go func(f string) {
				defer wg.Done()
				checkByFilepath(f, channel, *userAgent, summary)
			}(file)
		}
	}
}

//CheckWithDirectoryFlag to run link checking with -d flag
func CheckWithDirectoryFlag(globFlag *bool, args []string, channel chan models.LinkStatus, wg *sync.WaitGroup, userAgent *string, summary *models.CheckSummary) {
	if *globFlag {
		argsWithoutPattern := args[:len(args)-1]
		//Assign the glob pattern provided to a local variable
		pattern := args[len(args)-1]
		CheckWithGlobFlag(pattern, argsWithoutPattern, channel, wg, userAgent, summary)
	} else {
		for _, dirPath := range args {
			//Read all file from the directory path
			files, err := ioutil.ReadDir(dirPath)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			for _, file := range files {
				filepath := filepath.Join(dirPath, file.Name())
				wg.Add(1)
				go func(f string) {
					defer wg.Done()
					checkByFilepath(f, channel, *userAgent, summary)
				}(filepath)
			}
		}
	}
}

//CheckWithGlobFlag to run link checking with a glob pattern
func CheckWithGlobFlag(pattern string, dirList []string, channel chan models.LinkStatus, wg *sync.WaitGroup, userAgent *string, summary *models.CheckSummary) {
	//Create a globber object
	globber, _ := glob.New(glob.Default())
	for _, dirPath := range dirList {
		//Read all file from the directory path
		files, err := ioutil.ReadDir(dirPath)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		for _, file := range files {
			matched, _ := globber.Match(pattern, file.Name())
			//If matched then run the url check on that file
			if matched {
				filepath := filepath.Join(dirPath, file.Name())
				wg.Add(1)
				go func(f string) {
					defer wg.Done()
					checkByFilepath(f, channel, *userAgent, summary)
				}(filepath)
			}
		}
	}
}

//CheckWithIgnoreFlag to run link checking with -i flag
func CheckWithIgnoreFlag(ignoreList string, file string, channel chan models.LinkStatus, userAgent *string, summary *models.CheckSummary) {

	fileData := utils.ReadFromFile(file)

	if ignoreList != "" {

		regLinkIgnore := regexp.MustCompile("(?m)^.*(" + ignoreList + ").*$") // finds all urls in ignore list

		fileData = regLinkIgnore.ReplaceAllString(fileData, "") // the urls from ignorelist are taken out of urls to check

	} else {
		fmt.Println("The ignore file as no urls. Therefore no urls will be ignored")
	}
	links := utils.ParseLinks(fileData)
	//Loop to check all links
	for _, link := range links {
		go utils.CheckLink(link, channel, *userAgent)
	}

	//Receive the result from checkLink and update the link to correspondent lists
	i := 0
	for i < len(links) {
		ls := <-channel
		if ls.GetLiveStatus() == true {
			summary.RecordUpLink(ls.GetURL())
		} else {
			summary.RecordDownLink(ls.GetURL())
		}
		i++
	}
}

//CheckTelescopePosts to run link checking through the lastest 10 posts if Telescope
func CheckTelescopePosts(channel chan models.LinkStatus, userAgent *string, summary *models.CheckSummary) {
	rootURL := "http://localhost:3000"

	//Get the JSON array of 10 latest posts
	resp, err := http.Get(rootURL + "/posts")

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	//Store the body of response
	body, err := ioutil.ReadAll(resp.Body)

	//The type of object in the json array
	type Post struct {
		ID  string
		URL string
	}

	var posts []Post //Variable to store the decoded json

	//Decode the JSON
	if err := json.Unmarshal(body, &posts); err != nil {
		panic(err)
	}

	//Set up custom client to fetch from post links from above
	client := &http.Client{}

	for _, post := range posts {
		//Customize the request
		req, reqErr := http.NewRequest("GET", rootURL+post.URL, nil)

		if reqErr != nil {
			log.Fatal(reqErr)
			os.Exit(1)
		}
		//Set Header in request to accept response as plain text
		req.Header.Set("Accept", "text/html")
		resp, respErr := client.Do(req)
		if respErr != nil {
			log.Fatal(respErr)
			os.Exit(1)
		}

		data, _ := ioutil.ReadAll(resp.Body)

		links := utils.ParseLinks(string(data))

		fmt.Println("Checking link in post " + post.ID + ":")
		if len(links) == 0 {
			fmt.Println("--- No links to check ---")
		}
		for _, link := range links {
			go utils.CheckLink(link, channel, *userAgent)
		}

		//Receive the result from checkLink and ensure that all links in 1 post returned before moving to next post
		i := 0
		for i < len(links) {
			ls := <-channel
			if ls.GetLiveStatus() == true {
				summary.RecordUpLink(ls.GetURL())
			} else {
				summary.RecordDownLink(ls.GetURL())
			}
			i++
		}
		fmt.Println()
	}
}

//Helper function to remove duplicate code
func checkByFilepath(filepath string, channel chan models.LinkStatus, userAgent string, summary *models.CheckSummary) {
	//Parses links to local variable
	links := utils.ParseLinks(utils.ReadFromFile(filepath))

	//Loop to check all links
	for _, link := range links {
		go utils.CheckLink(link, channel, userAgent)
	}

	//Receive the result from checkLink and update the link to correspondent lists
	i := 0
	for i < len(links) {
		ls := <-channel
		if ls.GetLiveStatus() == true {
			summary.RecordUpLink(ls.GetURL())
		} else {
			summary.RecordDownLink(ls.GetURL())
		}
		i++
	}
}
