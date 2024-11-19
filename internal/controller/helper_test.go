package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func generateTestEchoContextFromRequest(req *http.Request) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	rw := httptest.NewRecorder()

	ctx := e.NewContext(req, rw)
	return ctx, rw
}

type controllerFunc[Service any] func(echo.Context, Service) error

func assertStatusCode[Service any](t *testing.T, req *http.Request, service Service, callable controllerFunc[Service], expectedStatusCode int) {
	ctx, rw := generateTestEchoContextFromRequest(req)

	err := callable(ctx, service)

	require.Nil(t, err)
	require.Equal(t, expectedStatusCode, rw.Code)
}

func assertStatusCodeAndBody[Service any](t *testing.T, req *http.Request, service Service, callable controllerFunc[Service], expectedStatusCode int, expectedBody []byte) {
	ctx, rw := generateTestEchoContextFromRequest(req)

	err := callable(ctx, service)

	require.Nil(t, err)
	require.Equal(t, expectedStatusCode, rw.Code)
	require.Equal(t, expectedBody, rw.Body.Bytes(), "Actual: %s", rw.Body.String())
}

func assertStatusCodeAndJsonBody[Service any](t *testing.T, req *http.Request, service Service, callable controllerFunc[Service], expectedStatusCode int, expectedJsonBody string) {
	ctx, rw := generateTestEchoContextFromRequest(req)

	err := callable(ctx, service)

	require.Nil(t, err)
	require.Equal(t, expectedStatusCode, rw.Code)
	require.JSONEq(t, expectedJsonBody, rw.Body.String(), "Actual: %s", rw.Body.String())
}
