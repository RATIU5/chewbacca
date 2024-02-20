package main

import (
	"fmt"
	"net/http"

	"github.com/gocolly/colly"
)

func main() {
	// Create a Collector
	c := colly.NewCollector(
		colly.Async(true), // Enable asynchronous requests
	)

	// Set up concurrency limit
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 10, // Adjust the Parallelism for more or fewer threads
	})

	// Common function to handle found URLs
	checkURL := func(element *colly.HTMLElement, attr string) {
		url := element.Attr(attr)
		// Absolute URL with respect to the page
		absoluteURL := element.Request.AbsoluteURL(url)
		// Create a request to check the URL
		element.Request.Visit(absoluteURL)
	}

	// Handle anchor tags and static resources
	c.OnHTML("a[href], img[src], link[href], script[src]", func(e *colly.HTMLElement) {
		if e.Name == "a" {
			checkURL(e, "href")
		} else if e.Name == "img" {
			checkURL(e, "src")
		} else if e.Name == "link" {
			checkURL(e, "href")
		} else if e.Name == "script" {
			checkURL(e, "src")
		}
	})

	// Handle checking for 404s on response
	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode == http.StatusNotFound {
			fmt.Println("404 found at:", r.Request.URL)
		}
	})

	// Error handling (especially useful for catching 404s directly)
	c.OnError(func(r *colly.Response, err error) {
		if r.StatusCode == http.StatusNotFound {
			fmt.Println("404 found at:", r.Request.URL)
		}
	})

	// Start scraping
	c.Visit("http://example.com") // Replace with your target URL

	// Wait until threads are finished
	c.Wait()
}
