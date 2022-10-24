package pkg

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/net/html"
)

var (
	DebugColor   = color.New(color.FgBlue)
	SuccessColor = color.New(color.FgGreen)
	ErrorColor   = color.New(color.FgRed)
)

/*type URLStatus struct {
	url    string
	err    string
	status int
}*/

type URLStatus struct {

	// err
	Err string `json:"err,omitempty"`

	// status
	Status int64 `json:"status,omitempty"`

	// url
	URL string `json:"url,omitempty"`
}

func checkWebsite(website string) string {
	if !strings.HasPrefix(website, "https://") {
		website = "https://" + website
	}

	return website
}

// get links inside a given website
func GetLinks(website string) ([]string, error) {
	var links []string

	// check the status of the website before getting its links
	res, err := http.Get(website)
	if err != nil {
		return links, err
	}

	body := res.Body

	htmlTokenizer := html.NewTokenizer(body)
	for {
		token := htmlTokenizer.Next()

		switch token {
		case html.ErrorToken:
			// remove duplicates
			testDuplicate := make(map[string]bool)
			uniqueLinks := []string{}
			for _, item := range links {
				if _, value := testDuplicate[item]; !value {
					testDuplicate[item] = true
					uniqueLinks = append(uniqueLinks, item)
				}
			}

			return uniqueLinks, nil
		case html.StartTagToken, html.EndTagToken:
			token := htmlTokenizer.Token()
			if "a" == token.Data {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						links = append(links, attr.Val)
					}

				}
			}

		}
	}
}

// TestLink tests the response of get request
// It returns the status code and the error of the request
func TestLink(website, url string) (int, error) {
	if strings.HasPrefix(url, "/") {
		url = website + url
	} else if !strings.HasPrefix(url, "https://") {
		url = website + "/" + url
	}

	res, err := http.Get(url)
	if res == nil {
		return 404, err
	}
	return res.StatusCode, err
}

// Checks the urls of a website
func Check(website string) ([]URLStatus, error) {
	DebugColor.Println("start checking: ", website)

	linksStatus := []URLStatus{}

	website = checkWebsite(website)

	links, err := GetLinks(website)
	if err != nil {
		return linksStatus, err
	}
	DebugColor.Println("len(links): ", fmt.Sprint(len(links)))

	results := make(chan URLStatus)

	// initialize routines
	for j := 0; j < len(links)-1; j++ {
		go func(index int) {
			// do processing
			url := links[index]
			var status int
			status, err = TestLink(website, url)
			results <- URLStatus{fmt.Sprint(err), int64(status), url}
		}(j)
	}

	DebugColor.Println("Waiting..")

	for j := 0; j < len(links)-1; j++ {
		r := <-results
		linksStatus = append(linksStatus, r)
		if r.Err != fmt.Sprint(nil) || r.Status >= 400 {
			ErrorColor.Println(r.URL, ": ", fmt.Sprint(r.Status))
			ErrorColor.Println("error: ", fmt.Sprint(r.Err))
		} else {
			SuccessColor.Println(r.URL, ": ", fmt.Sprint(r.Status))
		}
	}

	DebugColor.Println("Finished..")
	return linksStatus, nil
}

func main() {
	website := "codescalers.com"
	_, err := Check(website)
	if err != nil {
		fmt.Printf("Checking links of %v failed with error: %v\n", website, err)
	}
}
