package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/RATIU5/chewbacca/internal/view/components"
)

func StreamHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
	}

	ctx := r.Context() // Get the request's context

	for {
			select {
			case <-ctx.Done(): // Check if client connection is closed
					fmt.Println("Client has closed the connection")
					return
			default:
					var buf bytes.Buffer
					// Assuming components.RowShow().Render writes HTML content to buf
					components.RowShow().Render(ctx, &buf)

					// Properly format the data for SSE
					fmt.Fprintf(w, "data: %s\n\n", buf.String())

					// Flush the data immediately after writing
					flusher.Flush()

					time.Sleep(1 * time.Second) // Be cautious with production use
			}
	}
}