package handler

import (
	"github.com/RATIU5/chewbacca/internal/view/pages"
	"github.com/labstack/echo/v4"
)

type IndexHandler struct{}

func (i IndexHandler) HandleIndexShow(c echo.Context) error {

	addr := c.QueryParams().Get("addr")

	if (addr != "") {
		
	}
	return render(c, pages.ShowIndex())

}
