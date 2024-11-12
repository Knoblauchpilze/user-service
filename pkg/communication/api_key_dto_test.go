package communication

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUnit_ApiKeyDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	dto := ApiKeyDtoResponse{
		User:       uuid.MustParse("c74a22da-8a05-43a9-a8b9-717e422b0af4"),
		Key:        uuid.MustParse("872e9e40-ce61-497e-b606-c7a08a4faa14"),
		ValidUntil: someTime,
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	expectedJson := `
	{
		"user": "c74a22da-8a05-43a9-a8b9-717e422b0af4",
		"key": "872e9e40-ce61-497e-b606-c7a08a4faa14",
		"validUntil": "2024-11-12T19:09:36Z"
	}`
	assert.JSONEq(expectedJson, string(out))
}

func TestUnit_ToApiKeyDtoResponse(t *testing.T) {
	assert := assert.New(t)

	entity := persistence.ApiKey{
		Id:         uuid.New(),
		Key:        uuid.New(),
		ApiUser:    uuid.New(),
		ValidUntil: someTime.Add(2 * time.Hour),
	}

	actual := ToApiKeyDtoResponse(entity)

	assert.Equal(entity.ApiUser, actual.User)
	assert.Equal(entity.Key, actual.Key)
	assert.Equal(entity.ValidUntil, actual.ValidUntil)
}
