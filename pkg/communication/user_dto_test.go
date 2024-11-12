package communication

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var someTime = time.Date(2024, 11, 12, 19, 9, 36, 0, time.UTC)

func TestUnit_UserDtoRequest_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	dto := UserDtoRequest{
		Email:    "some@e.mail",
		Password: "secret",
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	expectedJson := `
	{
		"email": "some@e.mail",
		"password": "secret"
	}`
	assert.JSONEq(expectedJson, string(out))
}

func TestUnit_FromUserDtoRequest(t *testing.T) {
	assert := assert.New(t)

	beforeConversion := time.Now()

	dto := UserDtoRequest{
		Email:    "email",
		Password: "password",
	}

	actual := FromUserDtoRequest(dto)

	assert.Equal("email", actual.Email)
	assert.Equal("password", actual.Password)
	assert.True(actual.CreatedAt.After(beforeConversion))
	assert.Equal(actual.CreatedAt, actual.UpdatedAt)
}

func TestUnit_UserDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	dto := UserDtoResponse{
		Id:        uuid.MustParse("a590b448-d3cd-4dbc-a9e3-8d642b1a5814"),
		Email:     "some@e.mail",
		Password:  "secret",
		CreatedAt: someTime,
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	expectedJson := `
	{
		"id": "a590b448-d3cd-4dbc-a9e3-8d642b1a5814",
		"email": "some@e.mail",
		"password": "secret",
		"createdAt": "2024-11-12T19:09:36Z"
	}`
	assert.JSONEq(expectedJson, string(out))
}

func TestUnit_ToUserDtoResponse(t *testing.T) {
	assert := assert.New(t)

	entity := persistence.User{
		Id:       uuid.New(),
		Email:    "email",
		Password: "password",

		CreatedAt: someTime,
	}

	actual := ToUserDtoResponse(entity)

	assert.Equal(entity.Id, actual.Id)
	assert.Equal("email", actual.Email)
	assert.Equal("password", actual.Password)
	assert.Equal(someTime, actual.CreatedAt)
}
