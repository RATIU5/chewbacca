package model

import "net/url"

type Route struct {
	RootAddr    *url.URL
	CurrentAddr *url.URL
	Status      int16
	Title       string
}

func NewRoute(rootAddr *url.URL, currentAddr *url.URL, status int16, title string) *Route {
	return &Route{
		RootAddr:    rootAddr,
		CurrentAddr: currentAddr,
		Status:      status,
		Title:       title,
	}
}
