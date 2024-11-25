package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/db/pgx"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestIT_UserService_Create(t *testing.T) {
	id := uuid.New()
	userDtoRequest := communication.UserDtoRequest{
		Email:    fmt.Sprintf("my-user-%s", id),
		Password: "my-password",
	}

	service, conn := newTestUserRepository(t)
	out, err := service.Create(context.Background(), userDtoRequest)

	assert.Nil(t, err)

	assert.Equal(t, userDtoRequest.Email, out.Email)
	assert.Equal(t, userDtoRequest.Password, out.Password)
	assertUserExists(t, conn, out.Id)
}

func TestIT_UserService_Create_InvalidEmail(t *testing.T) {
	userDtoRequest := communication.UserDtoRequest{
		Email:    "",
		Password: "my-password",
	}

	service, _ := newTestUserRepository(t)
	_, err := service.Create(context.Background(), userDtoRequest)

	assert.True(t, errors.IsErrorWithCode(err, InvalidEmail), "Actual err: %v", err)
}

func TestIT_UserService_Create_InvalidPassword(t *testing.T) {
	userDtoRequest := communication.UserDtoRequest{
		Email:    "my-username",
		Password: "",
	}

	service, _ := newTestUserRepository(t)
	_, err := service.Create(context.Background(), userDtoRequest)

	assert.True(t, errors.IsErrorWithCode(err, InvalidPassword), "Actual err: %v", err)
}

func TestIT_UserService_Create_WhenUserAlreadyExists_ExpectFailure(t *testing.T) {
	service, conn := newTestUserRepository(t)
	user := insertTestUser(t, conn)
	userDtoRequest := communication.UserDtoRequest{
		Email:    user.Email,
		Password: "some-strong-password",
	}

	_, err := service.Create(context.Background(), userDtoRequest)

	assert.True(t, errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation), "Actual err: %v", err)
}

func TestIT_UserService_Get(t *testing.T) {
	service, conn := newTestUserRepository(t)
	user := insertTestUser(t, conn)

	actual, err := service.Get(context.Background(), user.Id)

	assert.Nil(t, err)
	assert.Equal(t, user.Id, actual.Id)
	assert.Equal(t, user.Email, actual.Email)
	assert.Equal(t, user.Password, actual.Password)
}

func TestIT_UserService_Get_WhenUserDoesNotExist_ExpectFailure(t *testing.T) {
	nonExistingId := uuid.MustParse("00000000-0000-1221-0000-000000000000")

	service, _ := newTestUserRepository(t)
	_, err := service.Get(context.Background(), nonExistingId)

	assert.True(t, errors.IsErrorWithCode(err, db.NoMatchingRows), "Actual err: %v", err)
}

func TestIT_UserService_List(t *testing.T) {
	service, conn := newTestUserRepository(t)
	user := insertTestUser(t, conn)

	ids, err := service.List(context.Background())

	assert.Nil(t, err)
	assert.Contains(t, ids, user.Id)
}

func TestIT_UserService_Update(t *testing.T) {
	service, conn := newTestUserRepository(t)
	user := insertTestUser(t, conn)

	id := uuid.New()
	updatedUser := communication.UserDtoRequest{
		Email:    fmt.Sprintf("updated-email-%s", id),
		Password: "this-is-a-better-password",
	}

	updated, err := service.Update(context.Background(), user.Id, updatedUser)

	assert.Nil(t, err)
	assert.Equal(t, updatedUser.Email, updated.Email)
	assert.Equal(t, updatedUser.Password, updated.Password)

	actual, err := service.Get(context.Background(), user.Id)
	assert.Nil(t, err)
	assert.Equal(t, updatedUser.Email, actual.Email)
	assert.Equal(t, updatedUser.Password, actual.Password)
}

func TestIT_UserService_Update_WhenUserDoesNotExist_ExpectFailure(t *testing.T) {
	nonExistentId := uuid.New()
	updatedUser := communication.UserDtoRequest{
		Email:    fmt.Sprintf("updated-email-%s", nonExistentId),
		Password: "this-is-a-better-password",
	}

	service, _ := newTestUserRepository(t)
	_, err := service.Update(context.Background(), nonExistentId, updatedUser)

	assert.True(t, errors.IsErrorWithCode(err, db.NoMatchingRows), "Actual err: %v", err)
}

func TestIT_UserService_Update_WhenUpdateFails_ExpectFailure(t *testing.T) {
	service, conn := newTestUserRepository(t)
	user := insertTestUser(t, conn)
	otherUser := insertTestUser(t, conn)

	updatedUser := communication.UserDtoRequest{
		Email:    otherUser.Email,
		Password: "this-is-a-better-password",
	}

	_, err := service.Update(context.Background(), user.Id, updatedUser)

	assert.True(t, errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation), "Actual err: %v", err)
}

func TestIT_UserService_Delete(t *testing.T) {
	service, conn := newTestUserRepository(t)
	user := insertTestUser(t, conn)

	err := service.Delete(context.Background(), user.Id)

	assert.Nil(t, err)
	assertUserDoesNotExist(t, conn, user.Id)
}

func TestIT_UserService_Delete_WhenUserDoesNotExist_ExpectSuccess(t *testing.T) {
	nonExistingId := uuid.MustParse("00000000-0000-1221-0000-000000000000")

	service, _ := newTestUserRepository(t)
	err := service.Delete(context.Background(), nonExistingId)

	assert.Nil(t, err)
}

func TestIT_UserService_Delete_WhenUserIsLoggedIn_ExpectApiKeyAlsoDeleted(t *testing.T) {
	service, conn := newTestUserRepository(t)
	user := insertTestUser(t, conn)
	apiKey := insertApiKeyForUser(t, conn, user.Id)

	err := service.Delete(context.Background(), user.Id)

	assert.Nil(t, err)
	assertApiKeyDoesNotExist(t, conn, apiKey.Id)
	assertUserDoesNotExist(t, conn, user.Id)
}

func TestIT_UserService_Login_ExpectCorrectUserAndValidity(t *testing.T) {
	service, conn := newTestUserRepository(t)
	user := insertTestUser(t, conn)

	userDtoRequest := communication.UserDtoRequest{
		Email:    user.Email,
		Password: user.Password,
	}

	apiKey, err := service.Login(context.Background(), userDtoRequest)

	assert.Nil(t, err)
	assert.Equal(t, user.Id, apiKey.User)
	validityDateWithSafetyMargin := time.Now().Add(55 * time.Minute)
	assert.True(t, apiKey.ValidUntil.After(validityDateWithSafetyMargin))
	assertApiKeyExistsByKey(t, conn, apiKey.Key)
}

func TestIT_UserService_Login_WhenUserDoesNotExist_ExpectFailure(t *testing.T) {
	userDtoRequest := communication.UserDtoRequest{
		Email:    fmt.Sprintf("not-an-existing-email-%s", uuid.New()),
		Password: "my-password",
	}

	service, _ := newTestUserRepository(t)
	_, err := service.Login(context.Background(), userDtoRequest)

	assert.True(t, errors.IsErrorWithCode(err, db.NoMatchingRows), "Actual err: %v", err)
}

func TestIT_UserService_Login_WhenCredentialsAreWrong_ExpectFailure(t *testing.T) {
	service, conn := newTestUserRepository(t)
	user := insertTestUser(t, conn)

	userDtoRequest := communication.UserDtoRequest{
		Email:    user.Email,
		Password: "not-the-right-password",
	}

	_, err := service.Login(context.Background(), userDtoRequest)

	assert.True(t, errors.IsErrorWithCode(err, InvalidCredentials), "Actual err: %v", err)
}

func TestIT_UserService_Login_WhenUserAlreadyLoggedIn_ExpectApiKeyValidityIsExtended(t *testing.T) {
	service, conn := newTestUserRepository(t)
	user := insertTestUser(t, conn)
	timeInThePast := time.Now().Add(-1 * time.Hour)
	apiKey := insertApiKeyForUserWithValidity(t, conn, user.Id, timeInThePast)

	userDtoRequest := communication.UserDtoRequest{
		Email:    user.Email,
		Password: user.Password,
	}

	updatedApiKey, err := service.Login(context.Background(), userDtoRequest)

	assert.Nil(t, err)
	assert.Equal(t, apiKey.Key, updatedApiKey.Key)
	assert.Equal(t, user.Id, updatedApiKey.User)
	assert.True(t, timeInThePast.Before(updatedApiKey.ValidUntil))
	validityDateWithSafetyMargin := time.Now().Add(55 * time.Minute)
	assert.True(t, updatedApiKey.ValidUntil.After(validityDateWithSafetyMargin))
}

func TestIT_UserService_Logout(t *testing.T) {
	service, conn := newTestUserRepository(t)
	user := insertTestUser(t, conn)
	apiKey := insertApiKeyForUser(t, conn, user.Id)

	err := service.Logout(context.Background(), user.Id)

	assert.Nil(t, err)
	assertApiKeyDoesNotExist(t, conn, apiKey.Id)
	assertUserExists(t, conn, user.Id)
}

func TestIT_UserService_Logout_WhenUserDoesNotExist_ExpectFailure(t *testing.T) {
	nonExistingId := uuid.MustParse("00000000-0000-1221-0000-000000000000")

	service, _ := newTestUserRepository(t)
	err := service.Logout(context.Background(), nonExistingId)

	assert.True(t, errors.IsErrorWithCode(err, db.NoMatchingRows), "Actual err: %v", err)
}

func TestIT_UserService_Logout_WhenNotLoggedIn_ExpectSuccess(t *testing.T) {
	service, conn := newTestUserRepository(t)
	user := insertTestUser(t, conn)

	err := service.Logout(context.Background(), user.Id)

	assert.Nil(t, err)
	assertUserExists(t, conn, user.Id)
}

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
