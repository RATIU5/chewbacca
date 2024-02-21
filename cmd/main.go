package main

import (
	"github.com/RATIU5/chewbacca/handler"
	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()

	indexHandler := handler.IndexHandler{}
	app.GET("/", indexHandler.HandleIndexShow)

	resultsHandler := handler.ResultsHandler{}
	app.GET("/results", resultsHandler.HandleResultsShow)

	app.Start(":3000")	
}