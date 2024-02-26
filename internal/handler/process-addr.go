package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/RATIU5/chewbacca/internal/model"
)

var (
	processedLinks sync.Map // To store processed links
	wg             sync.WaitGroup
)

func ProcessAddrHandler(w http.ResponseWriter, r *http.Request) {
	rootRoute := r.FormValue("addr")
	rootDomain := getRootDomain(rootRoute) // Implement this function based on your requirements

	// Initialize the map and start scanning from the root route
	processedLinks = sync.Map{}
	wg.Add(1)
	scanRoute(rootRoute, 0, 3, rootDomain) // Passing rootDomain to scanRoute

	wg.Wait() // Wait for all goroutines to finish

	// Iterate over processedLinks and print them
	processedLinks.Range(func(key, value interface{}) bool {
		route := value.(model.Route) // Correct type assertion
		fmt.Printf("Route: %s, URL: %s, Status: %d\n", route.Route, route.URL, route.Status)
		return true // Continue iteration
	})
}

func scanRoute(route string, currentDepth, maxDepth uint8, rootDomain string) {
	defer wg.Done()
	if currentDepth > maxDepth {
		return
	}

	// Check if the route is already processed
	if _, loaded := processedLinks.Load(route); loaded {
		return
	}

	resp, err := http.Get(route)
	if err != nil {
		fmt.Println("Error fetching route:", route, err)
		return
	}
	defer resp.Body.Close()

	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Error parsing the document:", err)
		return
	}

	status := int16(resp.StatusCode)
	newRoute := model.Route{
		Route:  route,
		URL:    route, // Assuming the URL is the route itself
		Status: status,
	}
	processedLinks.Store(route, newRoute)

	document.Find("a").Each(func(index int, element *goquery.Selection) {
		href, exists := element.Attr("href")
		if exists && IsValidURL(href) {
			// Correctly form the full URL if href is a relative path
			absoluteHref := href
			if !strings.HasPrefix(href, "http") {
				absoluteHref = rootDomain + href
			}
			// Ensure the link is within the root domain
			if strings.HasPrefix(absoluteHref, rootDomain) {
				wg.Add(1)
				go func(link string) {
					scanRoute(link, currentDepth+1, maxDepth, rootDomain)
				}(absoluteHref)
			}
		}
	})
}

func IsValidURL(link string) bool {
	// Your existing validation logic
	return !(strings.HasPrefix(link, "tel:") ||
		strings.HasPrefix(link, "mailto:") ||
		strings.HasPrefix(link, "#") ||
		strings.HasPrefix(link, "javascript:"))
}

// getRootDomain extracts the root domain from a given URL.
// Implement this based on your URL structure and validation needs.
func getRootDomain(urlAddr string) string {
	// Simple example implementation, adjust as needed
	parsedURL, err := url.Parse(urlAddr)
	if err != nil {
		return ""
	}
	return parsedURL.Scheme + "://" + parsedURL.Host
}
