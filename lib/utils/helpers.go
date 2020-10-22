package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gookit/color"
	"github.com/tonyvugithub/GoURLsCheckerCLI/models"
)

//ReadFromFile to read a file and converse to string data
func ReadFromFile(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return string(data)
}

//CheckLink to validate the link status
func CheckLink(link string, c chan models.LinkStatus, userAgent string) {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	req, err := http.NewRequest("HEAD", link, nil)

	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)

	if err != nil {
		color.Gray.Println("[ERROR] " + link)
		ls := models.LinkStatus{}
		ls.SetURL(link)
		ls.SetLiveStatus(false)
		c <- ls
		return
	}

	clicolorEnv := os.Getenv("CLICOLOR")

	if clicolorEnv == "0" {
		checkLinkNoColor(resp, link)
	} else if clicolorEnv == "1" || clicolorEnv == "" {
		checkLinkWithColor(resp, link)
	}

	ls := models.LinkStatus{}
	ls.SetURL(link)
	ls.SetLiveStatus(true)
	c <- ls
}

//checkLinkWithColor to output result of link checking with color if "CLICOLOR" option turned on
func checkLinkWithColor(resp *http.Response, link string) {
	statusFormatted := "[" + fmt.Sprint(resp.StatusCode, " ", http.StatusText(resp.StatusCode)) + "]"
	if resp.StatusCode == 200 {
		color.Green.Println(statusFormatted, link)
	} else if resp.StatusCode == 400 || resp.StatusCode == 404 {
		color.Red.Println(statusFormatted, link)
	} else {
		color.Gray.Println(statusFormatted, link)
	}
}

//checkLinkNoColor to output result of link checking with no color if "CLICOLOR" option turned off
func checkLinkNoColor(resp *http.Response, link string) {
	statusFormatted := "[" + fmt.Sprint(resp.StatusCode, " ", http.StatusText(resp.StatusCode)) + "]"
	if resp.StatusCode == 200 {
		fmt.Println(statusFormatted, link)
	} else if resp.StatusCode == 400 || resp.StatusCode == 404 {
		fmt.Println(statusFormatted, link)
	} else {
		fmt.Println(statusFormatted, link)
	}
}
