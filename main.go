package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

func main() {
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
