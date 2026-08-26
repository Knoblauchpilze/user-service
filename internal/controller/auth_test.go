package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/user-service/internal/service"
	"github.com/Knoblauchpilze/user-service/pkg/communication"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockAuthService struct {
	service.AuthService

	err error
}

func TestUnit_AuthController_WhenNoApiKeyProvided_ExpectBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	m := &mockAuthService{}
	handler := createServiceAwareHttpHandler[service.AuthService](authUser, m)

	r := createTestGinRouter(t, http.MethodGet, "/", handler)

	req := generateTestRequest(t, http.MethodGet)
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusBadRequest, rw.Code)
	actual := decodeResponseBody[string](t, rw)
	assert.Equal(t, "Invalid API key", actual)
}

func TestUnit_AuthController_WhenMultipleApiKeysProvided_ExpectBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	m := &mockAuthService{}
	handler := createServiceAwareHttpHandler[service.AuthService](authUser, m)

	r := createTestGinRouter(t, http.MethodGet, "/", handler)

	req := generateTestRequest(t, http.MethodGet, addSampleApiKeyHeader)
	addApiKeyHeader(t, req, uuid.NewString())
	addApiKeyHeader(t, req, uuid.NewString())
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusBadRequest, rw.Code)
	actual := decodeResponseBody[string](t, rw)
	assert.Equal(t, "Invalid API key", actual)
}

func TestUnit_AuthController_WhenApiKeyHasWrongSyntax_ExpectBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	m := &mockAuthService{}
	handler := createServiceAwareHttpHandler[service.AuthService](authUser, m)

	r := createTestGinRouter(t, http.MethodGet, "/", handler)

	req := generateTestRequest(t, http.MethodGet, addSampleApiKeyHeader)
	addApiKeyHeader(t, req, "not-a-uuid")
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusBadRequest, rw.Code)
	actual := decodeResponseBody[string](t, rw)
	assert.Equal(t, "Invalid API key", actual)
}

func TestUnit_AuthController_WhenUserNotAuthenticated_ExpectForbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)

	m := &mockAuthService{
		err: service.ErrUserNotAuthenticated,
	}
	handler := createServiceAwareHttpHandler[service.AuthService](authUser, m)

	r := createTestGinRouter(t, http.MethodGet, "/", handler)

	req := generateTestRequest(t, http.MethodGet, addSampleApiKeyHeader)
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusForbidden, rw.Code)
	actual := decodeResponseBody[errors.ErrorWithCode](t, rw)
	assert.Equal(t, errors.ErrorCode(1000), actual.Code)
	assert.Equal(t, "an unexpected error occurred", actual.Message)
}

func TestUnit_AuthController_WhenApiKeyIsExpired_ExpectForbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)

	m := &mockAuthService{
		err: service.ErrAuthenticationExpired,
	}
	handler := createServiceAwareHttpHandler[service.AuthService](authUser, m)

	r := createTestGinRouter(t, http.MethodGet, "/", handler)

	req := generateTestRequest(t, http.MethodGet, addSampleApiKeyHeader)
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusForbidden, rw.Code)
	actual := decodeResponseBody[errors.ErrorWithCode](t, rw)
	assert.Equal(t, errors.ErrorCode(1001), actual.Code)
	assert.Equal(t, "an unexpected error occurred", actual.Message)
}

func TestUnit_AuthController(t *testing.T) {
	gin.SetMode(gin.TestMode)

	m := &mockAuthService{}
	handler := createServiceAwareHttpHandler[service.AuthService](authUser, m)

	r := createTestGinRouter(t, http.MethodGet, "/", handler)

	req := generateTestRequest(t, http.MethodGet, addSampleApiKeyHeader)
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusNoContent, rw.Code)
}

func (m *mockAuthService) Authenticate(ctx context.Context, apiKey uuid.UUID) (communication.AuthorizationDtoResponse, error) {
	return communication.AuthorizationDtoResponse{}, m.err
}
