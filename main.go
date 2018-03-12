package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/NilsG-S/seq_down/crawler"
)

func main() {
	var err error

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter URL: ")
	url, _ := reader.ReadString('\n')

	fmt.Print("Enter content selector: ")
	cont, _ := reader.ReadString('\n')

	fmt.Print("Enter next selector: ")
	next, _ := reader.ReadString('\n')

	fmt.Print("Enter next attribute: ")
	attr, _ := reader.ReadString('\n')

	fmt.Print("Enter amount to download: ")
	countStr, _ := reader.ReadString('\n')
	count, _ := strconv.Atoi(countStr)
	fmt.Println(cont)

	crawler, err := crawler.New(url, next, attr, cont)
	if err != nil {
		fmt.Println(err)
	}
	crawler.Start(count)

	// TODO: this should be in a go routine
	for i := 0; i < count; i++ {
		body, ok := <-crawler.Output
		if !ok {
			fmt.Println("Output channel closed")
			return
		}
		defer body.Close()

		// TODO: The response doesn't supply a name, so either have the user provide one or reference the crawled page
		// TODO: The URL should provide a file type
		var file *os.File
		file, err = os.Create("/mnt/c/Users/nilsg/Downloads/test/" + strconv.Itoa(i) + ".jpg")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		_, err = io.Copy(file, body)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
