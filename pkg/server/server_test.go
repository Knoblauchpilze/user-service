package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/rest"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const reasonableTestTimeout = 5000 * time.Second
const reasonableTimeForServerToBeUpAndRunning = 100 * time.Millisecond

func TestUnit_Server_StopsWhenContextIsDone(t *testing.T) {
	s, ctx, cancel := createStoppableTestServer(context.Background())

	runServerAndExecuteHandler(t, ctx, s, cancel)
}

func TestUnit_Server_UnsupportedRoutes(t *testing.T) {
	s, _, _ := createStoppableTestServer(context.Background())

	handler := func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	}

	unsupportedMethods := []string{
		http.MethodHead,
		http.MethodPut,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}

	for _, method := range unsupportedMethods {
		t.Run(method, func(t *testing.T) {
			sampleRoute := rest.NewRoute(method, "/", handler)
			err := s.AddRoute(sampleRoute)
			assert.True(t, errors.IsErrorWithCode(err, UnsupportedMethod), "Actual err: %v", err)
		})
	}
}

func TestUnit_Server_ListensOnConfiguredPort(t *testing.T) {
	const port = 1234
	s, ctx, cancel := createStoppableTestServerWithPort(port, context.Background())

	var resp *http.Response
	var err error

	handler := func() {
		resp, err = http.Get(fmt.Sprintf("http://localhost:%d", port))
		cancel()
	}

	runServerAndExecuteHandler(t, ctx, s, handler)

	assert.Nil(t, err)
	assertResponseIs200Ok(t, resp)
}

func createStoppableTestServer(ctx context.Context) (Server, context.Context, context.CancelFunc) {
	return createStoppableTestServerWithPort(0, ctx)
}

func createStoppableTestServerWithPort(port uint16, ctx context.Context) (Server, context.Context, context.CancelFunc) {
	config := Config{
		Port:            port,
		ShutdownTimeout: 2 * time.Second,
	}

	cancellable, cancel := context.WithCancel(ctx)

	s := New(config)
	handler := func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	}
	sampleRoute := rest.NewRoute(http.MethodGet, "/", handler)
	s.AddRoute(sampleRoute)

	return s, cancellable, cancel
}

func runWithTimeout(handler func() error) (error, bool) {
	timer := time.After(reasonableTestTimeout)
	done := make(chan bool)

	var err error

	go func() {
		err = handler()
		done <- true
	}()

	select {
	case <-timer:
		return nil, true
	case <-done:
	}

	return err, false
}

func runServerWithTimeout(t *testing.T, ctx context.Context, s Server) {
	handler := func() error {
		return s.Start(ctx)
	}

	err, timeout := runWithTimeout(handler)

	require.False(t, timeout)
	require.Nil(t, err)
}

func runServerAndExecuteHandler(t *testing.T, ctx context.Context, s Server, handler func()) {
	go func() {
		time.Sleep(reasonableTimeForServerToBeUpAndRunning)
		handler()
	}()

	runServerWithTimeout(t, ctx, s)
}

func assertResponseIs200Ok(t *testing.T, resp *http.Response) {
	require.Equal(t, 200, resp.StatusCode)

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, "\"OK\"\n", string(data))
}
