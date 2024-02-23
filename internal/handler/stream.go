package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/RATIU5/chewbacca/internal/model"
	"golang.org/x/net/html"
)

var visited map[string]bool

func init() {
	visited = make(map[string]bool)
}

func processPage(ctx context.Context, w http.ResponseWriter, flusher http.Flusher, rootAddr, currentAddr url.URL, depth int) {
	currentURLString := currentAddr.String()

	// Check if the URL has already been visited.
	if visited[currentURLString] || depth > 10 {
		return
	}

	// Mark the URL as visited.
	visited[currentURLString] = true

	resp, err := http.Get(currentURLString)
	if err != nil {
		fmt.Fprintf(w, "data: Error fetching %s\n\n", currentURLString)
		flusher.Flush()
		return
	}
	defer resp.Body.Close()

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "text/html") {
		return
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	body := bytes.NewReader(bodyBytes)
	title := extractTitle(body)

	// Construct and stream the Route data.
	route := model.Route{
		RootAddr:    rootAddr,
		CurrentAddr: currentAddr,
		Status:      int16(resp.StatusCode),
		Title:       title,
	}
	streamRouteData(w, flusher, route)

	// Find, resolve, and process all unique links on the page.
	body.Seek(0, io.SeekStart)
	links := extractLinks(body, currentAddr)
	for _, link := range links {
		absoluteURL := resolveURL(link, currentAddr)
		if absoluteURL != "" && !visited[absoluteURL] {
			nextURL, err := url.Parse(absoluteURL)
			if err != nil {
				continue
			}
			processPage(ctx, w, flusher, rootAddr, *nextURL, depth+1)
		}
	}
}

func resolveURL(link string, base url.URL) string {
	linkURL, err := url.Parse(link)
	if err != nil {
		return ""
	}
	resolvedURL := base.ResolveReference(linkURL)
	return resolvedURL.String()
}

func extractLinks(body io.Reader, base url.URL) []string {
	var links []string
	z := html.NewTokenizer(body)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return links
		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()
			if t.Data == "a" {
				for _, attr := range t.Attr {
					if attr.Key == "href" {
						links = append(links, attr.Val)
						break
					}
				}
			}
		}
	}
}

func streamRouteData(w http.ResponseWriter, flusher http.Flusher, route model.Route) {
	// Assuming components.RowShow exists and can render the Route to HTML or a string
	// Adapt this line to your actual rendering logic
	fmt.Fprintf(w, "data: Processed: %s, Status: %d, Title: %s\n\n", route.CurrentAddr.String(), route.Status, route.Title)
	flusher.Flush()
}

func extractTitle(body io.Reader) string {
	z := html.NewTokenizer(body)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return ""
		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()
			if t.Data == "title" {
				z.Next()
				return z.Token().Data
			}
		}
	}
}
