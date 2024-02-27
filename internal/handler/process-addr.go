package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/RATIU5/chewbacca/internal/model"
	"github.com/RATIU5/chewbacca/internal/view/components"
	"github.com/a-h/templ"
)

var (
	baseDomain  string
	visited     sync.Map
	maxDepth    int = 0
	pageLinks       = make(map[string][]model.LinkInfo)
	pageLinksMu sync.Mutex
	rateLimit   = make(chan struct{}, 20) // Rate limiting concurrent requests
)

func ProcessAddrHandler(w http.ResponseWriter, r *http.Request) {
	startURL := r.FormValue("addr")
	u, err := url.Parse(startURL)
	if err != nil {
		panic(err)
	}
	baseDomain = u.Hostname()

	crawl(startURL, 0)

	// Print the links found for each page
	templ.Handler(components.TableResponseShow(pageLinks)).ServeHTTP(w, r)
	// for page, links := range pageLinks {
	// 	fmt.Println("Page:", page)
	// 	for _, link := range links {
	// 		fmt.Printf(" - URL: %s, Name: %s, Status: %d, Type: %s\n", link.URL, link.Name, link.Status, link.Type)
	// 	}
	// }
}

func crawl(link string, depth int) {
	if depth > maxDepth {
		return
	}

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
		href, exists := s.Attr("href")
		if !exists {
			src, srcExists := s.Attr("src")
			if srcExists {
				href = src
			}
		}
		if href == "" {
			return
		}

		absoluteURL := resolveURL(link, href)
		linkType := determineLinkType(s)
		linkName := s.Text()

		wg.Add(1)
		go func() {
			defer wg.Done()
			rateLimit <- struct{}{} // Acquire a slot
			status := fetchStatus(absoluteURL)
			<-rateLimit // Release the slot

			linkInfo := model.LinkInfo{
				URL:    absoluteURL,
				Name:   linkName,
				Status: status,
				Type:   linkType,
			}

			pageLinksMu.Lock()
			pageLinks[link] = append(pageLinks[link], linkInfo)
			pageLinksMu.Unlock()
		}()
	})

	wg.Wait()

	if depth < maxDepth {
		for _, linkInfo := range pageLinks[link] {
			if linkInfo.Type == "route" && sameDomain(linkInfo.URL, baseDomain) {
				crawl(linkInfo.URL, depth+1)
			}
		}
	}
}

func fetchStatus(link string) int {
	resp, err := http.Head(link) // HEAD request to minimize data transfer
	if err != nil {
		return 0 // Unable to fetch status, consider handling differently based on requirements
	}
	defer resp.Body.Close()
	return resp.StatusCode
}

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
