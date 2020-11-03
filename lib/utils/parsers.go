package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

//ParseLinks all valid URLs from a string data
func ParseLinks(data string) []string {
	//Create an regExp object
	re := regexp.MustCompile(`(?i)(?:(?:(?:https?|ftp):)\/\/)(?:\S+(?::\S*)?@)?(?:(x??!(?:10|127)(?:\.\d{1,3}){3})(x??!(?:169\.254|192\.168)(?:\.\d{1,3}){2})(x??!172\.(?:1[6-9]|2\d|3[0-1])(?:\.\d{1,3}){2})(?:[1-9]\d?|1\d\d|2[01]\d|22[0-3])(?:\.(?:1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.(?:[1-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|(?:(?:[a-z0-9\x{00a1}-\x{ffff}][a-z0-9\x{00a1}-\x{ffff}_-]{0,62})?[a-z0-9\x{00a1}-\x{ffff}]\.)+(?:[a-z\x{00a1}-\x{ffff}]{2,}\.?))(?::\d{2,5})?([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`)

	links := re.FindAllString(data, -1)

	return links
}

//ParseIgnoreListPattern to parse link to ignore in check from a text file
func ParseIgnoreListPattern(filePath string) string {
	reg := regexp.MustCompile(`(?m)^#.*$`) // regex to find all comments in ignore file
	fileData := ReadFromFile(filePath)
	fileDataReplace := reg.ReplaceAllString(fileData, "") // delete all comments leaving only links
	ignoreList := ParseLinks(fileDataReplace)             // parses all valid links

	str := strings.Join(ignoreList[:], "|")

	if str != "" {
		regLinkIgnore := regexp.MustCompile("(?m)^.*(" + str + ").*$") // finds all urls in ignore list

		fileDataReplace = regLinkIgnore.ReplaceAllString(fileDataReplace, "") // the urls from ignorelist are taken out of urls to chec
	}

	if strings.TrimSpace(fileDataReplace) != "" { // if filedata is not empty than a bad link still exists
		fmt.Println("Invalid ignore list")
		os.Exit(1)

	}
	return str
}
