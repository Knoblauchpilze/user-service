package middleware

import (
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
)

func createTestEchoHandlerFuncWithCalledBoolean() (echo.HandlerFunc, *bool) {
	called := false
	call := func(c echo.Context) error {
		called = true
		return c.NoContent(http.StatusOK)
	}
	return call, &called
}

func generateTestEchoContext() (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	return generateTestEchoContextFromRequest(req)
}

func generateTestEchoContextFromRequest(req *http.Request) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	rw := httptest.NewRecorder()

	ctx := e.NewContext(req, rw)

	return ctx, rw
}
