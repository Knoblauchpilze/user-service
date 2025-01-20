package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/user-service/internal/service"
	"github.com/Knoblauchpilze/user-service/pkg/communication"
	"github.com/google/uuid"
)

type mockAuthService struct {
	service.AuthService

	err error
}

func TestUnit_AuthController_WhenNoApiKeyProvided_ExpectBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	m := &mockAuthService{}
	expectedBody := []byte("\"Invalid API key\"\n")

	assertStatusCodeAndBody[service.AuthService](t, req, m, authUser, http.StatusBadRequest, expectedBody)
}

func TestUnit_AuthController_WhenMultipleApiKeysProvided_ExpectBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("X-Api-Key", "e6349328-543b-4b4e-8a3c-4caf7b413589")
	req.Header.Add("X-Api-Key", "de2108c2-f87b-4033-825c-4ccbbb8b778e")

	m := &mockAuthService{}
	expectedBody := []byte("\"Invalid API key\"\n")

	assertStatusCodeAndBody[service.AuthService](t, req, m, authUser, http.StatusBadRequest, expectedBody)
}

func TestUnit_AuthController_WhenApiKeyHasWrongSyntax_ExpectBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("X-Api-Key", "not-a-uuid")

	m := &mockAuthService{}
	expectedBody := []byte("\"Invalid API key\"\n")

	assertStatusCodeAndBody[service.AuthService](t, req, m, authUser, http.StatusBadRequest, expectedBody)
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

	assertStatusCodeAndJsonBody[service.AuthService](t, req, m, authUser, http.StatusForbidden, expectedBody)
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

	assertStatusCodeAndJsonBody[service.AuthService](t, req, m, authUser, http.StatusForbidden, expectedBody)
}

func TestUnit_AuthController(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("X-Api-Key", "e6349328-543b-4b4e-8a3c-4caf7b413589")

	m := &mockAuthService{}

	assertStatusCode[service.AuthService](t, req, m, authUser, http.StatusNoContent)
}

func (m *mockAuthService) Authenticate(ctx context.Context, apiKey uuid.UUID) (communication.AuthorizationDtoResponse, error) {
	return communication.AuthorizationDtoResponse{}, m.err
}
