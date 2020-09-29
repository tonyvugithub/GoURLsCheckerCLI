package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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
					fmt.Println("Cannot read", dirPath)
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
			fmt.Println("Printing report...")
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
