package controller

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/KnoblauchPilze/user-service/internal/service"
	"github.com/stretchr/testify/require"
)

type mockUserService struct {
	service.UserService
}

func TestUnit_UserController_CreateUser_WhenUserHasWrongSyntax_ExpectBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("not-a-user-dto-request"))

	m := &mockUserService{}
	expectedBody := []byte("\"Invalid user syntax\"\n")

	assertStatusCodeAndBody[service.UserService](t, req, m, createUser, http.StatusBadRequest, expectedBody)
}

func TestUnit_UserController_GetUser_WhenIdHasWrongSyntax_ExpectBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/not-a-uuid", nil)

	m := &mockUserService{}
	expectedBody := []byte("\"Invalid id syntax\"\n")

	assertStatusCodeAndBody[service.UserService](t, req, m, getUser, http.StatusBadRequest, expectedBody)
}

func TestUnit_UserController_UpdateUser_WhenIdHasWrongSyntax_ExpectBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPatch, "/not-a-uuid", nil)

	m := &mockUserService{}
	expectedBody := []byte("\"Invalid id syntax\"\n")

	assertStatusCodeAndBody[service.UserService](t, req, m, updateUser, http.StatusBadRequest, expectedBody)
}

func TestUnit_UserController_UpdateUser_WhenUserHasWrongSyntax_ExpectBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader("not-a-user-dto-request"))

	ctx, rw := generateTestEchoContextFromRequest(req)
	ctx.SetParamNames("id")
	ctx.SetParamValues("e6349328-543b-4b4e-8a3c-4caf7b413589")

	m := &mockUserService{}
	err := updateUser(ctx, m)

	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, rw.Code)
	expectedBody := []byte("\"Invalid user syntax\"\n")
	require.Equal(t, expectedBody, rw.Body.Bytes(), "Actual: %s", rw.Body.String())
}

func TestUnit_UserController_DeleteUser_WhenIdHasWrongSyntax_ExpectBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/not-a-uuid", nil)

	m := &mockUserService{}
	expectedBody := []byte("\"Invalid id syntax\"\n")

	assertStatusCodeAndBody[service.UserService](t, req, m, deleteUser, http.StatusBadRequest, expectedBody)
}

func TestUnit_UserController_LoginUserByEmail_WhenUserHasWrongSyntax_ExpectBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("not-a-user-dto-request"))

	m := &mockUserService{}
	expectedBody := []byte("\"Invalid user syntax\"\n")

	assertStatusCodeAndBody[service.UserService](t, req, m, loginUserByEmail, http.StatusBadRequest, expectedBody)
}

func TestUnit_UserController_LogoutUser_WhenIdHasWrongSyntax_ExpectBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/not-a-uuid", nil)

	m := &mockUserService{}
	expectedBody := []byte("\"Invalid id syntax\"\n")

	assertStatusCodeAndBody[service.UserService](t, req, m, logoutUser, http.StatusBadRequest, expectedBody)
}
