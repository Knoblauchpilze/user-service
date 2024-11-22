package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/internal/service"
	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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

func TestIT_UserController_Create(t *testing.T) {
	requestDto := communication.UserDtoRequest{
		Email:    fmt.Sprintf("my-email-%s", uuid.NewString()),
		Password: "my-password",
	}

	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(requestDto)
	require.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/", &body)
	req.Header.Set("Content-Type", "application/json")
	ctx, rw := generateTestEchoContextFromRequest(req)

	service, conn := createTestUserService(t)

	err = createUser(ctx, service)
	assert.Nil(t, err)

	var responseDto communication.UserDtoResponse
	err = json.Unmarshal(rw.Body.Bytes(), &responseDto)
	require.Nil(t, err)

	assert.Equal(t, http.StatusCreated, rw.Code)
	assert.Equal(t, requestDto.Email, responseDto.Email)
	assert.Equal(t, requestDto.Password, responseDto.Password)
	assertUserExists(t, conn, responseDto.Id)
}

func TestIT_UserController_Create_WhenEmailIsEmpty_ExpectFailure(t *testing.T) {
	requestDto := communication.UserDtoRequest{
		Email:    "",
		Password: "my-password",
	}

	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(requestDto)
	require.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/", &body)
	req.Header.Set("Content-Type", "application/json")
	ctx, rw := generateTestEchoContextFromRequest(req)

	service, _ := createTestUserService(t)

	err = createUser(ctx, service)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, rw.Code)
	assert.Equal(t, "\"Invalid email\"\n", rw.Body.String())
}

func TestIT_UserController_Create_WhenPasswordIsEmpty_ExpectFailure(t *testing.T) {
	requestDto := communication.UserDtoRequest{
		Email:    fmt.Sprintf("my-email-%s", uuid.NewString()),
		Password: "",
	}

	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(requestDto)
	require.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/", &body)
	req.Header.Set("Content-Type", "application/json")
	ctx, rw := generateTestEchoContextFromRequest(req)

	service, _ := createTestUserService(t)

	err = createUser(ctx, service)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, rw.Code)
	assert.Equal(t, "\"Invalid password\"\n", rw.Body.String())
}

func TestIT_UserController_Create_WhenEmailAlreadyExists_ExpectFailure(t *testing.T) {
	conn := newTestConnection(t)
	user := insertTestUser(t, conn)

	requestDto := communication.UserDtoRequest{
		Email:    user.Email,
		Password: "my-super-password",
	}

	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(requestDto)
	require.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/", &body)
	req.Header.Set("Content-Type", "application/json")
	ctx, rw := generateTestEchoContextFromRequest(req)

	service, _ := createTestUserService(t)

	err = createUser(ctx, service)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusConflict, rw.Code)
	assert.Equal(t, "\"Email already in use\"\n", rw.Body.String())
}

func TestUnit_UserController_GetUser_WhenIdHasWrongSyntax_ExpectBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/not-a-uuid", nil)

	m := &mockUserService{}
	expectedBody := []byte("\"Invalid id syntax\"\n")

	assertStatusCodeAndBody[service.UserService](t, req, m, getUser, http.StatusBadRequest, expectedBody)
}

func TestIT_UserController_GetUser(t *testing.T) {
	conn := newTestConnection(t)
	user := insertTestUser(t, conn)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx, rw := generateTestEchoContextFromRequest(req)
	ctx.SetParamNames("id")
	ctx.SetParamValues(user.Id.String())

	service, _ := createTestUserService(t)

	err := getUser(ctx, service)
	assert.Nil(t, err)

	var responseDto communication.UserDtoResponse
	err = json.Unmarshal(rw.Body.Bytes(), &responseDto)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, rw.Code)
	assert.Equal(t, user.Id, responseDto.Id)
	assert.Equal(t, user.Email, responseDto.Email)
	assert.Equal(t, user.Password, responseDto.Password)
	safetyMargin := 1 * time.Second
	assert.True(t, areTimeCloserThan(user.CreatedAt, responseDto.CreatedAt, safetyMargin))
}

func TestIT_UserController_GetUser_WhenUserDoesNotExist_ExpectFailure(t *testing.T) {
	// Non-existent id
	id := uuid.MustParse("00000000-1111-2222-1111-000000000000")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx, rw := generateTestEchoContextFromRequest(req)
	ctx.SetParamNames("id")
	ctx.SetParamValues(id.String())

	service, _ := createTestUserService(t)

	err := getUser(ctx, service)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, rw.Code)
	assert.Equal(t, "\"No such user\"\n", rw.Body.String())
}

func TestIT_UserController_ListUsers(t *testing.T) {
	conn := newTestConnection(t)
	u1 := insertTestUser(t, conn)
	u2 := insertTestUser(t, conn)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx, rw := generateTestEchoContextFromRequest(req)

	service, _ := createTestUserService(t)

	err := listUsers(ctx, service)
	assert.Nil(t, err)

	var allUsers []uuid.UUID
	err = json.Unmarshal(rw.Body.Bytes(), &allUsers)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, rw.Code)
	assert.GreaterOrEqual(t, len(allUsers), 2)
	assert.Contains(t, allUsers, u1.Id)
	assert.Contains(t, allUsers, u2.Id)
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

func TestIT_UserController_DeleteUser(t *testing.T) {
	conn := newTestConnection(t)
	user := insertTestUser(t, conn)

	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	ctx, rw := generateTestEchoContextFromRequest(req)
	ctx.SetParamNames("id")
	ctx.SetParamValues(user.Id.String())

	service, _ := createTestUserService(t)

	err := deleteUser(ctx, service)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNoContent, rw.Code)
	assertUserDoesNotExist(t, conn, user.Id)
}

func TestIT_UserController_DeleteUser_WhenLoggedIn_ExpectApiKeyAlsoDeleted(t *testing.T) {
	conn := newTestConnection(t)
	user := insertTestUser(t, conn)
	apiKey := insertApiKeyForUser(t, conn, user.Id)

	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	ctx, rw := generateTestEchoContextFromRequest(req)
	ctx.SetParamNames("id")
	ctx.SetParamValues(user.Id.String())

	service, _ := createTestUserService(t)

	err := deleteUser(ctx, service)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNoContent, rw.Code)
	assertApiKeyDoesNotExist(t, conn, apiKey.Id)
}

func TestIT_UserController_DeleteUser_WhenUserDoesNotExist_ExpectSuccess(t *testing.T) {
	// Non-existent id
	id := uuid.MustParse("00000000-1111-2222-1111-000000000000")

	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	ctx, rw := generateTestEchoContextFromRequest(req)
	ctx.SetParamNames("id")
	ctx.SetParamValues(id.String())

	service, _ := createTestUserService(t)

	err := deleteUser(ctx, service)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNoContent, rw.Code)
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

func createTestUserService(t *testing.T) (service.UserService, db.Connection) {
	conn := newTestConnection(t)

	repos := repositories.Repositories{
		ApiKey: repositories.NewApiKeyRepository(conn),
		User:   repositories.NewUserRepository(conn),
	}

	config := service.ApiKeyConfig{
		Validity: 1 * time.Hour,
	}

	return service.NewUserService(config, conn, repos), conn
}

func areTimeCloserThan(t1 time.Time, t2 time.Time, distance time.Duration) bool {
	diff := t1.Sub(t2).Abs()
	return diff <= distance
}
