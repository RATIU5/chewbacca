package model

type LinkInfo struct {
	URL    string
	Name   string
	Status int
	Type   string // "route", "image", "script", "stylesheet"
}
