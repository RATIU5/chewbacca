package handler

import (
	"github.com/RATIU5/chewbacca/internal/view/pages"
	"github.com/labstack/echo/v4"
)

type ResultsHandler struct{}

func (i ResultsHandler) ResultsAddr(c echo.Context) error {

	return render(c, pages.ShowResults())
}
