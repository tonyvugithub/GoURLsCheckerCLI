package outputs

import (
	"fmt"

	"github.com/gookit/color"
)

//DisplayHelpPanel display help panel
func DisplayHelpPanel() {
	fmt.Println("\n***Usage of link detector***")
	fmt.Print("\nGeneral form:\t")
	color.Yellow.Print("linkDetector [flag-options] [file-name]\n\n")
	fmt.Print("Flag options:\n\n")
	color.Yellow.Println("\tOption 1")
	color.Yellow.Println("\tOption 2")
	color.Yellow.Println("\tOption 3")
}
