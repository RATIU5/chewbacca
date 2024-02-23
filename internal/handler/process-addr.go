package handler

import (
	"net/http"
	"net/url"

	"github.com/RATIU5/chewbacca/internal/view/components"
	"github.com/a-h/templ"
)

func ProcessAddrHandler(w http.ResponseWriter, r *http.Request) {
	addrStr := r.FormValue("addr")

	addr, err := url.Parse(addrStr)
	if err != nil || addr.Host == "" {
		templ.Handler(components.ErrResponseShow("That doesn't seem to be a valid URL")).ServeHTTP(w, r)
	} else {
		templ.Handler(components.TableResponseShow(addr.Host)).ServeHTTP(w, r)
	}

}
