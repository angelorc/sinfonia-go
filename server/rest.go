package server

import (
	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func InitRest(e *echo.Echo) {
	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "<3")
	})
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
}
