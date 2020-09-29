package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/tonyvugithub/GoURLsCheckerCLI/helpers"
	"github.com/tonyvugithub/GoURLsCheckerCLI/models"
	"github.com/tonyvugithub/GoURLsCheckerCLI/outputs"
)

var (
	summary   models.CheckSummary
	upLinks   []string
	downLinks []string
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

	//Check if version flag was provided
	if *flagVersionLong || *flagVersionShort {
		//If exactly 2 arguments, print the name and version of the app
		if len(os.Args) == 2 {
			outputs.PrintVersion()
			os.Exit(0)
		} else {
			fmt.Println("There should be no extra argument after -v/-version")
		}
	}

	//Switch statement to consider what subcommand provided
	switch os.Args[1] {
	//Check Subcommand
	case "check":
		flags := os.Args[2:]
		dirFlag := checkCmd.Bool("d", false, "directory path input")
		fileFlag := checkCmd.Bool("f", false, "file path input")
		reportFlag := checkCmd.Bool("r", false, "check report")

		checkCmd.Parse(flags)
		args := checkCmd.Args()
		//helpers.CheckValidArgsLen(args)

		var wg sync.WaitGroup

		//if directory flag was provided, check it by directory path
		if *dirFlag && !*fileFlag {
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
						checkByFilepath(f, channel)
					}(filepath)
				}
			}
		}

		//If file flag was provided, check it by file path
		if *fileFlag && !*dirFlag {
			//Check by filenames
			for _, file := range args {
				wg.Add(1)
				go func(f string) {
					defer wg.Done()
					checkByFilepath(f, channel)
				}(file)
			}
			//Any other format would be invalid
		}

		wg.Wait()

		if *reportFlag && (*fileFlag || *dirFlag) {
			writeReportToFile()
		} else {
			fmt.Println("Invalid format!!! Please try again!!!")
		}
		break

	default:
		fmt.Println("Expected 'check' command")
		fmt.Println("Eg: $ linkDetector check ...")
		os.Exit(1)
	}

	numUpLinks := summary.GetNumUpLinks()
	numDownLinks := summary.GetNumDownLinks()
	fmt.Println("Total links:", numUpLinks+numDownLinks)
	fmt.Println("Up links:", numUpLinks)
	fmt.Println("Down links:", numDownLinks)
}

func checkByFilepath(filepath string, channel chan models.LinkStatus) {
	//Parses links to local variable
	links := helpers.ParseLinks(helpers.ReadFromFile(filepath))

	//Loop to check all links
	for _, link := range links {
		go helpers.CheckLink(link, channel)
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

	f.WriteString("\nDOWN LINKS list:\n")
	for _, link := range summary.GetUpLinks() {
		f.WriteString(link + "\n")
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
