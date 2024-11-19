package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KnoblauchPilze/user-service/internal/service"
	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type mockAuthService struct {
	service.AuthService

	err error
}

func TestUnit_AuthController_WhenNoApiKeyProvided_ExpectBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	m := &mockAuthService{}
	expectedBody := []byte("\"Invalid API key\"\n")

	assertStatusCodeAndBody(t, req, m, http.StatusBadRequest, expectedBody)
}

func TestUnit_AuthController_WhenMultipleApiKeysProvided_ExpectBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("X-Api-Key", "e6349328-543b-4b4e-8a3c-4caf7b413589")
	req.Header.Add("X-Api-Key", "de2108c2-f87b-4033-825c-4ccbbb8b778e")

	m := &mockAuthService{}
	expectedBody := []byte("\"Invalid API key\"\n")

	assertStatusCodeAndBody(t, req, m, http.StatusBadRequest, expectedBody)
}

func TestUnit_AuthController_WhenApiKeyHasWrongSyntax_ExpectBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("X-Api-Key", "not-a-uuid")

	m := &mockAuthService{}
	expectedBody := []byte("\"Invalid API key\"\n")

	assertStatusCodeAndBody(t, req, m, http.StatusBadRequest, expectedBody)
}

func TestUnit_AuthController_WhenUserNotAuthenticated_ExpectForbidden(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("X-Api-Key", "e6349328-543b-4b4e-8a3c-4caf7b413589")

	m := &mockAuthService{
		err: errors.NewCode(service.UserNotAuthenticated),
	}
	expectedBody := `
	{
		"Code": 1000,
		"Message": "An unexpected error occurred"
	}`

	assertStatusCodeAndJsonBody(t, req, m, http.StatusForbidden, expectedBody)
}

func TestUnit_AuthController_WhenApiKeyIsExpired_ExpectForbidden(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("X-Api-Key", "e6349328-543b-4b4e-8a3c-4caf7b413589")

	m := &mockAuthService{
		err: errors.NewCode(service.AuthenticationExpired),
	}
	expectedBody := `
	{
		"Code": 1001,
		"Message": "An unexpected error occurred"
	}`

	assertStatusCodeAndJsonBody(t, req, m, http.StatusForbidden, expectedBody)
}

func TestUnit_AuthController(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("X-Api-Key", "e6349328-543b-4b4e-8a3c-4caf7b413589")

	m := &mockAuthService{}

	assertStatusCode(t, req, m, http.StatusNoContent)
}

func (m *mockAuthService) Authenticate(ctx context.Context, apiKey uuid.UUID) (communication.AuthorizationDtoResponse, error) {
	return communication.AuthorizationDtoResponse{}, m.err
}

func assertStatusCode(t *testing.T, req *http.Request, service service.AuthService, expectedStatusCode int) {
	ctx, rw := generateTestEchoContextFromRequest(req)

	err := authUser(ctx, service)

	require.Nil(t, err)
	require.Equal(t, expectedStatusCode, rw.Code)
}

func assertStatusCodeAndBody(t *testing.T, req *http.Request, service service.AuthService, expectedStatusCode int, expectedBody []byte) {
	ctx, rw := generateTestEchoContextFromRequest(req)

	err := authUser(ctx, service)

	require.Nil(t, err)
	require.Equal(t, expectedStatusCode, rw.Code)
	require.Equal(t, expectedBody, rw.Body.Bytes())
}

func assertStatusCodeAndJsonBody(t *testing.T, req *http.Request, service service.AuthService, expectedStatusCode int, expectedJsonBody string) {
	ctx, rw := generateTestEchoContextFromRequest(req)

	err := authUser(ctx, service)

	require.Nil(t, err)
	require.Equal(t, expectedStatusCode, rw.Code)
	require.JSONEq(t, expectedJsonBody, rw.Body.String())
}
