package handler

import (
	"net/http"

	"github.com/RATIU5/chewbacca/internal/view/pages"
	"github.com/a-h/templ"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	templ.Handler(pages.ShowNotFound()).ServeHTTP(w, r)
}
