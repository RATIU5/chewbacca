package model

import "net/url"

type Route struct {
	RootAddr    url.URL
	ProcessAddr url.URL
	Status      int16
	Title       string
}

func NewRoute(rootAddr, processAddr url.URL, status int16, title string) *Route {
	return &Route{
		RootAddr:    rootAddr,
		ProcessAddr: processAddr,
		Status:      status,
		Title:       title,
	}
}
