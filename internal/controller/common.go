package controller

import (
	"github.com/labstack/echo/v4"
)

type serviceAwareHttpHandler[T any] func(echo.Context, T) error

func createServiceAwareHttpHandler[T any](handler serviceAwareHttpHandler[T], service T) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handler(c, service)
	}
}
