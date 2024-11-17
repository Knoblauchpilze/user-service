package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/rest"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const reasonableTestTimeout = 5000 * time.Second
const reasonableTimeForServerToBeUpAndRunning = 100 * time.Millisecond

type responseEnvelope struct {
	RequestId string          `json:"requestId"`
	Status    string          `json:"status"`
	Details   json.RawMessage `json:"details"`
}

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
	assertResponseStatusMatches(t, resp, http.StatusOK)
	actual := unmarshalResponseAndAssertRequestId(t, resp)
	assert.Equal(t, "SUCCESS", actual.Status)
	assert.Equal(t, `"OK"`, string(actual.Details))
}

func TestUnit_Server_WrapsResponseInEnvelope(t *testing.T) {
	const port = 1235
	s, ctx, cancel := createStoppableTestServerWithPort(port, context.Background())

	var resp *http.Response
	var err error

	handler := func() {
		resp, err = http.Get(fmt.Sprintf("http://localhost:%d", port))
		cancel()
	}

	runServerAndExecuteHandler(t, ctx, s, handler)

	assert.Nil(t, err)
	assertResponseStatusMatches(t, resp, http.StatusOK)
	actual := unmarshalResponseAndAssertRequestId(t, resp)
	assert.Equal(t, "SUCCESS", actual.Status)
	assert.Equal(t, `"OK"`, string(actual.Details))
}

func TestUnit_Server_WhenHandlerPanics_ExpectErrorResponseEnvelope(t *testing.T) {
	const port = 1236
	route := func(c echo.Context) error {
		panic(fmt.Errorf("this handler panics"))
	}
	s, ctx, cancel := createStoppableTestServerWithPortAndHandler(port, context.Background(), route)

	var resp *http.Response
	var err error

	handler := func() {
		resp, err = http.Get(fmt.Sprintf("http://localhost:%d", port))
		cancel()
	}

	runServerAndExecuteHandler(t, ctx, s, handler)

	assert.Nil(t, err)
	assertResponseStatusMatches(t, resp, http.StatusInternalServerError)
	actual := unmarshalResponseAndAssertRequestId(t, resp)
	assert.Equal(t, "ERROR", actual.Status)
	assert.Equal(t, `{"message":"this handler panics"}`, string(actual.Details))
}

func TestUnit_Server_WhenHandlerReturnsError_ExpectErrorResponseEnvelope(t *testing.T) {
	const port = 1237
	route := func(c echo.Context) error {
		return errors.NewCode(db.AlreadyCommitted)
	}
	s, ctx, cancel := createStoppableTestServerWithPortAndHandler(port, context.Background(), route)

	var resp *http.Response
	var err error

	handler := func() {
		resp, err = http.Get(fmt.Sprintf("http://localhost:%d", port))
		cancel()
	}

	runServerAndExecuteHandler(t, ctx, s, handler)

	assert.Nil(t, err)
	assertResponseStatusMatches(t, resp, http.StatusInternalServerError)
	actual := unmarshalResponseAndAssertRequestId(t, resp)
	assert.Equal(t, "ERROR", actual.Status)
	assert.Equal(t, `{"message":"An unexpected error occurred. Code: 102"}`, string(actual.Details))
}

func createStoppableTestServer(ctx context.Context) (Server, context.Context, context.CancelFunc) {
	return createStoppableTestServerWithPort(0, ctx)
}

func createStoppableTestServerWithPort(port uint16, ctx context.Context) (Server, context.Context, context.CancelFunc) {
	handler := func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	}

	return createStoppableTestServerWithPortAndHandler(port, ctx, handler)
}

func createStoppableTestServerWithPortAndHandler(port uint16, ctx context.Context, handler echo.HandlerFunc) (Server, context.Context, context.CancelFunc) {
	config := Config{
		Port:            port,
		ShutdownTimeout: 2 * time.Second,
	}

	cancellable, cancel := context.WithCancel(ctx)

	s := New(config)
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

func unmarshalResponseAndAssertRequestId(t *testing.T, resp *http.Response) responseEnvelope {
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var out responseEnvelope
	err = json.Unmarshal(data, &out)
	require.Nil(t, err)

	require.Regexp(t, `[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`, out.RequestId)

	return out
}

func assertResponseStatusMatches(t *testing.T, resp *http.Response, httpCode int) {
	require.Equal(t, httpCode, resp.StatusCode)
}
