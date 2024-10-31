package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

type ScrapedData struct {
	Title string
	URL   string
}

func scrapePage(url string, ch chan<- ScrapedData, wg *sync.WaitGroup) {
	defer wg.Done()
	c := colly.NewCollector()
	var data ScrapedData // to store scraped info

	c.OnHTML("title", func(e *colly.HTMLElement) {
		data.Title = e.Text
		data.URL = url
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Printf("Error while visiting %s: %s", url, err)
	})
	err := c.Visit(url)
	if err != nil {
		log.Printf("Could not visit %s: %s", url, err)
	}

	ch <- data
}

func main() {
	fmt.Println("Enter the URLs to scrape (separated by spaces):")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	urls := strings.Fields(input)

	ch := make(chan ScrapedData) //channel to recieve scraped data
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go scrapePage(url, ch, &wg)
	}

	go func() { //goroutine to close channel
		wg.Wait()
		close(ch)
	}()
	for data := range ch {
		fmt.Printf("Title: %s\nURL: %s\n\n", data.Title, data.URL)
	}
}
