package main

import (
	"fmt"
	"net/http"

	"github.com/RATIU5/chewbacca/internal/handler"
)

type server struct{}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		handler.IndexHandler(w, r)
	case "/process-addr":
		handler.ProcessAddrHandler(w, r)
	default:
		handler.NotFoundHandler(w, r)
	}
}

func main() {
	var s server
	http.Handle("/", &s)

	// Serve static assets through the /assets/ url path
	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	fmt.Println("Running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
