package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	// Starting URL
	baseURL := "https://bulbapedia.bulbagarden.net/wiki/List_of_Pokémon_by_effort_value_yield_(Generation_IX)"

	// Fetch the document
	doc, err := fetchDocument(baseURL)
	if err != nil {
		fmt.Println("Error fetching the document:", err)
		return
	}

	// Process the first 100 Pokémon images
	doc.Find("table.sortable tbody tr").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i == 0 {
			return true
		}
		if i > 100 {
			return false
		}
		imgTag := s.Find("td a img")
		src, exists := imgTag.Attr("src")
		if exists {
			fmt.Println("Downloading image:", src)
			downloadImage(src, fmt.Sprintf("pokemon_%d.png", i))
		} else {
			fmt.Println("Image src not found for index:", i)
		}
		return true
	})
}

// fetchDocument fetches the page and returns a goquery document
func fetchDocument(url string) (*goquery.Document, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// downloadImage downloads the image from the given URL and saves it to a file
func downloadImage(url, filename string) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error downloading the image:", err)
		return
	}
	defer response.Body.Close()

	// Create the file
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating the file:", err)
		return
	}
	defer file.Close()

	// Write the body to file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Println("Error writing the image to file:", err)
	}
}