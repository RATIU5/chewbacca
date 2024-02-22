package handler

import (
	"net/url"

	"github.com/RATIU5/chewbacca/internal/view/components"
	"github.com/labstack/echo/v4"
)

type ProcessAddrHandler struct{}

func (i ProcessAddrHandler) HandleProcessAddr(c echo.Context) error {

	// Not sure if this is necessary
	if c.Request().Method != "POST" {
		return c.Redirect(302, "/")
	}

	addrStr := c.FormValue("addr")

	addr, err := url.Parse(addrStr)
	if err != nil || addr.Host == "" {
		return render(c, components.ErrResponseShow("That doesn't seem to be a valid URL"))
	}

	c.Response().Header().Set("HX-Redirect", "/results?addr="+addr.Host)

	return nil
}
