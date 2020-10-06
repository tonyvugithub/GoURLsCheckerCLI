package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/mb0/glob"
	"github.com/tonyvugithub/GoURLsCheckerCLI/helpers"
	"github.com/tonyvugithub/GoURLsCheckerCLI/models"
	"github.com/tonyvugithub/GoURLsCheckerCLI/outputs"
)

var (
	summary   models.CheckSummary
	wg        sync.WaitGroup
	userAgent *string
)

func main() {

	channel := make(chan models.LinkStatus)
	//version flag
	flagVersionLong := flag.Bool("version", false, "version")
	flagVersionShort := flag.Bool("v", false, "version")
	//Create check sub-command
	checkCmd := flag.NewFlagSet("check", flag.ExitOnError)
	//Parse command-line args
	flag.Parse()

	//If there is only program name as argument, print help panel
	if len(os.Args) < 2 {
		outputs.DisplayHelpPanel()
		os.Exit(0)
	}

	//Check if version flag was provided avsc
	if *flagVersionLong || *flagVersionShort {
		//If exactly 2 arguments, print the name and version of the app
		if len(os.Args) == 2 {
			outputs.PrintVersion()
			os.Exit(0)
		} else {
			fmt.Println("There should be no extra argument after -v/-version")
			os.Exit(1)
		}
	}

	//Switch statement to consider what subcommand provided
	switch os.Args[1] {
	//Check Subcommand
	case "check":
		flags := os.Args[2:]
		dirFlag := checkCmd.Bool("d", false, "directory path input")
		fileFlag := checkCmd.Bool("f", false, "file path input")
		globFlag := checkCmd.Bool("g", false, "glob pattern")
		reportFlag := checkCmd.Bool("r", false, "check report")

		//Custom user-agent flag, using default user-agent for Go, access to http.defaultUserAgent deprecated
		userAgent = checkCmd.String("u", "Go-http-client/1.1", "custom user-agent")

		checkCmd.Parse(flags)

		args := checkCmd.Args()
		//helpers.CheckValidArgsLen(args)

		//if directory flag was provided, check it by directory path
		if *dirFlag && !*fileFlag {
			//If glob flag also provided
			if *globFlag {
				argsWithoutPattern := args[:len(args)-1]
				//Assign the glob pattern provided to a local variable
				pattern := args[len(args)-1]
				checkWithGlobPattern(pattern, argsWithoutPattern, channel)
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
							checkByFilepath(f, channel, *userAgent)
						}(filepath)
					}
				}
			}
		} else if *fileFlag && !*dirFlag {
			//If file flag was provided, check it by file path
			//Check by filenames
			if *globFlag {
				fmt.Println("Invalid format!!! Please try again!!!")
			} else {
				for _, file := range args {
					wg.Add(1)
					go func(f string) {
						defer wg.Done()
						checkByFilepath(f, channel, *userAgent)
					}(file)
				}
			}
		} else if *globFlag {
			//Assign the glob pattern provided to a local variable
			pattern := args[0]
			checkWithGlobPattern(pattern, []string{"."}, channel)
		} else {
			fmt.Println("Invalid format!!! Please try again!!!")
		}

		wg.Wait()

		//If there is a report flag then report to file report.txt
		if *reportFlag {
			writeReportToFile()
		}

		break

	default:
		fmt.Println("Expected 'check' command")
		fmt.Println("Eg: $ linkDetector check ...")
		break
	}
}

func checkWithGlobPattern(pattern string, dirList []string, channel chan models.LinkStatus) {
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
					checkByFilepath(f, channel, *userAgent)
				}(filepath)
			}
		}
	}
}

func checkByFilepath(filepath string, channel chan models.LinkStatus, userAgent string) {
	//Parses links to local variable
	links := helpers.ParseLinks(helpers.ReadFromFile(filepath))

	//Loop to check all links
	for _, link := range links {
		go helpers.CheckLink(link, channel, userAgent)
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

func writeReportToFile() {
	f, err := os.OpenFile("report.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}

	numUpLinks := summary.GetNumUpLinks()
	numDownLinks := summary.GetNumDownLinks()

	f.WriteString("CHECK REPORT\n\n")

	f.WriteString("Total number of links checked: " + fmt.Sprint(numUpLinks+numDownLinks) + "\n")
	f.WriteString("Total number of up links: " + fmt.Sprint(numUpLinks) + "\n")
	f.WriteString("Total number of down links: " + fmt.Sprint(numDownLinks) + "\n")

	f.WriteString("\nDOWN LINKS list:\n")
	for _, link := range summary.GetDownLinks() {
		f.WriteString(link + "\n")
	}

	f.WriteString("\nUP LINKS list:\n")
	for _, link := range summary.GetUpLinks() {
		f.WriteString(link + "\n")
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
