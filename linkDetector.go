package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tonyvugithub/GoURLsCheckerCLI/helpers"
	"github.com/tonyvugithub/GoURLsCheckerCLI/models"
	"github.com/tonyvugithub/GoURLsCheckerCLI/outputs"
)

func main() {
	var upLinks []string
	var downLinks []string

	channel := make(chan models.LinkStatus)
	//version flag
	flagVersionLong := flag.Bool("version", false, "version")
	flagVersionShort := flag.Bool("v", false, "version")
	//Create check sub-command
	checkCmd := flag.NewFlagSet("check", flag.ExitOnError)
	//Parse command-line args
	flag.Parse()

	if *flagVersionLong || *flagVersionShort {
		fmt.Println("Print version")
		return
	}

	if len(os.Args) < 2 {
		outputs.DisplayHelpPanel()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "check":
		flags := os.Args[2:]

		checkCmd.Parse(flags)
		args := checkCmd.Args()
		helpers.CheckValidArgsLen(args)

		links := helpers.ParseLinks(helpers.ReadFromFile(args[0]))

		for _, link := range links {
			go helpers.CheckLink(link, channel)
		}

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

		fmt.Println("Total links:", len(links))
		fmt.Println("Up links:", len(upLinks))
		fmt.Println("Down links:", len(downLinks))
	default:
		fmt.Println("Expected 'check' command")
		fmt.Println("Eg: $ linkDetector check ...")
		os.Exit(1)
	}
}
