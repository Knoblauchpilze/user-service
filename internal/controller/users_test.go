package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	eassert "github.com/Knoblauchpilze/easy-assert/assert"
	"github.com/Knoblauchpilze/user-service/internal/service"
	"github.com/Knoblauchpilze/user-service/pkg/communication"
	"github.com/Knoblauchpilze/user-service/pkg/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockUserService struct {
	service.UserService
}

func TestUnit_UserController_CreateUser_WhenUserHasWrongSyntax_ExpectBadRequest(t *testing.T) {
	m := &mockUserService{}
	handler := createServiceAwareHttpHandler[service.UserService](createUser, m)

	r := createTestGinRouter(t, http.MethodPost, "/", handler)

	req := generateTestRequestWithJsonBody(t, http.MethodPost, "not-a-user-dto-request")
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusBadRequest, rw.Code)
	actual := decodeResponseBody[string](t, rw)
	assert.Equal(t, "Invalid user syntax", actual)
}

func TestIT_UserController_Create(t *testing.T) {
	m := &mockUserService{}
	handler := createServiceAwareHttpHandler[service.UserService](createUser, m)

	r := createTestGinRouter(t, http.MethodPost, "/", handler)

	requestDto := communication.UserDtoRequest{
		Email:    fmt.Sprintf("my-email-%s", uuid.NewString()),
		Password: "my-password",
	}
	req := generateTestRequestWithJsonBody(t, http.MethodPost, requestDto)
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusCreated, rw.Code)
	actual := decodeResponseBody[communication.UserDtoResponse](t, rw)
	assert.Equal(t, requestDto.Email, actual.Email)
	assert.Equal(t, requestDto.Password, actual.Password)
	conn := newTestConnection(t)
	assertUserExists(t, conn, actual.Id)
}

func TestIT_UserController_Create_WhenEmailIsEmpty_ExpectFailure(t *testing.T) {
	m := &mockUserService{}
	handler := createServiceAwareHttpHandler[service.UserService](createUser, m)

	r := createTestGinRouter(t, http.MethodPost, "/", handler)

	requestDto := communication.UserDtoRequest{
		Email:    "",
		Password: "my-password",
	}
	req := generateTestRequestWithJsonBody(t, http.MethodPost, requestDto)
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusBadRequest, rw.Code)
	actual := decodeResponseBody[string](t, rw)
	assert.Equal(t, "Invalid email", actual)
}

func TestIT_UserController_Create_WhenPasswordIsEmpty_ExpectFailure(t *testing.T) {
	m := &mockUserService{}
	handler := createServiceAwareHttpHandler[service.UserService](createUser, m)

	r := createTestGinRouter(t, http.MethodPost, "/", handler)

	requestDto := communication.UserDtoRequest{
		Email:    fmt.Sprintf("my-email-%s", uuid.NewString()),
		Password: "",
	}
	req := generateTestRequestWithJsonBody(t, http.MethodPost, requestDto)
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusBadRequest, rw.Code)
	actual := decodeResponseBody[string](t, rw)
	assert.Equal(t, "Invalid password", actual)
}

func TestIT_UserController_Create_WhenEmailAlreadyExists_ExpectFailure(t *testing.T) {
	service, conn := createTestUserService(t)
	user := insertTestUser(t, conn)

	handler := createServiceAwareHttpHandler(createUser, service)

	r := createTestGinRouter(t, http.MethodPost, "/", handler)

	requestDto := communication.UserDtoRequest{
		Email:    user.Email,
		Password: "my-super-password",
	}
	req := generateTestRequestWithJsonBody(t, http.MethodPost, requestDto)
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusConflict, rw.Code)
	actual := decodeResponseBody[string](t, rw)
	assert.Equal(t, "Email already in use", actual)
}

func TestUnit_UserController_GetUser_WhenIdHasWrongSyntax_ExpectBadRequest(t *testing.T) {
	m := &mockUserService{}
	handler := createServiceAwareHttpHandler[service.UserService](getUser, m)

	r := createTestGinRouter(t, http.MethodGet, "/:id", handler)

	req := generateTestRequest(t, http.MethodGet)
	addIdPathParam(t, req, "not-a-uuid")
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusBadRequest, rw.Code)
	actual := decodeResponseBody[string](t, rw)
	assert.Equal(t, "Invalid id syntax", actual)
}

func TestIT_UserController_GetUser(t *testing.T) {
	service, conn := createTestUserService(t)
	user := insertTestUser(t, conn)
	handler := createServiceAwareHttpHandler(getUser, service)

	r := createTestGinRouter(t, http.MethodGet, "/:id", handler)

	req := generateTestRequest(t, http.MethodGet)
	addIdPathParam(t, req, user.Id.String())
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusOK, rw.Code)
	actual := decodeResponseBody[communication.UserDtoResponse](t, rw)
	assert.Equal(t, user.Id, actual.Id)
	assert.Equal(t, user.Email, actual.Email)
	assert.Equal(t, user.Password, actual.Password)
	safetyMargin := 1 * time.Second
	assert.True(t, eassert.AreTimeCloserThan(user.CreatedAt, actual.CreatedAt, safetyMargin))
}

func TestIT_UserController_GetUser_WhenUserDoesNotExist_ExpectFailure(t *testing.T) {
	service, _ := createTestUserService(t)
	handler := createServiceAwareHttpHandler(getUser, service)

	r := createTestGinRouter(t, http.MethodGet, "/:id", handler)

	req := generateTestRequest(t, http.MethodGet)
	addIdPathParam(t, req, uuid.MustParse("00000000-1111-2222-1111-000000000000").String())
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusNotFound, rw.Code)
	actual := decodeResponseBody[string](t, rw)
	assert.Equal(t, "No such user", actual)
}

func TestIT_UserController_ListUsers(t *testing.T) {
	service, conn := createTestUserService(t)
	user1 := insertTestUser(t, conn)
	user2 := insertTestUser(t, conn)
	handler := createServiceAwareHttpHandler(listUsers, service)

	r := createTestGinRouter(t, http.MethodGet, "/", handler)

	req := generateTestRequest(t, http.MethodGet)
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusOK, rw.Code)
	actual := decodeResponseBody[[]uuid.UUID](t, rw)
	assert.GreaterOrEqual(t, len(actual), 2)
	assert.Contains(t, actual, user1.Id)
	assert.Contains(t, actual, user2.Id)
}

func TestUnit_UserController_UpdateUser_WhenIdHasWrongSyntax_ExpectBadRequest(t *testing.T) {
	m := &mockUserService{}
	handler := createServiceAwareHttpHandler[service.UserService](updateUser, m)

	r := createTestGinRouter(t, http.MethodPatch, "/:id", handler)

	req := generateTestRequest(t, http.MethodPatch)
	addIdPathParam(t, req, "not-a-uuid")
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusBadRequest, rw.Code)
	actual := decodeResponseBody[string](t, rw)
	assert.Equal(t, "Invalid id syntax", actual)
}

func TestUnit_UserController_UpdateUser_WhenUserHasWrongSyntax_ExpectBadRequest(t *testing.T) {
	service, conn := createTestUserService(t)
	user := insertTestUser(t, conn)
	handler := createServiceAwareHttpHandler(updateUser, service)

	r := createTestGinRouter(t, http.MethodPatch, "/:id", handler)

	req := generateTestRequestWithJsonBody(t, http.MethodPatch, "not-a-user-dto-request")
	addIdPathParam(t, req, user.Id.String())
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusBadRequest, rw.Code)
	actual := decodeResponseBody[string](t, rw)
	assert.Equal(t, "Invalid user syntax", actual)
}

func TestIT_UserController_UpdateUser(t *testing.T) {
	service, conn := createTestUserService(t)
	user := insertTestUser(t, conn)
	handler := createServiceAwareHttpHandler(updateUser, service)

	r := createTestGinRouter(t, http.MethodPatch, "/:id", handler)

	requestDto := communication.UserDtoRequest{
		Email:    fmt.Sprintf("my-other-email-%s", uuid.NewString()),
		Password: "my-other-password",
	}
	req := generateTestRequestWithJsonBody(t, http.MethodPatch, requestDto)
	addIdPathParam(t, req, user.Id.String())
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusOK, rw.Code)
	actual := decodeResponseBody[communication.UserDtoResponse](t, rw)
	assert.Equal(t, user.Email, actual.Email)
	assert.Equal(t, user.Password, actual.Password)
	assertEmailForUser(t, conn, user.Id, requestDto.Email)
}

func TestIT_UserController_UpdateUser_WhenUserDoesNotExist_ExpectFailure(t *testing.T) {
	service, conn := createTestUserService(t)
	handler := createServiceAwareHttpHandler(updateUser, service)

	r := createTestGinRouter(t, http.MethodPatch, "/:id", handler)

	id := uuid.MustParse("00000000-1111-2222-1111-000000000000")
	requestDto := communication.UserDtoRequest{
		Email:    fmt.Sprintf("my-email-%s", uuid.NewString()),
		Password: "my-new-password",
	}
	req := generateTestRequestWithJsonBody(t, http.MethodPatch, requestDto)
	addIdPathParam(t, req, id.String())
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusNotFound, rw.Code)
	actual := decodeResponseBody[string](t, rw)
	assert.Equal(t, "No such user", actual)
	assertUserDoesNotExist(t, conn, id)
}

func TestUnit_UserController_DeleteUser_WhenIdHasWrongSyntax_ExpectBadRequest(t *testing.T) {
	m := &mockUserService{}
	handler := createServiceAwareHttpHandler[service.UserService](deleteUser, m)

	r := createTestGinRouter(t, http.MethodDelete, "/:id", handler)

	req := generateTestRequest(t, http.MethodDelete)
	addIdPathParam(t, req, "not-a-uuid")
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusBadRequest, rw.Code)
	actual := decodeResponseBody[string](t, rw)
	assert.Equal(t, "Invalid id syntax", actual)
}

func TestIT_UserController_DeleteUser(t *testing.T) {
	service, conn := createTestUserService(t)
	user := insertTestUser(t, conn)

	handler := createServiceAwareHttpHandler(deleteUser, service)

	r := createTestGinRouter(t, http.MethodDelete, "/:id", handler)

	req := generateTestRequest(t, http.MethodDelete)
	addIdPathParam(t, req, user.Id.String())
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusNoContent, rw.Code)
	assertUserDoesNotExist(t, conn, user.Id)
}

func TestIT_UserController_DeleteUser_WhenLoggedIn_ExpectApiKeyAlsoDeleted(t *testing.T) {
	service, conn := createTestUserService(t)
	user := insertTestUser(t, conn)
	apiKey := insertApiKeyForUser(t, conn, user.Id)

	handler := createServiceAwareHttpHandler(deleteUser, service)

	r := createTestGinRouter(t, http.MethodDelete, "/:id", handler)

	req := generateTestRequest(t, http.MethodDelete)
	addIdPathParam(t, req, user.Id.String())
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusNoContent, rw.Code)
	assertApiKeyDoesNotExist(t, conn, apiKey.Id)
}

func TestIT_UserController_DeleteUser_WhenUserDoesNotExist_ExpectSuccess(t *testing.T) {
	service, conn := createTestUserService(t)
	handler := createServiceAwareHttpHandler(deleteUser, service)

	r := createTestGinRouter(t, http.MethodDelete, "/:id", handler)

	id := uuid.MustParse("00000000-1111-2222-1111-000000000000")
	req := generateTestRequest(t, http.MethodDelete)
	addIdPathParam(t, req, id.String())
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusNoContent, rw.Code)
	assertUserDoesNotExist(t, conn, id)
}

func TestUnit_UserController_LoginUserByEmail_WhenUserHasWrongSyntax_ExpectBadRequest(t *testing.T) {
	m := &mockUserService{}
	handler := createServiceAwareHttpHandler[service.UserService](loginUserByEmail, m)

	r := createTestGinRouter(t, http.MethodPatch, "/", handler)

	req := generateTestRequestWithJsonBody(t, http.MethodPatch, "not-a-user-dto-request")
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusBadRequest, rw.Code)
	actual := decodeResponseBody[string](t, rw)
	assert.Equal(t, "Invalid user syntax", actual)
}

// TODO: Add gin.SetMode(gin.Test??)
// TODO: To update
func TestIT_UserController_LoginUserByEmail(t *testing.T) {
	conn := newTestConnection(t)
	user := insertTestUser(t, conn)

	requestDto := communication.UserDtoRequest{
		Email:    user.Email,
		Password: user.Password,
	}

	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(requestDto)
	require.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/", &body)
	req.Header.Set("Content-Type", "application/json")
	ctx, rw := generateTestEchoContextFromRequest(req)

	service, _ := createTestUserService(t)

	err = loginUserByEmail(ctx, service)
	assert.Nil(t, err)

	var responseDto communication.ApiKeyDtoResponse
	err = json.Unmarshal(rw.Body.Bytes(), &responseDto)
	require.Nil(t, err)

	assert.Equal(t, http.StatusCreated, rw.Code)
	assertEmailForUser(t, conn, user.Id, requestDto.Email)
	assert.Equal(t, user.Id, responseDto.User)
	assertApiKeyExistsByKey(t, conn, responseDto.Key)
	expectedApproximateValidity := time.Now().Add(1 * time.Hour)
	safetyMargin := 5 * time.Second
	assert.True(t, eassert.AreTimeCloserThan(responseDto.ValidUntil, expectedApproximateValidity, safetyMargin))
}

func TestIT_UserController_LoginUserByEmail_WhenUserDoesNotExist_ExpectFailure(t *testing.T) {
	requestDto := communication.UserDtoRequest{
		Email:    fmt.Sprintf("some-email-%s", uuid.NewString()),
		Password: "my-password",
	}

	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(requestDto)
	require.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/", &body)
	req.Header.Set("Content-Type", "application/json")
	ctx, rw := generateTestEchoContextFromRequest(req)

	service, _ := createTestUserService(t)

	err = loginUserByEmail(ctx, service)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, rw.Code)
	assert.Equal(t, "\"No such user\"\n", rw.Body.String())
}

func TestIT_UserController_LoginUserByEmail_WhenPasswordDoesNotMatch_ExpectFailure(t *testing.T) {
	conn := newTestConnection(t)
	user := insertTestUser(t, conn)

	requestDto := communication.UserDtoRequest{
		Email:    user.Email,
		Password: fmt.Sprintf("%s-and-stuff", user.Password),
	}

	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(requestDto)
	require.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/", &body)
	req.Header.Set("Content-Type", "application/json")
	ctx, rw := generateTestEchoContextFromRequest(req)

	service, _ := createTestUserService(t)

	err = loginUserByEmail(ctx, service)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, rw.Code)
	assert.Equal(t, "\"Invalid credentials\"\n", rw.Body.String())
}

// TOD: End to update

func TestUnit_UserController_LogoutUser_WhenIdHasWrongSyntax_ExpectBadRequest(t *testing.T) {
	m := &mockUserService{}
	handler := createServiceAwareHttpHandler[service.UserService](logoutUser, m)

	r := createTestGinRouter(t, http.MethodDelete, "/:id", handler)

	req := generateTestRequest(t, http.MethodDelete)
	addIdPathParam(t, req, "not-a-uuid")
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusBadRequest, rw.Code)
	actual := decodeResponseBody[string](t, rw)
	assert.Equal(t, "Invalid id syntax", actual)
}

func TestIT_UserController_LogoutUser(t *testing.T) {
	service, conn := createTestUserService(t)
	user := insertTestUser(t, conn)
	apiKey := insertApiKeyForUser(t, conn, user.Id)

	handler := createServiceAwareHttpHandler(logoutUser, service)

	r := createTestGinRouter(t, http.MethodDelete, "/:id", handler)

	req := generateTestRequest(t, http.MethodDelete)
	addIdPathParam(t, req, user.Id.String())
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusNoContent, rw.Code)
	assertUserExists(t, conn, user.Id)
	assertApiKeyDoesNotExist(t, conn, apiKey.Id)
}

func TestIT_UserController_LogoutUser_WhenNotLoggedIn_ExpectSuccess(t *testing.T) {
	service, conn := createTestUserService(t)
	user := insertTestUser(t, conn)

	handler := createServiceAwareHttpHandler(logoutUser, service)

	r := createTestGinRouter(t, http.MethodDelete, "/:id", handler)

	req := generateTestRequest(t, http.MethodDelete)
	addIdPathParam(t, req, user.Id.String())
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusNoContent, rw.Code)
	assertUserExists(t, conn, user.Id)
}

func TestIT_UserController_LogoutUser_WhenUserDoesNotExist_ExpectFailure(t *testing.T) {
	service, _ := createTestUserService(t)
	handler := createServiceAwareHttpHandler(logoutUser, service)

	r := createTestGinRouter(t, http.MethodDelete, "/:id", handler)

	id := uuid.MustParse("00000000-1111-2222-1111-000000000000")
	req := generateTestRequest(t, http.MethodDelete)
	addIdPathParam(t, req, id.String())
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusNotFound, rw.Code)
	actual := decodeResponseBody[string](t, rw)
	assert.Equal(t, "No such user", actual)
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
