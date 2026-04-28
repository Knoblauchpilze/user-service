package communication

import (
	"time"

	"github.com/Knoblauchpilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type ApiKeyDtoResponse struct {
	User       uuid.UUID `json:"user" format:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Key        uuid.UUID `json:"key" format:"uuid" example:"f47ac10b-58cc-4372-a567-0e02b2c3d479"`
	ValidUntil time.Time `json:"validUntil" format:"date-time" example:"2026-04-28T20:56:59Z"`
}

func ToApiKeyDtoResponse(apiKey persistence.ApiKey) ApiKeyDtoResponse {
	return ApiKeyDtoResponse{
		User:       apiKey.ApiUser,
		Key:        apiKey.Key,
		ValidUntil: apiKey.ValidUntil,
	}
}
