package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

type mockPgxTransaction struct {
	pgx.Tx

	queryCalled    int
	execCalled     int
	rollbackCalled int
	commitCalled   int

	sql       string
	arguments []interface{}

	tag pgconn.CommandTag
	err error
}

func TestTransaction_Query_DelegatesToTransaction(t *testing.T) {
	assert := assert.New(t)

	mt := mockPgxTransaction{}
	tx := transactionImpl{
		tx: &mt,
	}

	tx.Query(context.Background(), exampleSqlQuery)

	assert.Equal(1, mt.queryCalled)
}

func TestTransaction_Query_PropagatesSqlQuery(t *testing.T) {
	assert := assert.New(t)

	mt := mockPgxTransaction{}
	tx := transactionImpl{
		tx: &mt,
	}

	tx.Query(context.Background(), exampleSqlQuery)

	assert.Equal(exampleSqlQuery, mt.sql)
}

func TestTransaction_Query_PropagatesSqlArguments(t *testing.T) {
	assert := assert.New(t)

	mt := mockPgxTransaction{}
	tx := transactionImpl{
		tx: &mt,
	}

	tx.Query(context.Background(), exampleSqlQuery, 1, "test-str")

	assert.Equal([]interface{}{1, "test-str"}, mt.arguments)
}

func TestTransaction_Query_PropagatesError(t *testing.T) {
	assert := assert.New(t)

	mt := mockPgxTransaction{
		err: errDefault,
	}
	tx := transactionImpl{
		tx: &mt,
	}

	actual := tx.Query(context.Background(), exampleSqlQuery)

	assert.Equal(errDefault, tx.err)
	assert.Equal(errDefault, actual.Err())
}

func TestTransaction_Exec_DelegatesToTransaction(t *testing.T) {
	assert := assert.New(t)

	mt := mockPgxTransaction{}
	tx := transactionImpl{
		tx: &mt,
	}

	tx.Exec(context.Background(), exampleExecQuery)

	assert.Equal(1, mt.execCalled)
}

func TestTransaction_Exec_PropagatesSqlQuery(t *testing.T) {
	assert := assert.New(t)

	mt := mockPgxTransaction{}
	tx := transactionImpl{
		tx: &mt,
	}

	tx.Exec(context.Background(), exampleExecQuery)

	assert.Equal(exampleExecQuery, mt.sql)
}

func TestTransaction_Exec_PropagatesSqlArguments(t *testing.T) {
	assert := assert.New(t)

	mt := mockPgxTransaction{}
	tx := transactionImpl{
		tx: &mt,
	}

	tx.Exec(context.Background(), exampleExecQuery, 1, "test-str")

	assert.Equal([]interface{}{1, "test-str"}, mt.arguments)
}

func TestTransaction_Exec_PropagatesError(t *testing.T) {
	assert := assert.New(t)

	mt := mockPgxTransaction{
		err: errDefault,
	}
	tx := transactionImpl{
		tx: &mt,
	}

	_, err := tx.Exec(context.Background(), exampleExecQuery)

	assert.Equal(errDefault, tx.err)
	assert.Equal(errDefault, err)
}

func TestTransaction_Exec_PropagatesCommandTag(t *testing.T) {
	assert := assert.New(t)

	mt := mockPgxTransaction{
		tag: pgconn.NewCommandTag("INSERT 0 1"),
	}
	tx := transactionImpl{
		tx: &mt,
	}

	actual, _ := tx.Exec(context.Background(), exampleExecQuery)

	assert.Equal(1, actual)
}

func TestTransaction_Close_WhenError_CallsRollback(t *testing.T) {
	assert := assert.New(t)

	mt := mockPgxTransaction{}
	tx := transactionImpl{
		tx:  &mt,
		err: errDefault,
	}

	tx.Close(context.Background())

	assert.Equal(1, mt.rollbackCalled)
}

func TestTransaction_Close_WhenNoError_CallsCommit(t *testing.T) {
	assert := assert.New(t)

	mt := mockPgxTransaction{}
	tx := transactionImpl{
		tx: &mt,
	}

	tx.Close(context.Background())

	assert.Equal(1, mt.commitCalled)
}

func (m *mockPgxTransaction) Rollback(ctx context.Context) error {
	m.rollbackCalled++
	return m.err
}

func (m *mockPgxTransaction) Commit(ctx context.Context) error {
	m.commitCalled++
	return m.err
}

func (m *mockPgxTransaction) Query(ctx context.Context, sql string, arguments ...interface{}) (pgx.Rows, error) {
	m.queryCalled++
	m.sql = sql
	m.arguments = append(m.arguments, arguments...)
	return nil, m.err
}

func (m *mockPgxTransaction) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	m.execCalled++
	m.sql = sql
	m.arguments = append(m.arguments, arguments...)
	return m.tag, m.err
}
