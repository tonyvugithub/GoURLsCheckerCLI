package helpers

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/tonyvugithub/GoURLsCheckerCLI/outputs"
)

//ReadFromFile ...
func ReadFromFile(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	return string(data)
}

//CheckValidArgsLen ...
func CheckValidArgsLen(args []string) {
	if len(args) == 0 {
		outputs.DisplayHelpPanel()
		os.Exit(0)
	}
	if len(args) > 1 {
		fmt.Printf("Too many arguments! Expected exactly 1, Received %+v\n", len(args))
		os.Exit(1)
	}
}
