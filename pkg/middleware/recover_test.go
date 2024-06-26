package middleware

import (
	"fmt"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRecover_CallsNextMiddleware(t *testing.T) {
	assert := assert.New(t)
	ctx, _ := generateTestEchoContext()
	next, called := createHandlerFuncWithCalledBoolean()

	em := Recover()
	callable := em(next)
	callable(ctx)

	assert.True(*called)
}

func TestRecover_WhenNoErrorReturnsNoError(t *testing.T) {
	assert := assert.New(t)
	ctx, _ := generateTestEchoContext()
	next := createHandlerFuncReturning(nil)

	em := Recover()
	callable := em(next)
	actual := callable(ctx)

	assert.Nil(actual)
}

func TestRecover_WhenNoErrorDoesNotCallContextError(t *testing.T) {
	assert := assert.New(t)
	errHandler, called, _ := createErrorHandlerFunc()
	ctx, _ := generateTestEchoContextWithErrorHandler(errHandler)
	next := createHandlerFuncReturning(nil)

	em := Recover()
	callable := em(next)
	callable(ctx)

	assert.False(*called)
}

func TestRecover_PropagatesError(t *testing.T) {
	assert := assert.New(t)
	ctx, _ := generateTestEchoContext()
	next := createHandlerFuncReturning(errDefault)

	em := Recover()
	callable := em(next)
	actual := callable(ctx)

	assert.Equal(actual, errDefault)
}

func createPanickingHandlerFunc(err interface{}) echo.HandlerFunc {
	return func(c echo.Context) error {
		panic(err)
	}
}

func TestRecover_SetsContextErrorOnPanic(t *testing.T) {
	assert := assert.New(t)
	errHandler, called, reportedErr := createErrorHandlerFunc()
	ctx, _ := generateTestEchoContextWithErrorHandler(errHandler)
	next := createPanickingHandlerFunc(errDefault)

	em := Recover()
	callable := em(next)
	callable(ctx)

	assert.True(*called)
	assert.Equal(errDefault, *reportedErr)
}

func TestRecover_ConvertsToErrorUnknownPanic(t *testing.T) {
	assert := assert.New(t)
	errHandler, called, reportedErr := createErrorHandlerFunc()
	ctx, _ := generateTestEchoContextWithErrorHandler(errHandler)
	next := createPanickingHandlerFunc(36)

	em := Recover()
	callable := em(next)
	callable(ctx)

	expected := fmt.Errorf("%v", 36)
	assert.True(*called)
	assert.Equal(expected, *reportedErr)
}
