package model

import "net/url"

type Route struct {
	RootAddr     *url.URL
	RootTitle    string
	CurrentAddr  *url.URL
	CurrentTitle string
	Status       int16
}

type Link struct {
	Link  string
	Title string
}
