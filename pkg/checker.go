package pkg

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

const (
	SuccessColor = "\033[1;32m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

/*type URLStatus struct {
	url    string
	err    string
	status int
	index  int
}*/

type URLStatus struct {

	// err
	Err string `json:"err,omitempty"`

	// index
	Index int64 `json:"index,omitempty"`

	// status
	Status int64 `json:"status,omitempty"`

	// url
	URL string `json:"url,omitempty"`
}

func checkWebsite(website string) (string, error) {
	var err error
	websiteTries := [3]string{website, "https://" + website, "http://" + website}
	for _, web := range websiteTries {
		_, err = http.Get(web)
		if err == nil {
			website = web
		}
	}
	return website, err
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
	} else if !strings.HasPrefix(url, "http") {
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
	fmt.Printf(DebugColor, "start checking: "+website)
	fmt.Println("")

	linksStatus := []URLStatus{}

	website, err := checkWebsite(website)
	if err != nil {
		return linksStatus, err
	}

	links, err := GetLinks(website)
	if err != nil {
		return linksStatus, err
	}
	fmt.Printf(DebugColor, "len(links):: "+fmt.Sprint(len(links)))
	fmt.Println("")

	results := make(chan URLStatus)

	// initialize routines
	for j := 0; j < len(links)-1; j++ {
		go func(index int) {
			// do processing
			url := links[index]
			var status int
			status, err = TestLink(website, url)
			results <- URLStatus{fmt.Sprint(err), int64(index), int64(status), url}
		}(j)
	}

	fmt.Printf(DebugColor, "Waiting..")
	fmt.Println("")

	for j := 0; j < len(links)-1; j++ {
		r := <-results
		linksStatus = append(linksStatus, r)
		if r.Err != "" || r.Status >= 400 {
			fmt.Printf(ErrorColor, r.URL+": "+fmt.Sprint(r.Status))
			fmt.Println("")
			fmt.Printf(ErrorColor, "error: "+fmt.Sprint(r.Err))
		} else {
			fmt.Printf(SuccessColor, r.URL+": "+fmt.Sprint(r.Status))
		}
		fmt.Println("")
	}

	fmt.Printf(DebugColor, "Finished..")
	fmt.Println("")
	return linksStatus, nil
}

func main() {
	website := "codescalers.com"
	_, err := Check(website)
	if err != nil {
		fmt.Printf("Checking links of %v failed with error: %v\n", website, err)
	}
}
