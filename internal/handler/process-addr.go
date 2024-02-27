package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/RATIU5/chewbacca/internal/model"
	"github.com/RATIU5/chewbacca/internal/view/components"
	"github.com/a-h/templ"
)

var (
	baseDomain  string
	visited     sync.Map
	maxDepth    int = 3
	pageLinks       = make(map[string][]model.LinkInfo)
	pageLinksMu sync.Mutex
	rateLimit   = make(chan struct{}, 20) // Rate limiting concurrent requests
)

// ProcessAddrHandler handles the /process-addr route
func ProcessAddrHandler(w http.ResponseWriter, r *http.Request) {
	runtime.GOMAXPROCS(4)
	startURL := r.FormValue("addr")
	u, err := url.Parse(startURL)
	if err != nil {
		panic(err)
	}
	baseDomain = u.Hostname()

	crawl(startURL, 0)

	templ.Handler(components.TableResponseShow(pageLinks)).ServeHTTP(w, r)
}

// crawl crawls the given link and its children recursively up to the given depth
//
// It populates the pageLinks map with the links found
func crawl(link string, depth int) {
	if depth > maxDepth {
		return
	}

	// Store the link in the visited map to avoid fetching it again
	if _, loaded := visited.LoadOrStore(link, true); loaded {
		return
	}

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

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Error parsing the page:", link, err)
		return
	}

	var wg sync.WaitGroup
	doc.Find("a, link, script, img").Each(func(i int, s *goquery.Selection) {
		var href string
		var exists bool
		href, exists = s.Attr("href")
		if !exists {
			href, exists = s.Attr("src")
		}
		if !exists || href == "" {
			return
		}

		absoluteURL := resolveURL(link, href)
		linkType := determineLinkType(s)

		wg.Add(1)
		go func() {
			defer wg.Done()
			rateLimit <- struct{}{} // Acquire a slot
			status := fetchStatus(absoluteURL)
			<-rateLimit // Release the slot

			if status != 200 {
				linkInfo := model.LinkInfo{
					URL:    absoluteURL,
					Name:   s.Text(),
					Status: status,
					Type:   linkType,
				}

				pageLinksMu.Lock()
				pageLinks[link] = append(pageLinks[link], linkInfo)
				pageLinksMu.Unlock()
			}

			// Recursively crawl if the link is a route, regardless of its status code
			if linkType == "route" && depth+1 <= maxDepth && sameDomain(absoluteURL, baseDomain) {
				crawl(absoluteURL, depth+1)
			}
		}()
	})

	wg.Wait()
}

func fetchStatus(link string) int {
	resp, err := http.Head(link) // HEAD request to minimize data transfer
	if err != nil {
		return 0 // Unable to fetch status
	}
	defer resp.Body.Close()
	return resp.StatusCode
}

// determineLinkType determines the type of the given HTML element
//
// It can be either "route", "image", "script", "stylesheet" or "unknown"
func determineLinkType(s *goquery.Selection) string {
	if s.Is("a") {
		return "route"
	} else if s.Is("img") {
		return "image"
	} else if s.Is("script") {
		return "script"
	} else if s.Is("link") {
		rel, _ := s.Attr("rel")
		if rel == "stylesheet" {
			return "stylesheet"
		}
	}
	return "unknown"
}

// resolveURL resolves the given href URL against the base URL
//
// It returns the resolved URL as a string
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

// sameDomain checks if the given link is from the same domain as the base domain
//
// It does this by comparing the hostnames of the two URLs
func sameDomain(link, baseDomain string) bool {
	u, err := url.Parse(link)
	if err != nil {
		return false
	}
	return u.Hostname() == baseDomain
}
