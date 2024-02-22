package handler

import (
	"github.com/RATIU5/chewbacca/internal/view/pages"
	"github.com/labstack/echo/v4"
)

type ProcessAddrHandler struct{}

func (i ProcessAddrHandler) HandleProcessAddr(c echo.Context) error {
	return render(c, pages.ShowNotFound())
}
