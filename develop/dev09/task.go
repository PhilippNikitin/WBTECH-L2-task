package main

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

// Crawler содержит информацию о том, что нужно скачать
type Crawler struct {
	BaseURL    string
	Downloaded map[string]bool
	Mutex      sync.Mutex
}

// NewCrawler создает новый экземпляр Crawler
func NewCrawler(baseURL string) *Crawler {
	return &Crawler{
		BaseURL:    baseURL,
		Downloaded: make(map[string]bool),
	}
}

// DownloadFile скачивает файл по URL и сохраняет его в указанное имя
func (c *Crawler) DownloadFile(url, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %s", resp.Status)
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

// CrawlPage анализирует HTML страницы и скачивает все найденные ссылки
func (c *Crawler) CrawlPage(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to crawl page: %s", resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}

	var links []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	c.Mutex.Lock()
	c.Downloaded[url] = true
	c.Mutex.Unlock()

	for _, link := range links {
		if !strings.HasPrefix(link, "http") {
			link = c.BaseURL + link
		}

		if !c.Downloaded[link] {
			filename := strings.ReplaceAll(link, "/", "_")
			err := c.DownloadFile(link, filename)
			if err != nil {
				fmt.Printf("Failed to download %s: %v\n", link, err)
			} else {
				fmt.Printf("Downloaded %s\n", link)
			}

			c.Mutex.Lock()
			c.Downloaded[link] = true
			c.Mutex.Unlock()
		}
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <url>")
		return
	}

	baseURL := os.Args[1]
	crawler := NewCrawler(baseURL)

	err := crawler.CrawlPage(baseURL)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
