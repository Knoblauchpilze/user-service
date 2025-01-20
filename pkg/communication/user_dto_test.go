package communication

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Knoblauchpilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var someTime = time.Date(2024, 11, 12, 19, 9, 36, 0, time.UTC)

func TestUnit_UserDtoRequest_MarshalsToCamelCase(t *testing.T) {
	dto := UserDtoRequest{
		Email:    "some@e.mail",
		Password: "secret",
	}

	out, err := json.Marshal(dto)

	assert.Nil(t, err)
	expectedJson := `
	{
		"email": "some@e.mail",
		"password": "secret"
	}`
	assert.JSONEq(t, expectedJson, string(out))
}

func TestUnit_FromUserDtoRequest(t *testing.T) {
	beforeConversion := time.Now()

	dto := UserDtoRequest{
		Email:    "email",
		Password: "password",
	}

	actual := FromUserDtoRequest(dto)

	assert.Equal(t, "email", actual.Email)
	assert.Equal(t, "password", actual.Password)
	assert.True(t, actual.CreatedAt.After(beforeConversion))
	assert.Equal(t, actual.CreatedAt, actual.UpdatedAt)
}

func TestUnit_UserDtoResponse_MarshalsToCamelCase(t *testing.T) {
	dto := UserDtoResponse{
		Id:        uuid.MustParse("a590b448-d3cd-4dbc-a9e3-8d642b1a5814"),
		Email:     "some@e.mail",
		Password:  "secret",
		CreatedAt: someTime,
	}

	out, err := json.Marshal(dto)

	assert.Nil(t, err)
	expectedJson := `
	{
		"id": "a590b448-d3cd-4dbc-a9e3-8d642b1a5814",
		"email": "some@e.mail",
		"password": "secret",
		"createdAt": "2024-11-12T19:09:36Z"
	}`
	assert.JSONEq(t, expectedJson, string(out))
}

func TestUnit_ToUserDtoResponse(t *testing.T) {
	entity := persistence.User{
		Id:       uuid.New(),
		Email:    "email",
		Password: "password",

		CreatedAt: someTime,
	}

	actual := ToUserDtoResponse(entity)

	assert.Equal(t, entity.Id, actual.Id)
	assert.Equal(t, "email", actual.Email)
	assert.Equal(t, "password", actual.Password)
	assert.Equal(t, someTime, actual.CreatedAt)
}
