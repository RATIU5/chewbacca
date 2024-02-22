package main

import (
	"net/http"

	"github.com/RATIU5/chewbacca/internal/handler"
	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()

	// Serve static assets under /assets in the url
	app.Static("assets", "internal/assets/dist")

	// Serve the root index page
	indexHandler := handler.IndexHandler{}
	app.GET("/", indexHandler.HandleIndexShow)
	
	processAddrHandler := handler.ProcessAddrHandler{}
	app.POST("/process-addr", processAddrHandler.HandleProcessAddr)

	// Serve error pages
	app.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}

		if code == http.StatusNotFound {
			notFoundHandler := handler.NotFoundHandler{}
			notFoundHandler.HandleNotFoundShow(c)
		}
	}

	app.Start("localhost:3000")
}
