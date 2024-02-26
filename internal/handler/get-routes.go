package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

var visitedUrls = make(map[string]struct{})

func GetRoutesHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now() // Capture the start time

	queryParams := r.URL.Query()
	addrUrl, err := url.ParseRequestURI(queryParams.Get("addr"))
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	baseURL := FormatURL(addrUrl.String())
	addUrlToCache(baseURL)

	host := addrUrl.Hostname()

	domainWithoutWWW := host
	domainWithWWW := "www." + host
	if strings.HasPrefix(host, "www.") {
		domainWithoutWWW = strings.TrimPrefix(host, "www.")
		domainWithWWW = host
	}

	c := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(3),
		colly.AllowedDomains(domainWithoutWWW, domainWithWWW),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 10, // Adjust based on your needs
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		if !ValidateURL(link) {
			return
		}

		url := e.Request.AbsoluteURL(link)
		url = FormatURL(url)

		if !urlInCache(url) {
			addUrlToCache(url)
			// fmt.Println("Visiting", url)
			e.Request.Visit(url)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(r.Request.URL.String() + ": " + err.Error())
	})

	c.Visit(addrUrl.String())
	c.Wait()

	elapsedTime := time.Since(startTime) // Calculate the elapsed time
	fmt.Println("Total URLs visited:", len(visitedUrls))
	fmt.Printf("Total execution time: %s\n", elapsedTime)

	visitedUrls = make(map[string]struct{})
}

func ValidateURL(uri string) bool {
	return !(strings.HasPrefix(uri, "tel:") || strings.HasPrefix(uri, "mailto:") || strings.HasPrefix(uri, "#"))
}

func FormatURL(url string) string {
	var urlFormatted string = url
	if idx := strings.Index(url, "#"); idx != -1 {
		urlFormatted = url[:idx] + "/"
	}

	if !strings.Contains(path.Base(url), ".") && !strings.Contains(url, "?") && !strings.HasSuffix(url, "/") {
		urlFormatted += "/"
	}

	return strings.Replace(urlFormatted, "www.", "", 1)
}

func addUrlToCache(url string) {
	visitedUrls[url] = struct{}{}
}

func urlInCache(url string) bool {
	_, found := visitedUrls[url]
	return found
}
