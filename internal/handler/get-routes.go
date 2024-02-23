package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/gocolly/colly"
)

func GetRoutesHandler(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()

	addrUrl, err := url.ParseRequestURI(queryParams.Get("addr"))
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	c := colly.NewCollector(
		colly.AllowedDomains(addrUrl.Host),
		colly.Async(true),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if !ValidateURL(e.Attr("href")) {
			return
		}

		url := e.Request.AbsoluteURL(e.Attr("href"))

		// Format the URL
		url = FormatURL(url)

		fmt.Println("Visiting", url)

		// e.Request.Visit(url)
	})

	// Target all image elements
	c.OnHTML("img[src]", func(e *colly.HTMLElement) {
		if !ValidateURL(e.Attr("src")) {
			return
		}

		url := e.Request.AbsoluteURL(e.Attr("src"))

		// Format the URL
		url = FormatURL(url)

		fmt.Println("Visiting", url)

		e.Request.Visit(url)
	})

	c.Visit(addrUrl.String())
}

func ValidateURL(uri string) bool {
	// If the link is a telephone number or email address, ignore it
	if strings.HasPrefix(uri, "tel:") || strings.HasPrefix(uri, "mailto:") {
		return false
	}

	// If the link is to an id on the same page, ignore it
	if strings.HasPrefix(uri, "#") {
		return false
	}
	return true
}

func FormatURL(url string) string {
	var urlFormatted string = url
	// Trim any possible id from the URL
	if strings.Contains(url, "#") {
		urlFormatted = url[:strings.Index(url, "#")] + "/"
	}

	// Add a trailing slash if the URL doesn't have:
	// - a file extension
	// - a query strin
	if !strings.Contains(path.Base(url), ".") && !strings.Contains(url, "?") && !strings.Contains(url, "#") && !strings.HasSuffix(url, "/") {
		urlFormatted += "/"
	}

	return urlFormatted
}
