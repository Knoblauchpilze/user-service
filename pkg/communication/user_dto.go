package communication

import (
	"time"

	"github.com/Knoblauchpilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

// https://stackoverflow.com/questions/18635671/how-to-define-multiple-name-tags-in-a-struct
type UserDtoRequest struct {
	Email    string `json:"email" form:"email" example:"user@example.com"`
	Password string `json:"password" form:"password" example:"SecurePassword123"`
}

type UserDtoResponse struct {
	Id       uuid.UUID `json:"id" format:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email    string    `json:"email" example:"user@example.com"`
	Password string    `json:"password" example:"SecurePassword123"`

	CreatedAt time.Time `json:"createdAt" format:"date-time" example:"2026-04-27T20:56:59Z"`
}

func FromUserDtoRequest(user UserDtoRequest) persistence.User {
	t := time.Now()
	return persistence.User{
		Id:       uuid.New(),
		Email:    user.Email,
		Password: user.Password,

		CreatedAt: t,
		UpdatedAt: t,
	}
}

func ToUserDtoResponse(user persistence.User) UserDtoResponse {
	return UserDtoResponse{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,

		CreatedAt: user.CreatedAt,
	}
}
