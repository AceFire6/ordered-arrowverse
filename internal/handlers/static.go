package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func adsRoute(c echo.Context) error {
	return c.File("")
}

func privacyPolicyRoute(c echo.Context) error {
	return c.Render(http.StatusOK, "privacy-policy", nil)
}

func cookiePolicyRoute(c echo.Context) error {
	return c.Render(http.StatusOK, "cookie-policy", nil)
}
