package handler

import (
	"github.com/RATIU5/chewbacca/view/pages"
	"github.com/labstack/echo/v4"
)

type IndexHandler struct {}

func(i IndexHandler) HandleIndexShow(c echo.Context) error {
	return render(c, pages.ShowIndex())
} 