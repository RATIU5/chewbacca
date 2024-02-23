package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/RATIU5/chewbacca/internal/model"
	"github.com/RATIU5/chewbacca/internal/view/components"
)

func StreamHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	url1, _ := url.Parse("https://example.com/")
	url2, _ := url.Parse("https://example.com/test")

	route := model.NewRoute(*url1, *url2, 200, "Example Website")

	for {
		var wr io.Writer
		components.RowShow(*route).Render(context.Background(), wr)

		fmt.Println(wr)

		// Flush the data immediately
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		} else {
			fmt.Println("Failed to flush")
		}

		time.Sleep(1 * time.Second)
	}
}
