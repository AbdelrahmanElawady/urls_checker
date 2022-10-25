package pkg

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/fatih/color"
	"golang.org/x/net/html"
)

var (
	DebugColor   = color.New(color.FgBlue)
	SuccessColor = color.New(color.FgGreen)
	ErrorColor   = color.New(color.FgRed)
)

var links map[string]bool = map[string]bool{}

type link struct {
	website string
	URL     string
}

type URLStatus struct {
	Err    string `json:"err,omitempty"`
	Status int64  `json:"status,omitempty"`
	URL    string `json:"url,omitempty"`
}

// Union 2 maps together
func union(m1, m2 map[string]bool) map[string]bool {
	for key, val := range m1 {
		m2[key] = val
	}

	return m2
}

// CheckWebsite checks if the website has a prefix https ..?
func checkWebsite(website string) string {
	if !strings.HasPrefix(website, "https://") {
		website = "https://" + website
	}

	return website
}

// ExtractLinks from a given website
func ExtractLinks(website string, body io.Reader, linkToCheck chan link) map[string]bool {
	htmlTokenizer := html.NewTokenizer(body)
	for {
		token := htmlTokenizer.Next()

		switch token {
		case html.ErrorToken:
			close(linkToCheck)
			return links
		case html.StartTagToken, html.EndTagToken:
			token := htmlTokenizer.Token()
			if "a" == token.Data {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						if _, ok := links[attr.Val]; !ok {
							linkToCheck <- link{website, attr.Val}
							links[attr.Val] = true
						}
					}
				}
			}
		}
	}
}

// Get links inside a given website
func GetLinksOfSource(website string, linkToCheck chan link) error {

	if _, ok := links[website]; ok {
		//close(linkToCheck)
		return nil
	}

	// check the status of the website before getting its links
	res, err := http.Get(website)
	if err != nil {
		close(linkToCheck)
		ErrorColor.Println("Invalid website", err)
		return err
	}

	body := res.Body
	go ExtractLinks(website, body, linkToCheck)
	/*
		for link := range linkToCheck {
			fmt.Println("rawdahanem", link)
			if strings.HasPrefix(link.URL, "/") {
				go GetLinksOfSource(website+link.URL, linkToCheck)
			} else if !strings.HasPrefix(link.URL, "https://") {
				go GetLinksOfSource(website+"/"+link.URL, linkToCheck)
			}
		}
	*/
	return err
}

// TestLink tests the response of get request
// It returns the status code and the error of the request
func TestLink(website, url string) (int, error) {
	if strings.HasPrefix(url, "/") {
		url = website + url
	} else if strings.HasPrefix(url, "http://") {
		return 400, errors.New("url: " + url + " is broken. It starts with http://")
	} else if !strings.HasPrefix(url, "https://") {
		url = website + "/" + url
	}

	res, err := http.Get(url)
	if res == nil {
		return 404, err
	}

	if strings.HasPrefix(url, website) && url != website {
		CheckInternalUrls(url)
	}

	return res.StatusCode, err
}

// Result of the testing links
func result(results chan URLStatus, done chan bool) {
	for result := range results {
		if result.Err != fmt.Sprint(nil) || result.Status >= 400 {
			ErrorColor.Println(result.URL, ": ", fmt.Sprint(result.Status))
			ErrorColor.Println("error: ", fmt.Sprint(result.Err))
		} else {
			SuccessColor.Println(result.URL, ": ", fmt.Sprint(result.Status))
		}
	}
	done <- true
}

// The worker tests the extracted links
func worker(results chan URLStatus, linkToCheck chan link, wg *sync.WaitGroup) {
	for testLink := range linkToCheck {
		status, err := TestLink(testLink.website, testLink.URL)
		results <- URLStatus{fmt.Sprint(err), int64(status), testLink.URL}
	}
	wg.Done()
}

// Creates worker pool
func createWorkerPool(results chan URLStatus, linkToCheck chan link, noOfWorkers int) {
	var wg sync.WaitGroup
	for i := 0; i < noOfWorkers; i++ {
		wg.Add(1)
		go worker(results, linkToCheck, &wg)
	}
	wg.Wait()
	close(results)
}

// Checks the internal urls of a website
func CheckInternalUrls(website string) error {

	linkToCheck := make(chan link)
	results := make(chan URLStatus)

	website = checkWebsite(website)

	go GetLinksOfSource(website, linkToCheck)

	done := make(chan bool)
	go result(results, done)

	noOfWorkers := 10
	createWorkerPool(results, linkToCheck, noOfWorkers)

	<-done

	return nil
}

// Checks the urls of a website
func Check(website string) error {
	DebugColor.Println("start checking: ", website)

	err := CheckInternalUrls(website)

	DebugColor.Println("Finished..")
	return err
}

func checkerMain() {
	website := "codescalers.com"
	err := Check(website)
	if err != nil {
		fmt.Printf("Checking links of %v failed with error: %v\n", website, err)
	}
}
