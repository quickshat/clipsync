package main

import (
	"net/http"

	"strings"

	"github.com/labstack/echo"
)

func onlyLocal(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if strings.HasPrefix(c.Request().RemoteAddr, "127.0.0.1:") {
			next(c)
		}
		return c.NoContent(http.StatusUnauthorized)
	}
}
