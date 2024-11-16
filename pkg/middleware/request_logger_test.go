package middleware

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnit_RequestLogger_CallsNextMiddleware(t *testing.T) {
	callable, called, ctx := createCallableHandler()

	err := callable(ctx)

	assert := assert.New(t)
	assert.Nil(err)
	assert.True(*called)
}

func TestUnit_RequestLogger_PrintsRequestTiming(t *testing.T) {
	callable, _, ctx := createCallableHandler()

	var out bytes.Buffer
	log := logger.New(&out)
	ctx.SetLogger(logger.Wrap(log))

	err := callable(ctx)
	afterCall := time.Now()

	assert := assert.New(t)
	assert.Nil(err)

	type message struct {
		Level   string
		Time    time.Time
		Message string
	}

	var actual message
	err = json.Unmarshal(out.Bytes(), &actual)
	require.Nil(t, err)

	assert.Equal("info", actual.Level)
	safetyMargin := 5 * time.Second
	assert.True(areTimeCloserThan(actual.Time, afterCall, safetyMargin), "%v and %v are not within %v", afterCall, actual.Time, safetyMargin)
	assert.Regexp(`GET example.com/ processed in [0-9]+.?s -> \x1b\[1;32m200\x1b\[0m`, actual.Message)
}

func createCallableHandler() (echo.HandlerFunc, *bool, echo.Context) {
	next, called := createTestEchoHandlerFuncWithCalledBoolean()
	ctx, _ := generateTestEchoContext()

	middlewareFunc := RequestLogger()
	callable := middlewareFunc(next)

	return callable, called, ctx
}

func areTimeCloserThan(t1 time.Time, t2 time.Time, distance time.Duration) bool {
	diff := t1.Sub(t2).Abs()
	return diff <= distance
}
