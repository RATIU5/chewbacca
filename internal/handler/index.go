package handler

import (
	"net/http"

	"github.com/RATIU5/chewbacca/internal/view/pages"
	"github.com/a-h/templ"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	templ.Handler(pages.ShowIndex()).ServeHTTP(w, r)
}
