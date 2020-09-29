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

	//If there is only program name as argument
	if len(os.Args) < 2 {
		outputs.DisplayHelpPanel()
		os.Exit(0)
	}

	//Check if version flag was provided
	if *flagVersionLong || *flagVersionShort {
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

		checkCmd.Parse(flags)
		args := checkCmd.Args()
		//helpers.CheckValidArgsLen(args)

		var wg sync.WaitGroup

		//if directory flag was provided, check it by directory path
		if *dirFlag {
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
			//If file flag was provided, check it by file path
		} else if *fileFlag {
			//Check by filenames
			for _, file := range args {
				wg.Add(1)
				go func(f string) {
					defer wg.Done()
					checkByFilepath(f, channel)
				}(file)
			}
		} else {
			fmt.Println("Invalid format!!! Please try again!!!")
		}
		wg.Wait()

		fmt.Println("Total links:", len(upLinks)+len(downLinks))
		fmt.Println("Up links:", len(upLinks))
		fmt.Println("Down links:", len(downLinks))
	default:
		fmt.Println("Expected 'check' command")
		fmt.Println("Eg: $ linkDetector check ...")
		os.Exit(1)
	}
}

func checkByFilepath(filepath string, channel chan models.LinkStatus) {
	links := helpers.ParseLinks(helpers.ReadFromFile(filepath))

	for _, link := range links {
		go helpers.CheckLink(link, channel)
	}
	//Receive the result from checkLink and update the link to correspondent lists
	i := 0
	for i < len(links) {
		ls := <-channel
		if ls.GetLiveStatus() == false {
			downLinks = append(downLinks, ls.GetURL())
		} else {
			upLinks = append(upLinks, ls.GetURL())
		}
		i++
	}
}
