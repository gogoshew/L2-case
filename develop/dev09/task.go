package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

var (
	fileName    string
	fullUrlFile string
)

func main() {
	fullUrlFile = "http://www.golang-book.com/public/pdf/gobook.pdf"

	// Get filename from path
	fileUrl, err := url.Parse(fullUrlFile)
	if err != nil {
		log.Fatal(err)
	}
	myPath := fileUrl.Path
	seg := strings.Split(myPath, "/")
	fileName = seg[len(seg)-1]

	// Create blank file
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			req.URL.Opaque = req.URL.Path
			return nil
		},
	}

	// Put content on file
	resp, err := client.Get(fullUrlFile)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	size, err := io.Copy(file, resp.Body)
	defer file.Close()

	fmt.Printf("Downloaded a file %s with %d size", fileName, size)
}
