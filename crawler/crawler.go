package crawler

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// TODO: maybe have the crawler put the response bodies in a chan of io.Reader
// Then let the consumer handle file saving
// TODO: Add stop methods (number of downloads, linked page not the same, re-encounter page, etc.)
// TODO: Use embedding with interfaces to allow functionality customization

type Crawler struct {
	// Indicates whether the target website accepts ranges
	ranges bool
	// Number of downloaded items
	count int
	// Buffered channel of URLs to target
	targets chan string
	// Buffered channel of content Readers
	output chan io.Reader
	// Selector (id, class, etc.) of the HTML element for next
	next string
	// Link can be href, src, etc.
	nextAttr string
	// Selector (id, class, etc.) of the HTML element for content
	// TODO: Generalize for attr types?
	content string
}

func New(init, next, attr, cont string) (*Crawler, error) {
	var err error
	var out *Crawler = &Crawler{
		next:     next,
		nextAttr: attr,
		content:  cont,
		ranges:   false,
		count, 0,
		targets: make(chan string, 128),
		output:  make(chan io.Reader, 128),
	}

	var doc *goquery.Document
	doc, err = goquery.NewDocument(init)
	if err != nil {
		return nil, fmt.Errorf("Couldn't get initial page: %v", err)
	}

	// Get content URL
	url, ok := doc.Find(cont).First().Attr("src")
	if !ok {
		return nil, fmt.Errorf("Couldn't get content source")
	}

	// Checking for "Accepts-Ranges"
	var res *http.Response
	res, err = http.Head(url)
	if err != nil {
		return nil, fmt.Errorf("Couldn't make HEAD request: %v", err)
	}

	if v := res.Header.Get("Accepts-Ranges"); v != "" {
		out.ranges = true
	}

	// Put initial html page into targets
	out.targets <- init

	return out, nil
}

// A function to spin off go routines for each new target
func (c *Crawler) Start(count int) {
	for c.count < count {
		url, ok := <-c.targets
		// If channel is closed externally
		if !ok {
			return
		}

		go Handle(url)

		c.count++
	}
}

// A function to be spun off for a target (gets next target as well)
func (c *Crawler) Handle(url string) {
	var err error

	var doc *goquery.Document
	doc, err = goquery.NewDocument(url)
	if err != nil {
		// TODO: Find some better way of handling these errors
		fmt.Printf("Download of %v failed", url)
		return
	}

	get, ok := doc.Find(c.content).First().Attr("src")
	if !ok {
		fmt.Printf("Download of %v failed", url)
		return
	}

	var res *http.Response
	res, err = http.Get(get)
	if err != nil {
		fmt.Printf("Download of %v failed", url)
		return
	}

	c.output <- res.Body
}

func (c *Crawler) GetOutput() chan io.Reader {
	return c.output
}

// A function to get the next target (should be called first)

// A function to get the current content (should handle range restarts if necessary)
