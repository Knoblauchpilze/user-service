package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/db/pgx"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUnit_UserService_Create(t *testing.T) {
	id := uuid.New()
	userDtoRequest := communication.UserDtoRequest{
		Email:    fmt.Sprintf("my-user-%s", id),
		Password: "my-password",
	}

	service, conn := newTestUserRepository(t)
	out, err := service.Create(context.Background(), userDtoRequest)

	assert := assert.New(t)
	assert.Nil(err)

	assert.Equal(userDtoRequest.Email, out.Email)
	assert.Equal(userDtoRequest.Password, out.Password)
	assertUserExists(t, conn, out.Id)
}

func TestUnit_UserService_Create_InvalidEmail(t *testing.T) {
	userDtoRequest := communication.UserDtoRequest{
		Email:    "",
		Password: "my-password",
	}

	service, _ := newTestUserRepository(t)
	_, err := service.Create(context.Background(), userDtoRequest)

	assert := assert.New(t)
	assert.True(errors.IsErrorWithCode(err, InvalidEmail), "Actual err: %v", err)
}

func TestUnit_UserService_Create_InvalidPassword(t *testing.T) {
	userDtoRequest := communication.UserDtoRequest{
		Email:    "my-username",
		Password: "",
	}

	service, _ := newTestUserRepository(t)
	_, err := service.Create(context.Background(), userDtoRequest)

	assert := assert.New(t)
	assert.True(errors.IsErrorWithCode(err, InvalidPassword), "Actual err: %v", err)
}

func TestUnit_UserService_Create_WhenUserAlreadyExists_ExpectFailure(t *testing.T) {
	service, conn := newTestUserRepository(t)
	user := insertTestUser(t, conn)
	userDtoRequest := communication.UserDtoRequest{
		Email:    user.Email,
		Password: "some-strong-password",
	}

	_, err := service.Create(context.Background(), userDtoRequest)

	assert := assert.New(t)
	assert.True(errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation), "Actual err: %v", err)
}

func TestUnit_UserService_Get(t *testing.T) {
	service, conn := newTestUserRepository(t)
	user := insertTestUser(t, conn)

	actual, err := service.Get(context.Background(), user.Id)

	assert := assert.New(t)
	assert.Nil(err)
	assert.Equal(user.Id, actual.Id)
	assert.Equal(user.Email, actual.Email)
	assert.Equal(user.Password, actual.Password)
}

func TestUnit_UserService_Get_WhenUserDoesNotExist_ExpectFailure(t *testing.T) {
	nonExistingId := uuid.MustParse("00000000-0000-1221-0000-000000000000")

	service, _ := newTestUserRepository(t)
	_, err := service.Get(context.Background(), nonExistingId)

	assert := assert.New(t)
	assert.True(errors.IsErrorWithCode(err, db.NoMatchingRows))
}

func TestUnit_UserService_List(t *testing.T) {
	service, conn := newTestUserRepository(t)
	user := insertTestUser(t, conn)

	ids, err := service.List(context.Background())

	assert := assert.New(t)
	assert.Nil(err)
	assert.Contains(ids, user.Id)
}

func TestUnit_UserService_Update(t *testing.T) {
	service, conn := newTestUserRepository(t)
	user := insertTestUser(t, conn)

	id := uuid.New()
	updatedUser := communication.UserDtoRequest{
		Email:    fmt.Sprintf("updated-email-%s", id),
		Password: "this-is-a-better-password",
	}

	updated, err := service.Update(context.Background(), user.Id, updatedUser)

	assert := assert.New(t)
	assert.Nil(err)
	assert.Equal(updatedUser.Email, updated.Email)
	assert.Equal(updatedUser.Password, updated.Password)

	actual, err := service.Get(context.Background(), user.Id)
	assert.Nil(err)
	assert.Equal(updatedUser.Email, actual.Email)
	assert.Equal(updatedUser.Password, actual.Password)
}

func TestUnit_UserService_Update_WhenUserDoesNotExist_ExpectFailure(t *testing.T) {
	nonExistentId := uuid.New()
	updatedUser := communication.UserDtoRequest{
		Email:    fmt.Sprintf("updated-email-%s", nonExistentId),
		Password: "this-is-a-better-password",
	}

	service, _ := newTestUserRepository(t)
	_, err := service.Update(context.Background(), nonExistentId, updatedUser)

	assert := assert.New(t)
	assert.True(errors.IsErrorWithCode(err, db.NoMatchingRows), "Actual err: %v", err)
}

func TestUnit_UserService_Update_WhenUpdateFails_ExpectFailure(t *testing.T) {
	service, conn := newTestUserRepository(t)
	user := insertTestUser(t, conn)
	otherUser := insertTestUser(t, conn)

	updatedUser := communication.UserDtoRequest{
		Email:    otherUser.Email,
		Password: "this-is-a-better-password",
	}

	_, err := service.Update(context.Background(), user.Id, updatedUser)

	assert := assert.New(t)
	assert.True(errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation), "Actual err: %v", err)
}

func TestUnit_UserService_Delete(t *testing.T) {
	service, conn := newTestUserRepository(t)
	user := insertTestUser(t, conn)

	err := service.Delete(context.Background(), user.Id)

	assert := assert.New(t)
	assert.Nil(err)
	assertUserDoesNotExist(t, conn, user.Id)
}

func TestUnit_UserService_Delete_WhenUserDoesNotExist_ExpectSuccess(t *testing.T) {
	nonExistingId := uuid.MustParse("00000000-0000-1221-0000-000000000000")

	service, _ := newTestUserRepository(t)
	err := service.Delete(context.Background(), nonExistingId)

	assert := assert.New(t)
	assert.Nil(err)
}

func TestUnit_UserService_Delete_WhenUserIsLoggedIn_ExpectApiKeyAlsoDeleted(t *testing.T) {
	service, conn := newTestUserRepository(t)
	user := insertTestUser(t, conn)
	apiKey := insertApiKeyForUser(t, conn, user.Id)

	err := service.Delete(context.Background(), user.Id)

	assert := assert.New(t)
	assert.Nil(err)
	assertApiKeyDoesNotExist(t, conn, apiKey.Id)
	assertUserDoesNotExist(t, conn, user.Id)
}

// Login
// LoginById
// Logout

func newTestUserRepository(t *testing.T) (UserService, db.Connection) {
	conn := newTestConnection(t)

	repos := repositories.Repositories{
		ApiKey: repositories.NewApiKeyRepository(conn),
		User:   repositories.NewUserRepository(conn),
	}

	apiKeyConfig := ApiKeyConfig{
		Validity: 1 * time.Hour,
	}

	return NewUserService(apiKeyConfig, conn, repos), conn
}
