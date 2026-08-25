package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/db/postgresql"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/Knoblauchpilze/user-service/pkg/persistence"
	"github.com/Knoblauchpilze/user-service/pkg/repositories"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/require"
)

var (
	sampleApiKey = uuid.MustParse("e6349328-543b-4b4e-8a3c-4caf7b413589")

	dbTestConfig = postgresql.NewConfigForLocalhost("db_user_service", "user_service_manager", "manager_password")
)

func newTestConnection(t *testing.T) db.Connection {
	t.Helper()

	conn, err := db.New(t.Context(), dbTestConfig)
	require.Nil(t, err, "Actual err: %v", err)

	t.Cleanup(func() {
		conn.Close(t.Context())
	})

	return conn
}

func generateTestRequest(
	t *testing.T,
	method string,
	modifiers ...func(*testing.T, *http.Request),
) *http.Request {
	t.Helper()

	ctx := rest.WithContextLogger(t.Context(), slog.Default())
	req := httptest.NewRequestWithContext(ctx, method, "/", nil)

	for _, modifier := range modifiers {
		modifier(t, req)
	}

	return req
}

func generateTestRequestWithJsonBody[T any](
	t *testing.T,
	method string,
	data T,
) *http.Request {
	ctx := rest.WithContextLogger(t.Context(), slog.Default())
	req := httptest.NewRequestWithContext(ctx, method, "/", encodeBody(t, data))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func addSampleApiKeyHeader(t *testing.T, req *http.Request) {
	t.Helper()
	req.Header.Add("X-Api-Key", sampleApiKey.String())
}

func addApiKeyHeader(t *testing.T, req *http.Request, apiKey string) {
	t.Helper()
	req.Header.Add("X-Api-Key", apiKey)
}

func generateTestEchoContextFromRequest(req *http.Request) (*echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	rw := httptest.NewRecorder()

	ctx := e.NewContext(req, rw)
	return ctx, rw
}

func createTestGinRouter(
	t *testing.T,
	method string,
	path string,
	handler gin.HandlerFunc,
	middlewares ...gin.HandlerFunc,
) *gin.Engine {
	t.Helper()

	r := gin.New()

	for _, middleware := range middlewares {
		r.Use(middleware)
	}

	r.Handle(method, path, handler)

	return r
}

func encodeBody[T any](t *testing.T, data T) io.Reader {
	t.Helper()

	out, err := json.Marshal(data)
	require.NoError(t, err, "Actual err: %v", err)

	return bytes.NewReader(out)
}

func decodeResponseBody[T any](t *testing.T, w *httptest.ResponseRecorder) T {
	t.Helper()

	var responseBody T

	rawBody, err := io.ReadAll(w.Result().Body)
	require.NoError(t, err, "Actual err: %v", err)

	err = json.Unmarshal(rawBody, &responseBody)
	require.NoError(t, err, "Actual err: %v", err)

	return responseBody
}

type controllerFunc[Service any] func(*gin.Context, Service)

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

func insertTestUser(t *testing.T, conn db.Connection) persistence.User {
	repo := repositories.NewUserRepository(conn)

	id := uuid.New()
	user := persistence.User{
		Id:        id,
		Email:     fmt.Sprintf("my-user-%s", id),
		Password:  "my-password",
		CreatedAt: time.Now(),
	}
	out, err := repo.Create(context.Background(), user)
	require.Nil(t, err)

	assertUserExists(t, conn, out.Id)

	return out
}

func assertUserExists(t *testing.T, conn db.Connection, id uuid.UUID) {
	value, err := db.QueryOne[uuid.UUID](context.Background(), conn, "SELECT id FROM api_user WHERE id = $1", id)
	require.Nil(t, err)
	require.Equal(t, id, value)
}

func assertUserDoesNotExist(t *testing.T, conn db.Connection, id uuid.UUID) {
	value, err := db.QueryOne[int](context.Background(), conn, "SELECT COUNT(id) FROM api_user WHERE id = $1", id)
	require.Nil(t, err)
	require.Zero(t, value)
}

func assertEmailForUser(t *testing.T, conn db.Connection, user uuid.UUID, expectedEmail string) {
	value, err := db.QueryOne[string](context.Background(), conn, "SELECT email FROM api_user WHERE id = $1", user)
	require.Nil(t, err)
	require.Equal(t, expectedEmail, value)
}

func insertApiKeyForUser(t *testing.T, conn db.Connection, userId uuid.UUID) persistence.ApiKey {
	repo := repositories.NewApiKeyRepository(conn)

	apiKey := persistence.ApiKey{
		Id:         uuid.New(),
		Key:        uuid.New(),
		ApiUser:    userId,
		ValidUntil: time.Date(2024, 11, 22, 17, 00, 10, 0, time.UTC),
	}

	out, err := repo.Create(context.Background(), apiKey)
	require.Nil(t, err)

	assertApiKeyExists(t, conn, out.Id)

	return out
}

func assertApiKeyExists(t *testing.T, conn db.Connection, id uuid.UUID) {
	value, err := db.QueryOne[uuid.UUID](context.Background(), conn, "SELECT id FROM api_key WHERE id = $1", id)
	require.Nil(t, err)
	require.Equal(t, id, value)
}

func assertApiKeyExistsByKey(t *testing.T, conn db.Connection, key uuid.UUID) {
	value, err := db.QueryOne[uuid.UUID](context.Background(), conn, "SELECT key FROM api_key WHERE key = $1", key)
	require.Nil(t, err)
	require.Equal(t, key, value)
}

func assertApiKeyDoesNotExist(t *testing.T, conn db.Connection, id uuid.UUID) {
	value, err := db.QueryOne[int](context.Background(), conn, "SELECT COUNT(id) FROM api_key WHERE id = $1", id)
	require.Nil(t, err)
	require.Zero(t, value)
}
