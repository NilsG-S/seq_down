package crawler

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// TODO: maybe have the crawler put the response bodies in a chan of io.Reader
// Then let the consumer handle file saving

type Crawler struct {
	// Indicates whether the target website accepts ranges
	ranges bool
	// Buffered channel of URLs to target
	targets chan string
	// ID of the HTML element that gets the next sequential content
	// Doesn't always work with id, can be class and href
	next     string
	nextAttr string
	// ID of the content HTML element
	// Doesn't always work with id, img, and src
	// TODO: future improvement?
	content string
	// contentAttr string
	// Do I need something to handle file naming?
	// If I want a count, I'll need to mutex it
	// TODO: Add stop methods (number of downloads, linked page not the same, re-encounter page, etc.)
	// TODO: Use embedding with interfaces to allow functionality customization
}

func New(init, next, attr, cont string) (*Crawler, error) {
	var err error
	var ok bool
	var out *Crawler = &Crawler{
		next:     next,
		nextAttr: attr,
		content:  cont,
		ranges:   false,
		targets:  make(chan string, 128),
	}

	var doc *goquery.Document
	doc, err = goquery.NewDocument(init)
	if err != nil {
		return nil, fmt.Errorf("Couldn't get initial page: %v", err)
	}

	// Get content URL
	var url string
	url, ok = doc.Find(cont).First().Attr("src")
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
	target <- init

	return out, nil
}

// A function to spin off go routines for each new target

// A function to be spun off for a target (gets next target as well)

// A function to get the next target (should be called first)

// A function to get the current content (should handle range restarts if necessary)
