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

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func retryVisit(c *colly.Collector, url string, retries int) error {
	var err error
	for i := 0; i < retries; i++ {
		err = c.Visit(url)
		if err == nil {
			return nil
		}
		log.Printf("Retry %d for %s failed: %s", i+1, url, err)
	}
	return err
}

func writeSummary(url string, success bool, outputDir string) {
	summaryFile := filepath.Join(outputDir, "summary.txt")
	f, err := os.OpenFile(summaryFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error creating summary file: %s", err)
		return
	}
	defer f.Close()

	status := "SUCCESS"
	if !success {
		status = "FAILED"
	}
	f.WriteString(fmt.Sprintf("URL: %s - %s\n", url, status))
}

func scrapePage(url string, outputDir string, wg *sync.WaitGroup) {
	defer wg.Done()
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36"),
	)

	var pageContent string // to store scraped info

	c.OnResponse(func(r *colly.Response) {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(r.Body)))
		if err != nil {
			log.Printf("Error parsing HTML for %s: %s", url, err)
			writeSummary(url, false, outputDir)
			return
		}
		pageContent, _ = doc.Html()
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Printf("Error while visiting %s: %s", url, err)
		writeSummary(url, false, outputDir)
	})
	err := retryVisit(c, url, 3)
	if err != nil {
		log.Printf("Could not visit %s: %s", url, err)
		writeSummary(url, false, outputDir)
		return
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
		writeSummary(url, false, outputDir)
		return
	}
	fmt.Printf("Saved %s to %s\n", url, filePath)
	writeSummary(url, true, outputDir)
}

func readLinksFromFile(fileName string) []string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file %s: %s", fileName, err)
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url != "" {
			urls = append(urls, url)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file %s: %s", fileName, err)
	}
	return urls
}

func initLogger() { //redirects all errors to log.txt file
	logFile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to create log file: %s", err)
	}
	log.SetOutput(logFile)
}

func main() {
	initLogger()
	urls := readLinksFromFile("links.txt") //read urls from links.txt

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
