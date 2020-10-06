package helpers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gookit/color"
	"github.com/tonyvugithub/GoURLsCheckerCLI/models"
)

//ReadFromFile ...
func ReadFromFile(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return string(data)
}

//ParseLinks ...
func ParseLinks(data string) []string {
	//Create an regExp object
	re := regexp.MustCompile(`(?i)(?:(?:(?:https?|ftp):)\/\/)(?:\S+(?::\S*)?@)?(?:(x??!(?:10|127)(?:\.\d{1,3}){3})(x??!(?:169\.254|192\.168)(?:\.\d{1,3}){2})(x??!172\.(?:1[6-9]|2\d|3[0-1])(?:\.\d{1,3}){2})(?:[1-9]\d?|1\d\d|2[01]\d|22[0-3])(?:\.(?:1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.(?:[1-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|(?:(?:[a-z0-9\x{00a1}-\x{ffff}][a-z0-9\x{00a1}-\x{ffff}_-]{0,62})?[a-z0-9\x{00a1}-\x{ffff}]\.)+(?:[a-z\x{00a1}-\x{ffff}]{2,}\.?))(?::\d{2,5})?([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`)

	links := re.FindAllString(data, -1)

	return links
}

//CheckLink ...
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

//checkLinkWithColor ...
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

//checkLinkNoColor ...
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
