package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func test1() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter URL: ")
	url, _ := reader.ReadString('\n')

	res, _ := http.Head(url)
	maps := res.Header
	/*
		When checking for the "Accepts-Ranges" header, you
		have to target the file server, not the HTML page.
		HTML pages don't typically need the header.
	*/

	for k, v := range maps {
		fmt.Println(k, v)
	}
}

func test2() {
	// Sometimes image links require authorization that is generated on page load
	// TODO: can a web crawler get this authorization?
	// EDIT: yes, it seems that it can
	url := ""
	var err error

	var res *http.Response
	res, err = http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	// The response doesn't supply a name, so either have the user provide one or reference the crawled page
	// The URL should provide a file type
	var file *os.File
	file, err = os.Create("/mnt/c/Users/nilsg/Downloads/test.html")
	if err != nil {
		fmt.Println(err)
	}

	_, err = io.Copy(file, res.Body)
	if err != nil {
		fmt.Println(err)
	}
	file.Close()

	fmt.Println("Success!")
}

func test3() {
	url := ""
	var err error

	var res *http.Response
	res, err = http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	z := html.NewTokenizer(res.Body)
	for {
		tt := z.Next()
		switch tt {
		case html.StartTagToken:
			tn, _ := z.TagName()
			fmt.Println(tn)
		case html.ErrorToken:
			fmt.Println("Error!")
			return
		case html.EndTagToken:
			fmt.Println("Success!")
			return
		}
	}
}

func main() {
	doc, err := goquery.NewDocument("")
	if err != nil {
		fmt.Println("An error has occurred!", err)
	}

	doc.Find("#image").Each(func(i int, s *goquery.Selection) {
		h, e := s.Attr("src")
		if !e {
			fmt.Println("Couldn't get attribute!")
		}

		fmt.Println(h)
	})
}
