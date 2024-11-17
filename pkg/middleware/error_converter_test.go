package middleware

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUnit_ErrorConverter_CallsNextMiddleware(t *testing.T) {
	callable, called, ctx := createCallableHandler(ErrorConverter)

	err := callable(ctx)

	assert := assert.New(t)
	assert.Nil(err)
	assert.True(*called)
}

func TestUnit_ErrorConverter_WrapsUnknownErrorIntoHttpError(t *testing.T) {
	next := createErrorHandler(fmt.Errorf("some error"))
	middleware := ErrorConverter()
	callable := middleware(next)
	ctx, _ := generateTestEchoContext()

	err := callable(ctx)

	assertIsHttpErrorWithMessageAndCode(t, err, "some error", http.StatusInternalServerError)
}

func TestUnit_ErrorConverter_WrapsErrorWithCodeIntoHttpError(t *testing.T) {
	next := createErrorHandler(errors.NewCode(UncaughtPanic))
	middleware := ErrorConverter()
	callable := middleware(next)
	ctx, _ := generateTestEchoContext()

	err := callable(ctx)

	assertIsHttpErrorWithMessageAndCode(t, err, "(400) An unexpected error occurred", http.StatusInternalServerError)
}

func createErrorHandler(err error) echo.HandlerFunc {
	handler := func(c echo.Context) error {
		return err
	}

	return handler
}
