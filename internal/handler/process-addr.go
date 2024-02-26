package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/RATIU5/chewbacca/internal/model"
	"github.com/RATIU5/chewbacca/internal/view/components"
	"github.com/a-h/templ"
	"github.com/gocolly/colly"
)

func ProcessAddrHandler(w http.ResponseWriter, r *http.Request) {
	var badURLsList []model.Route
	var mutex sync.Mutex // For safe concurrent access to badURLsList

	startTime := time.Now()
	queryParams := r.FormValue("addr")

	if queryParams == "" {
		templ.Handler(components.ErrResponseShow("A URL to search was not provided")).ServeHTTP(w, r)
		return
	}

	addrUrl, err := url.Parse(queryParams)
	if err != nil {
		templ.Handler(components.ErrResponseShow("An invalid URL was provided")).ServeHTTP(w, r)
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
		r.Ctx.Put("rootURL", r.URL.String())
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		targetURL := e.Attr("href")
		if !IsValidURL(targetURL) {
			return
		}

		targetURL = e.Request.AbsoluteURL(targetURL)
		targetURL = FormatURL(targetURL)
		fmt.Println("Visiting URL:", targetURL, "with title:", e.Text)
		link := model.Link{
			Link:  targetURL,
			Title: e.Text,
		}
		e.Request.Visit(targetURL)
	})

	c.OnError(func(r *colly.Response, err error) {
		mutex.Lock()
		defer mutex.Unlock()
		referrerURL, er := url.Parse(r.Ctx.Get("referrerURL"))
		if er != nil {
			fmt.Println("Error parsing root URL:", err)
		} else {
			url := strings.Split(referrerURL.String(), "?")[0]
			badURLsList = append(badURLsList,
				model.Route{
					RootAddr:     referrerURL,
					RootTitle:    url,
					CurrentAddr:  r.Request.URL,
					CurrentTitle: r.Ctx.Get("currTitle"),
					Status:       int16(r.StatusCode),
				})
		}
	})

	c.Visit(addrUrl.String())
	c.Wait()

	elapsedTime := time.Since(startTime)
	fmt.Printf("Total execution time: %s\n", elapsedTime)

	fmt.Println("Total Bad URLs:", len(badURLsList))
	templ.Handler(components.TableResponseShow(badURLsList)).ServeHTTP(w, r)
}

func IsValidURL(link string) bool {
	return !(strings.HasPrefix(link, "tel:") ||
		strings.HasPrefix(link, "mailto:") ||
		strings.HasPrefix(link, "#") ||
		strings.HasPrefix(link, "javascript:"))
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
