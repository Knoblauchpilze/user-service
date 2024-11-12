package db

import (
	"context"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/db/pgx"
	"github.com/KnoblauchPilze/user-service/pkg/db/postgresql"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/google/uuid"
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
	conn := newTestConnection(t)

	err := conn.Ping(context.Background())
	assert := assert.New(t)
	assert.Nil(err)
}

func TestIT_Connection_Close(t *testing.T) {
	conn := newTestConnection(t)

	err := conn.Ping(context.Background())
	require.Nil(t, err)

	conn.Close(context.Background())
	err = conn.Ping(context.Background())
	assert := assert.New(t)
	assert.True(errors.IsErrorWithCode(err, NotConnected), "Actual err: %v", err)
}

func TestIT_Connection_BeginTx_TimeStampIsValid(t *testing.T) {
	conn := newTestConnection(t)

	beforeTx := time.Now()
	tx, err := conn.BeginTx(context.Background())

	assert := assert.New(t)
	assert.Nil(err)
	assert.True(beforeTx.Before(tx.TimeStamp()))
}

func TestIT_Connection_BeginTx_ClosedConnection(t *testing.T) {
	conn := newTestConnection(t)
	conn.Close(context.Background())

	tx, err := conn.BeginTx(context.Background())

	assert := assert.New(t)
	assert.Nil(tx)
	assert.True(errors.IsErrorWithCode(err, NotConnected), "Actual err: %v", err)
}

func TestIT_Connection_Exec_Select(t *testing.T) {
	conn := newTestConnection(t)

	affectedRows, err := conn.Exec(context.Background(), "SELECT COUNT(*) FROM my_table WHERE name = 'test-name'")

	assert := assert.New(t)
	assert.Equal(int64(1), affectedRows)
	assert.Nil(err)
}

func TestIT_Connection_Exec_Insert(t *testing.T) {
	conn := newTestConnection(t)

	id := uuid.New()
	// Also using a uuid for the name to easily generate characters
	name := uuid.New()
	affectedRows, err := conn.Exec(context.Background(), "INSERT INTO my_table VALUES ($1, $2)", id, name)

	assert := assert.New(t)
	assert.Equal(int64(1), affectedRows)
	assert.Nil(err)
}

func TestIT_Connection_Exec_InsertDuplicate(t *testing.T) {
	conn := newTestConnection(t)

	_, name := insertTestData(t, conn)
	id := uuid.New()

	affectedRows, err := conn.Exec(context.Background(), "INSERT INTO my_table VALUES ($1, $2)", id, name)

	assert := assert.New(t)
	assert.Equal(int64(0), affectedRows)
	assert.True(errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation), "Actual err: %v", err)
}

func TestIT_Connection_Exec_Update(t *testing.T) {
	conn := newTestConnection(t)
	id, _ := insertTestData(t, conn)
	newName := uuid.New().String()

	affectedRows, err := conn.Exec(context.Background(), "UPDATE my_table SET name = $1 WHERE id = $2", newName, id)
	assert := assert.New(t)
	assert.Equal(int64(1), affectedRows)
	assert.Nil(err)

	assertNameForId(t, conn, id, newName)
}

func TestIT_Connection_Exec_UpdateDuplicate(t *testing.T) {
	conn := newTestConnection(t)
	id, name := insertTestData(t, conn)

	affectedRows, err := conn.Exec(context.Background(), "UPDATE my_table SET name = $1 WHERE id = $2", "test-name", id)
	assert := assert.New(t)
	assert.Equal(int64(0), affectedRows)
	assert.True(errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation), "Actual err: %v", err)

	assertNameForId(t, conn, id, name.String())
}

func TestIT_Connection_Exec_Delete(t *testing.T) {
	conn := newTestConnection(t)
	id, _ := insertTestData(t, conn)

	affectedRows, err := conn.Exec(context.Background(), "DELETE FROM my_table WHERE id = $1", id)
	assert := assert.New(t)
	assert.Equal(int64(1), affectedRows)
	assert.Nil(err)

	assertIdDoesNotExist(t, conn, id)
}

func TestIT_Connection_Exec_WithArguments(t *testing.T) {
	conn := newTestConnection(t)

	affectedRows, err := conn.Exec(context.Background(), "SELECT COUNT(*) FROM my_table WHERE name = $1", "test-name")

	assert := assert.New(t)
	assert.Equal(int64(1), affectedRows)
	assert.Nil(err)
}

func TestIT_Connection_Exec_WrongSyntax(t *testing.T) {
	conn := newTestConnection(t)

	affectedRows, err := conn.Exec(context.Background(), "DESELECT COUNT(*) FROM my_table WHERE name = 'test-name'")

	assert := assert.New(t)
	assert.Equal(int64(0), affectedRows)
	assert.True(errors.IsErrorWithCode(err, pgx.GenericSqlError), "Actual err: %v", err)
}
