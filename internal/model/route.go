package model

type Route struct {
	RootAddr    string
	CurrentAddr string
	Status      int16
	Title       string
}

func NewRoute(rootAddr string, currentAddr string, status int16, title string) *Route {
	return &Route{
		RootAddr:    rootAddr,
		CurrentAddr: currentAddr,
		Status:      status,
		Title:       title,
	}
}
