package middleware

import (
	"github.com/labstack/echo/v4"
)

func ErrorConverter() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				return wrapToHttpError(err)
			}

			return nil
		}
	}
}
