# **GoScraper - Multithreaded Web Scraper with File Server**  

## **Overview**  
GoScraper is a fast, multithreaded web scraper built in **Go** using the **Colly** library. It reads URLs from a `links.txt` file, scrapes the HTML content, saves it locally, and serves the scraped files via an HTTP server.

## **Features**  
- **Multithreaded Scraping** - Concurrently scrapes multiple websites for efficiency.  
- **Colly Web Scraping** - Uses **Go Colly**, a web scraping library.  
- **Automatic URL Handling** - Reads target websites from `links.txt`.  
- **Custom User-Agent** - Mimics real browsers to avoid bot detection.  
- **File-Based Storage** - Saves scraped HTML files in `scrapedHTML/`.

## **Technologies Used**  
- **Go** (Golang)  
- **Colly** (Web Scraping)  
- **HTTP Server** (for serving scraped files)   

## **Installation**  

### **Clone the Repository**  
```sh
git clone https://github.com/yourusername/go-scraper.git
cd go-scraper
```

### **Install Dependencies**  
Ensure you have **Go** installed. Then, install Colly:  
```sh
go get -u github.com/gocolly/colly
```

### **Build the Scraper**  
```sh
go build -o scraper scraper.go
```

### **Prepare `links.txt`**  
Add website URLs **(one per line)** in `links.txt`. Example:  
```
https://example.com
https://golang.org
https://github.com
```

## **Usage**  

### **Run the Scraper**  
```sh
./scraper
```
- Scrapes all websites listed in `links.txt`.  
- Saves HTML files in `scrapedHTML/`.  
- Starts a web server at **http://localhost:8080/** to serve scraped content.    

## **How It Works**  

1. **Reads URLs from `links.txt`**  
2. **Scrapes each website using Colly**  
3. **Saves the HTML content locally**  
4. **Hosts the scraped files via an HTTP server**  

## **License**  
This project is open-source under the **MIT License**.
