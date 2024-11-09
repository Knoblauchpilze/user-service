package db

import (
	"context"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/db/postgresql"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnit_New_InvalidConfiguration(t *testing.T) {
	config := postgresql.Config{
		Host: ":/not-a-host",
	}

	conn, err := New(context.Background(), config)

	assert := assert.New(t)
	assert.Nil(conn)
	assert.NotNil(err)
}

func TestIT_New_ValidConfiguration(t *testing.T) {
	conn, err := New(context.Background(), dbTestConfig)

	assert := assert.New(t)
	assert.NotNil(conn)
	assert.Nil(err)
}

func TestIT_Connection_Ping(t *testing.T) {
	conn, err := New(context.Background(), dbTestConfig)
	require.Nil(t, err)

	err = conn.Ping(context.Background())
	assert := assert.New(t)
	assert.Nil(err)
}

func TestIT_Connection_Close(t *testing.T) {
	conn, err := New(context.Background(), dbTestConfig)
	require.Nil(t, err)

	err = conn.Ping(context.Background())
	require.Nil(t, err)

	conn.Close(context.Background())
	err = conn.Ping(context.Background())
	assert := assert.New(t)
	assert.True(errors.IsErrorWithCode(err, NotConnected))
}
