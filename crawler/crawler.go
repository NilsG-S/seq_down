package crawler

// TODO: maybe have the crawler put the response bodies in a chan of io.Reader
// Then let the consumer handle file saving

type Crawler struct {
	// Indicates whether the target website accepts ranges
	ranges bool
	// Buffered channel of URLs to target
	targets chan string
	// ID of the HTML element that gets the next sequential content
	next string
	// ID of the content HTML element
	content string
	// Do I need something to handle file naming?
	// If I want a count, I'll need to mutex it
}

/*
	New
	This function takes the following input:
	- Initial target html page
	- ID of the next link
	- ID of the content link

	It determines whether ranges is true or false
	by targeting the content link. It then puts
	the initial target into targets.
*/

// A function to spin off go routines for each new target

// A function to be spun off for a target (gets next target as well)

// A function to get the next target (should be called first)

// A function to get the current content (should handle range restarts if necessary)
