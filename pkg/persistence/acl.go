package persistence

import (
	"time"

	"github.com/google/uuid"
)

type Acl struct {
	Id          uuid.UUID
	UserGroup   uuid.UUID
	Policy      string
	Resource    string
	Description string
	CreatedAt   time.Time
	Permissions []string
}
