package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/RATIU5/chewbacca/internal/model"
	"github.com/RATIU5/chewbacca/internal/view/components"
	"github.com/a-h/templ"
	"github.com/gocolly/colly"
)

var currentURL string
var badURLsList []model.Route

func TestHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	queryParams := r.URL.Query()

	if queryParams.Get("addr") == "" {
		http.Error(w, "Missing addr query parameter", http.StatusBadRequest)
		return
	}

	addrUrl, err := url.Parse(queryParams.Get("addr"))
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	domainWithoutWWW := strings.TrimPrefix(addrUrl.Hostname(), "www.")
	domainWithWWW := "www." + strings.TrimPrefix(addrUrl.Hostname(), "www.")

	c := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(4),
		colly.AllowedDomains(domainWithWWW, domainWithoutWWW),
		colly.CacheDir("./cache"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*.*",
		Parallelism: 10,
	})

	c.OnRequest(func(r *colly.Request) {
		currentURL = r.URL.String()
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if !IsValidURL(link) {
			return
		}

		absoluteURL := e.Request.AbsoluteURL(link)
		formattedURL := FormatURL(absoluteURL)
		e.Request.Visit(formattedURL)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(r.Request.URL.String() + ": " + err.Error())
		badURLsList = append(badURLsList, model.Route{RootAddr: currentURL, CurrentAddr: r.Request.URL.String(), Status: int16(r.StatusCode)})
	})

	c.Visit(addrUrl.String())
	c.Wait()

	elapsedTime := time.Since(startTime)
	fmt.Printf("Total execution time: %s\n", elapsedTime)

	fmt.Println("Total URLs visited:", len(badURLsList))
	templ.Handler(components.TableResponseShow(badURLsList)).ServeHTTP(w, r)
}

func IsValidURL(link string) bool {
	return !(strings.HasPrefix(link, "tel:") ||
		strings.HasPrefix(link, "mailto:") ||
		strings.HasPrefix(link, "#") ||
		strings.HasPrefix(link, "javascript:"))
}
