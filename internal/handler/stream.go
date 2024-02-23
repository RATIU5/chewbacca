package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"github.com/RATIU5/chewbacca/internal/model"
	"github.com/RATIU5/chewbacca/internal/view/components"
	"github.com/gocolly/colly"
)

func StreamHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	queryParams := r.URL.Query()
	addrStr := queryParams.Get("addr")
	addr, err := url.Parse(addrStr)
	if err != nil {
		http.Error(w, "Failed to parse URL", http.StatusInternalServerError)
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()

	c := colly.NewCollector()

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnError(func(r *colly.Response, err error) {
		route := model.NewRoute(*addr, *r.Request.URL, int16(r.StatusCode), "")
		if ctx.Err() != nil {
			return
		} else {
			var buf bytes.Buffer
			components.RowShow(*route).Render(ctx, &buf)

			// Properly format the data for SSE
			fmt.Fprintf(w, "data: %s\n\n", buf.String())

			flusher.Flush()
		}
	})

	c.OnRequest(func(r *colly.Request) {
		route := model.NewRoute(*addr, *r.URL, 200, "")
		if ctx.Err() != nil {
			return
		} else {
			var buf bytes.Buffer
			components.RowShow(*route).Render(ctx, &buf)

			// Properly format the data for SSE
			fmt.Fprintf(w, "data: %s\n\n", buf.String())

			flusher.Flush()
		}
	})

	c.Visit(addrStr)

	if ctx.Err() != nil {
		return
	} else {
		var buf bytes.Buffer
		components.ShowTermStream().Render(ctx, &buf)

		// Properly format the data for SSE
		fmt.Fprintf(w, "data: %s\n\n", buf.String())

		flusher.Flush()
	}
}
