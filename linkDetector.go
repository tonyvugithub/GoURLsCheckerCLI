package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"github.com/gookit/color"
)

//Set up link count info "class"
type linkCountInfo struct {
	upLink    int
	downLink  int
	totalLink int
}

//linkStatus "class"
type linkStatus struct {
	url  string
	live bool
}

func (pt *linkCountInfo) incrementUpLink() {
	(*pt).upLink++
}

func (pt *linkCountInfo) incrementDownLink() {
	(*pt).downLink++
}

func (pt *linkCountInfo) setTotalLink(total int) {
	(*pt).totalLink = total
}

//Function to check whether a link is broken
func checkLink(link string, c chan linkStatus) {
	resp, err := http.Get(link)

	if err != nil {
		color.Gray.Println("[ERROR] " + link)
		c <- linkStatus{
			url:  link,
			live: false,
		}
		return
	}

	statusFormatted := "[" + fmt.Sprint(resp.StatusCode, " ", http.StatusText(resp.StatusCode)) + "]"
	if resp.StatusCode == 200 {
		color.Green.Println(statusFormatted, link)
	} else if resp.StatusCode == 400 || resp.StatusCode == 404 {
		color.Red.Println(statusFormatted, link)
	} else {
		color.Gray.Println(statusFormatted, link)
	}

	c <- linkStatus{
		url:  link,
		live: true,
	}
}

//Function to read content from a file
func readFromFile(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	return string(data)
}

//Function to check if user provided the correct number of args
func checkValidArgsLen(args []string) {
	if len(args) == 0 {
		displayHelpPanel()
		os.Exit(0)
	}
	if len(args) > 1 {
		fmt.Printf("Too many arguments! Expected exactly 1, Received %+v\n", len(args))
		os.Exit(1)
	}
}

func displayHelpPanel() {
	fmt.Println("\n***Usage of link detector***")
	fmt.Print("\nGeneral form:\t")
	color.Yellow.Print("linkDetector [flag-options] [file-name]\n\n")
	fmt.Print("Flag options:\n\n")
	color.Yellow.Println("\tOption 1")
	color.Yellow.Println("\tOption 2")
	color.Yellow.Println("\tOption 3")
}

//Function to parse a valid URL
func parseLinks(data string) []string {
	//Create an regExp object
	//re := regexp.MustCompile(`(?is)(((http|ftp|https):\/{2})+(([0-9a-z_-]+\.)+(aero|asia|biz|cat|com|coop|edu|gov|info|int|jobs|mil|mobi|museum|name|net|org|pro|tel|travel|ac|ad|ae|af|ag|ai|al|am|an|ao|aq|ar|as|at|au|aw|ax|az|ba|bb|bd|be|bf|bg|bh|bi|bj|bm|bn|bo|br|bs|bt|bv|bw|by|bz|ca|cc|cd|cf|cg|ch|ci|ck|cl|cm|cn|co|cr|cu|cv|cx|cy|cz|cz|de|dj|dk|dm|do|dz|ec|ee|eg|er|es|et|eu|fi|fj|fk|fm|fo|fr|ga|gb|gd|ge|gf|gg|gh|gi|gl|gm|gn|gp|gq|gr|gs|gt|gu|gw|gy|hk|hm|hn|hr|ht|hu|id|ie|il|im|in|io|iq|ir|is|it|je|jm|jo|jp|ke|kg|kh|ki|km|kn|kp|kr|kw|ky|kz|la|lb|lc|li|lk|lr|ls|lt|lu|lv|ly|ma|mc|md|me|mg|mh|mk|ml|mn|mn|mo|mp|mr|ms|mt|mu|mv|mw|mx|my|mz|na|nc|ne|nf|ng|ni|nl|no|np|nr|nu|nz|nom|pa|pe|pf|pg|ph|pk|pl|pm|pn|pr|ps|pt|pw|py|qa|re|ra|rs|ru|rw|sa|sb|sc|sd|se|sg|sh|si|sj|sj|sk|sl|sm|sn|so|sr|st|su|sv|sy|sz|tc|td|tf|tg|th|tj|tk|tl|tm|tn|to|tp|tr|tt|tv|tw|tz|ua|ug|uk|us|uy|uz|va|vc|ve|vg|vi|vn|vu|wf|ws|ye|yt|yu|za|zm|zw|arpa)(:[0-9]+)?((\/([~0-9a-zA-Z\#\+\%@\.\/_-]+))?(\?[0-9a-zA-Z\+\%@\/&\[\];=_-]+)?)?))\b`)

	re := regexp.MustCompile(`(?im)https?:\/\/(w{3}.)?.+\.([A-Za-z]+)(:([0-9]+))?\/?[^\s\"\]]*\??(.+=.+&?)?`)

	links := re.FindAllString(data, -1)

	return links
}

func main() {

	var upLinks []string
	var downLinks []string

	channel := make(chan linkStatus)

	//Parse command-line args
	flag.Parse()
	args := flag.Args()

	checkValidArgsLen(args)

	links := parseLinks(readFromFile(args[0]))

	for _, link := range links {
		go checkLink(link, channel)
	}

	i := 0
	for i < len(links) {
		ls := <-channel
		if ls.live == false {
			downLinks = append(downLinks, ls.url)
		} else {
			upLinks = append(upLinks, ls.url)
		}
		i++
	}

	fmt.Println("Total links:", len(links))
	fmt.Println("Up links:", len(upLinks))
	fmt.Println("Down links:", len(downLinks))
}
