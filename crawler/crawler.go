package crawler

import (
	"fmt"
	"io"
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
	Targets chan string
	// Buffered channel of content Readers
	Output chan io.ReadCloser
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
		count:    0,
		Targets:  make(chan string, 128),
		Output:   make(chan io.ReadCloser, 128),
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

	/*
		When checking for the "Accepts-Ranges" header, you
		have to target the file server, not the HTML page.
		HTML pages don't typically need the header.
	*/
	var res *http.Response
	res, err = http.Head(url)
	if err != nil {
		return nil, fmt.Errorf("Couldn't make HEAD request: %v", err)
	}

	if v := res.Header.Get("Accepts-Ranges"); v != "" {
		out.ranges = true
	}

	// Put initial html page into targets
	out.Targets <- init

	return out, nil
}

// A function to spin off go routines for each new target
func (c *Crawler) Start(count int) {
	// TODO: Remove start, have user implement functionality?
	for c.count < count {
		url, ok := <-c.Targets
		// If channel is closed externally
		if !ok {
			return
		}

		go c.Handle(url)

		c.count++
	}
}

// A function to be spun off for a target (gets next target as well)
func (c *Crawler) Handle(url string) {
	var err error

	// TODO: Use channels to handle these errors
	var doc *goquery.Document
	doc, err = goquery.NewDocument(url)
	if err != nil {
		fmt.Printf("Couldn't get HTML for %v", url)
		return
	}

	go c.Next(doc)
	go c.Content(doc)
}

// A function to get the next target (should be called first)
func (c *Crawler) Next(doc *goquery.Document) {
	// TODO: Use channels to handle errors
	next, ok := doc.Find(c.next).First().Attr(c.nextAttr)
	if !ok {
		fmt.Printf("Couldn't find next for %v", doc.Url)
		return
	}

	c.Targets <- next
}

// A function to get the current content (should handle range restarts if necessary)
func (c *Crawler) Content(doc *goquery.Document) {
	// TODO: Use channels to handle errors
	get, ok := doc.Find(c.content).First().Attr("src")
	if !ok {
		fmt.Printf("Couldn't find src for %v", doc.Url)
		return
	}

	res, err := http.Get(get)
	if err != nil {
		fmt.Printf("Download of %v failed", doc.Url)
		return
	}

	c.Output <- res.Body
}
