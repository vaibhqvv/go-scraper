package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

func scrapePage(url string, outputDir string, wg *sync.WaitGroup) {
	defer wg.Done()
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36"),
	)

	var pageContent string // to store scraped info

	c.OnResponse(func(r *colly.Response) {
		pageContent = string(r.Body)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Printf("Error while visiting %s: %s", url, err)
	})
	err := c.Visit(url)
	if err != nil {
		log.Printf("Could not visit %s: %s", url, err)
	}

	// to generate filenaem for scraped info
	fileName := strings.ReplaceAll(url, "https://", "")
	fileName = strings.ReplaceAll(fileName, "http://", "")
	fileName = strings.ReplaceAll(fileName, "/", "_")
	fileName = strings.Split(fileName, "?")[0] //to remove anything after '?'
	fileName += ".html"

	filePath := filepath.Join(outputDir, fileName)

	err = os.WriteFile(filePath, []byte(pageContent), 0644)
	if err != nil {
		log.Printf("Error saving %s: %s", filePath, err)
		return
	}
	fmt.Printf("Saved %s to %s\n", url, filePath)
}

func main() {
	fmt.Println("Enter the URLs to scrape (separated by spaces):")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	urls := strings.Fields(input)

	outputDir := "scrapedHTML"
	err := os.MkdirAll(outputDir, os.ModePerm) //if directory isnt there, create one
	if err != nil {
		log.Fatalf("Failed to create directory %s: %s", outputDir, err)
	}
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go scrapePage(url, outputDir, &wg)
	}
	wg.Wait()
	fmt.Println("All pages have been scraped and saved.")

	http.Handle("/", http.FileServer(http.Dir(outputDir)))
	fmt.Println("Serving scraped files at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
