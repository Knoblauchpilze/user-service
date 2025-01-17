package persistence

import (
	"time"

	"github.com/google/uuid"
)

type UserGroup struct {
	Id          uuid.UUID
	ApiUser     uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
	Acls        []Acl
}
