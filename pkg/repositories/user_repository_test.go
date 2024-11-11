package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestIT_UserRepository(t *testing.T) {
	repo, conn := newTestUserRepository(t)

	user := persistence.User{
		Id:        uuid.New(),
		Email:     "my-email-" + uuid.New().String(),
		Password:  "my-password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   6,
	}

	actual, err := repo.Create(context.Background(), user)
	assert := assert.New(t)
	assert.Nil(err)

	assert.Equal(user, actual)

	assertUserExists(t, conn, user.Id)
}

func newTestUserRepository(t *testing.T) (UserRepository, db.Connection) {
	conn := newTestConnection(t)
	return NewUserRepository(conn), conn
}
