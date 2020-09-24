package outputs

import (
	"fmt"

	"github.com/tonyvugithub/GoURLsCheckerCLI/version"

	"github.com/gookit/color"
)

//DisplayHelpPanel display help panel
func DisplayHelpPanel() {
	fmt.Println("\n***Usage of link detector***")
	fmt.Print("\nGeneral form:\t")
	color.Yellow.Print("linkDetector check [file-name]\n\n")
	fmt.Print("Flag options:\n\n")
	color.Yellow.Print("\t-v / -version")
	fmt.Println("\t: Display app version")
	fmt.Println("\t  Eg: linkDetector -v / -version")

}

//PrintVersion to display version of app
func PrintVersion() {
	fmt.Println("GO URLs CHECKER CLI")
	fmt.Printf("Version: %+v\n", version.BuildVersion)
}

//PrintNoArgsError to display error for no argument
func PrintNoArgsError() {
	fmt.Println("Please provide 1 argument")
	fmt.Println("Eg: ./linkDetector check [file-name]")
}
