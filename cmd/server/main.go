package main

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var (
	baseDomain  string
	visited     sync.Map
	maxDepth    int        = 3 // Example depth
	pageLinks              = make(map[string][]string)
	pageLinksMu sync.Mutex // Protects pageLinks
)

func main() {
	startURL := "http://maloufhome.com"
	u, err := url.Parse(startURL)
	if err != nil {
		panic(err)
	}
	baseDomain = u.Hostname()

	crawl(startURL, 0)

	// Print the links found for each page
	for page, links := range pageLinks {
		fmt.Println("Page:", page)
		for _, link := range links {
			fmt.Println(" -", link)
		}
	}
}

func crawl(link string, depth int) {
	if depth > maxDepth {
		return
	}

	// Check if already visited
	if _, loaded := visited.LoadOrStore(link, true); loaded {
		return
	}

	// Fetch the webpage
	resp, err := http.Get(link)
	if err != nil {
		fmt.Println("Error fetching:", link, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Non-OK HTTP status:", resp.StatusCode, link)
		return
	}

	// Parse the webpage
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Error parsing the page:", link, err)
		return
	}

	var linksOnPage []string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			absoluteURL := resolveURL(link, href)
			linksOnPage = append(linksOnPage, absoluteURL)
		}
	})

	// Safely add linksOnPage to the global map
	pageLinksMu.Lock()
	pageLinks[link] = linksOnPage
	pageLinksMu.Unlock()

	var wg sync.WaitGroup
	for _, url := range linksOnPage {
		// Only process if the link is from the same domain
		if sameDomain(url, baseDomain) {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				crawl(url, depth+1)
			}(url)
		}
	}

	wg.Wait()
}

func resolveURL(base, href string) string {
	baseURL, err := url.Parse(base)
	if err != nil {
		return ""
	}
	hrefURL, err := url.Parse(href)
	if err != nil {
		return ""
	}
	return baseURL.ResolveReference(hrefURL).String()
}

func sameDomain(link, baseDomain string) bool {
	u, err := url.Parse(link)
	if err != nil {
		return false
	}
	return u.Hostname() == baseDomain
}
