package server

import "github.com/labstack/echo/v4"

func New() {
	e := echo.New()

	e.GET("/limited", func(c echo.Context) error {
		return c.String(200, "Limited, don't over use me!")
	})
	e.GET("/unlimited", func(c echo.Context) error {
		return c.String(200, "Unlimited! Let's Go!")
	})
}
