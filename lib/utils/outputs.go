package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/tonyvugithub/GoURLsCheckerCLI/models"
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

//WriteReportToFile to output the summary into a report
func WriteReportToFile(summary models.CheckSummary) {
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
