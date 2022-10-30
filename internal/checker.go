package internal

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"golang.org/x/net/html"
)

var (
	DebugColor   = color.New(color.FgBlue)
	SuccessColor = color.New(color.FgGreen)
	ErrorColor   = color.New(color.FgRed)
)

type CheckerService struct {
	noOfWorkers int
	links       map[string]bool

	linkToCheck chan Link
	results     chan URLStatus
	wg          sync.WaitGroup

	conn   *websocket.Conn
	apiOut []URLStatus
	client http.Client
}

type Link struct {
	website string
	URL     string
}

type URLStatus struct {
	Err    string
	Status int64
	URL    string
}

// NewParser creates new instance from the parser
func NewCheckerService(noOfWorkers int) CheckerService {
	return CheckerService{
		noOfWorkers: noOfWorkers,
		links:       map[string]bool{},
		linkToCheck: make(chan Link),
		results:     make(chan URLStatus),

		client: http.Client{Timeout: 5 * time.Second},
	}
}

// Start checker service
func (c *CheckerService) Start() {
	DebugColor.Println("start checking: ")

	web := <-c.linkToCheck
	c.checkInternalUrls(web.website)

	DebugColor.Println("Finished..")
}

// Add a new site to the checker service
func (c *CheckerService) AddSite(site string) {
	go func(website string) {
		website = checkWebsite(website)
		c.linkToCheck <- Link{website, ""}
	}(site)
}

func (c *CheckerService) AddSocket(conn *websocket.Conn) {
	c.conn = conn
}

func (c *CheckerService) AddApiOutput() {
	c.apiOut = []URLStatus{}
}

func (c *CheckerService) GetApiOutput() []URLStatus {
	return c.apiOut
}

// ExtractLinks from a given website
func (c *CheckerService) extractLinks(website string, body io.Reader) {
	htmlTokenizer := html.NewTokenizer(body)
	for {
		token := htmlTokenizer.Next()

		switch token {
		case html.ErrorToken:
			return
		case html.StartTagToken, html.EndTagToken:
			token := htmlTokenizer.Token()
			if "a" == token.Data {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						if _, ok := c.links[attr.Val]; !ok {
							c.linkToCheck <- Link{website, attr.Val}
							c.links[attr.Val] = true
						}
					}
				}
			}
		}
	}
}

// Get links inside a given website
func (c *CheckerService) getLinksOfSource(website string) {

	if _, ok := c.links[website]; ok {
		return
	}

	// check the status of the website before getting its links
	res, err := c.client.Get(website)
	if err != nil {
		c.results <- URLStatus{fmt.Sprint(err), int64(404), website}
		return
	}

	body := res.Body
	go c.extractLinks(website, body)
}

// TestLink tests the response of get request
// It returns the status code and the error of the request
func (c *CheckerService) testLink(website, url string) (int, error) {
	if strings.HasPrefix(url, "/") {
		url = website + url
	} else if strings.HasPrefix(url, "http://") {
		return 400, errors.New("url: " + url + " is broken. It starts with http://")
	} else if !strings.HasPrefix(url, "https://") {
		url = website + "/" + url
	}

	res, err := c.client.Get(url)
	if res == nil {
		return 404, err
	}

	if strings.HasPrefix(url, website) && url != website {
		c.getLinksOfSource(url)
	}
	return res.StatusCode, err
}

// Result of the testing links
func (c *CheckerService) result(done chan bool) {
	loop := true
	for loop {
		select {
		case result := <-c.results:

			if c.conn != nil {
				err := c.conn.WriteJSON(result)
				if err != nil {
					ErrorColor.Println("write failed: ", err)
				}
			} else if c.apiOut != nil {
				c.apiOut = append(c.apiOut, result)
			} else if result.Err != fmt.Sprint(nil) || result.Status >= 400 {
				ErrorColor.Println(result.URL, ": ", fmt.Sprint(result.Status))
				ErrorColor.Println("error: ", fmt.Sprint(result.Err))
			} else {
				SuccessColor.Println(result.URL, ": ", fmt.Sprint(result.Status))
			}

		case <-time.After(4 * time.Second):
			loop = false
		}

	}
	close(c.linkToCheck)
	done <- true
}

// The worker tests the extracted links
func (c *CheckerService) worker(wg *sync.WaitGroup) {
	for linkToTest := range c.linkToCheck {
		status, err := c.testLink(linkToTest.website, linkToTest.URL)
		//time.Sleep(2 * time.Second)
		c.results <- URLStatus{fmt.Sprint(err), int64(status), linkToTest.URL}
	}
	wg.Done()
}

// Creates worker pool
func (c *CheckerService) createWorkerPool() {
	for i := 0; i < c.noOfWorkers; i++ {
		c.wg.Add(1)
		go c.worker(&c.wg)
	}
	c.wg.Wait()
	close(c.results)
}

// Checks the internal urls of a website
func (c *CheckerService) checkInternalUrls(website string) {

	go c.getLinksOfSource(website)

	done := make(chan bool)
	go c.result(done)

	c.createWorkerPool()

	<-done
}
