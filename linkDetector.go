package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/tonyvugithub/GoURLsCheckerCLI/models"
	"github.com/tonyvugithub/GoURLsCheckerCLI/outputs"
	"github.com/zg3d/GoURLsCheckerCLI/lib/features"
	"github.com/zg3d/GoURLsCheckerCLI/lib/utils"
)

var (
	summary   models.CheckSummary
	wg        sync.WaitGroup
	userAgent *string
)

func main() {
	//Create the channel for routine communication
	channel := make(chan models.LinkStatus)
	//version flag
	versionFlag := flag.Bool("v", false, "version")
	//Create check sub-command
	checkCmd := flag.NewFlagSet("check", flag.ExitOnError)
	//Parse command-line args
	flag.Parse()

	//If there is only program name as argument, print help panel
	if len(os.Args) < 2 {
		outputs.DisplayHelpPanel()
		os.Exit(0)
	}

	//Display the version of the app if -v provided
	features.CheckVersion(versionFlag)

	//Switch statement to consider what subcommand provided
	switch os.Args[1] {
	//Check Subcommand
	case "check":
		flags := os.Args[2:]
		dirFlag := checkCmd.Bool("d", false, "directory path input")
		fileFlag := checkCmd.Bool("f", false, "file path input")
		globFlag := checkCmd.Bool("g", false, "glob pattern")
		reportFlag := checkCmd.Bool("r", false, "check report")
		ignoreFlag := checkCmd.Bool("i", false, "ignore url list")

		//Custom user-agent flag, using default user-agent for Go, access to http.defaultUserAgent deprecated
		userAgent = checkCmd.String("u", "Go-http-client/1.1", "custom user-agent")

		checkCmd.Parse(flags)

		args := checkCmd.Args()

		if *dirFlag && !*fileFlag {

			//if directory flag was provided, check it by directory path
			features.CheckWithDirectoryFlag(globFlag, args, channel, &wg, userAgent, &summary)

		} else if *fileFlag && !*dirFlag {

			//If file flag was provided, check it by Check by filepaths
			features.CheckWithFileFlag(globFlag, args, channel, &wg, userAgent, &summary)

		} else if *globFlag {

			//If provided only -g flag then check the current directory
			//Assign the glob pattern provided to a local variable
			pattern := args[0]
			features.CheckWithGlobFlag(pattern, []string{"."}, channel, &wg, userAgent, &summary)

		} else if *ignoreFlag {

			ignoreList := utils.ParseIgnoreListPattern(args[0]) // gets string with all links seprated by |
			file := args[1]

			features.CheckWithIgnoreFlag(ignoreList, file, channel, userAgent, &summary)

		} else {
			fmt.Println("Invalid format!!! Please try again!!!")
		}

		wg.Wait()

		//If there is a -r flag then report to file report.txt
		if *reportFlag {
			utils.WriteReportToFile(summary)
		}

		break

	default:
		fmt.Println("Expected 'check' command")
		fmt.Println("Eg: $ linkDetector check ...")
		break
	}

	//Exit with error code
	numDownLinks := summary.GetNumDownLinks()
	if numDownLinks > 0 {
		fmt.Println("Exit with status code 1")
		os.Exit(1)
	}

	fmt.Println("Exit with status code 0")
	os.Exit(0)
}
