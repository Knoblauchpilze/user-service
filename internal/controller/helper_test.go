package controller

import (
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
)

func generateTestEchoContextFromRequest(req *http.Request) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	rw := httptest.NewRecorder()

	ctx := e.NewContext(req, rw)
	return ctx, rw
}
