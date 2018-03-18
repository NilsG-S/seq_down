package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/NilsG-S/seq_down/crawler"
)

func main() {
	var (
		err   error
		path  string
		url   string
		cont  string
		next  string
		attr  string
		count int
		pre   string
	)

	fmt.Print("Enter save path: ")
	_, err = fmt.Scanf("%s\n", &path)
	if err != nil {
		fmt.Println("Reading save path failed: ", err)
		return
	}

	fmt.Print("Enter URL: ")
	_, err = fmt.Scanf("%s\n", &url)
	if err != nil {
		fmt.Println("Reading URL failed: ", err)
		return
	}

	fmt.Print("Enter content selector: ")
	_, err = fmt.Scanf("%s\n", &cont)
	if err != nil {
		fmt.Println("Reading content selector failed: ", err)
		return
	}

	fmt.Print("Enter next selector: ")
	_, err = fmt.Scanf("%s\n", &next)
	if err != nil {
		fmt.Println("Reading next selector failed: ", err)
		return
	}

	fmt.Print("Enter next attribute: ")
	_, err = fmt.Scanf("%s\n", &attr)
	if err != nil {
		fmt.Println("Reading next attribute failed: ", err)
		return
	}

	fmt.Print("Enter amount to download: ")
	_, err = fmt.Scanf("%d\n", &count)
	if err != nil {
		fmt.Println("Reading count failed: ", err)
		return
	}

	fmt.Print("Enter file prefix: ")
	_, err = fmt.Scanf("%s\n", &pre)
	if err != nil {
		fmt.Println("Reading prefix failed: ", err)
		return
	}

	craw, err := crawler.New(url, next, attr, cont)
	if err != nil {
		fmt.Println("Setup failed: ", err)
		return
	}
	go craw.Start(count)

	// TODO: this should be in a go routine
	for i := 0; i < count; i++ {
		body, ok := <-craw.Output
		if !ok {
			fmt.Println("Output channel closed")
			return
		}
		defer body.Close()

		// TODO: The response doesn't supply a name, so either have the user provide one or reference the crawled page
		// TODO: The URL should provide a file type
		var file *os.File
		file, err = os.Create(path + pre + strconv.Itoa(i) + ".jpg")
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
