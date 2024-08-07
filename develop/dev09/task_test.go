package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestDownloadFile(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/test" {
			w.Write([]byte("test content"))
		}
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	crawler := NewCrawler(server.URL)
	err := crawler.DownloadFile(server.URL+"/test", "testfile.txt")
	if err != nil {
		t.Fatalf("Failed to download file: %v", err)
	}

	content, err := os.ReadFile("testfile.txt")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if string(content) != "test content" {
		t.Fatalf("Expected file content 'test content', got %s", string(content))
	}

	err = os.Remove("testfile.txt")
	if err != nil {
		t.Fatalf("Failed to remove file: %v", err)
	}
}

func TestCrawlPage(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/page1" {
			w.Write([]byte(`<a href="/page2">Page 2</a>`))
		} else if r.URL.Path == "/page2" {
			w.Write([]byte(`<a href="/page1">Page 1</a>`))
		}
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	crawler := NewCrawler(server.URL)
	err := crawler.CrawlPage(server.URL + "/page1")
	if err != nil {
		t.Fatalf("Failed to crawl page: %v", err)
	}

	if !crawler.Downloaded[server.URL+"/page1"] {
		t.Fatalf("Expected page1 to be downloaded")
	}

	if !crawler.Downloaded[server.URL+"/page2"] {
		t.Fatalf("Expected page2 to be downloaded")
	}
}
