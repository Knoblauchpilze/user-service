package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var defaultApiKeyId = uuid.MustParse("cc1742fa-77b4-4f5f-ac92-058c2e47a5d6")
var defaultApiKeyValue = uuid.MustParse("b01b9b1f-b651-4702-9b58-905b19584d69")
var defaultApiKey = persistence.ApiKey{
	Id:      defaultApiKeyId,
	Key:     defaultApiKeyValue,
	ApiUser: defaultUserId,
}

func TestApiKeyRepository_Create_DbInteraction(t *testing.T) {
	expectedSql := `
INSERT INTO api_key (id, key, api_user, valid_until)
	VALUES($1, $2, $3, $4)
	ON CONFLICT (api_user) DO UPDATE
	SET
		valid_until = excluded.valid_until
	WHERE
		api_key.api_user = excluded.api_user
	RETURNING
		api_key.key
`

	s := RepositoryPoolTestSuite{
		sqlMode: QueryBased,
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewApiKeyRepository(pool)
			_, err := repo.Create(context.Background(), defaultApiKey)
			return err
		},
		expectedSql: expectedSql,
		expectedArguments: []interface{}{
			defaultApiKey.Id,
			defaultApiKey.Key,
			defaultApiKey.ApiUser,
			defaultApiKey.ValidUntil,
		},
	}

	suite.Run(t, &s)
}

func TestApiKeyRepository_Create_RetrievesGeneratedApiKey(t *testing.T) {
	s := RepositorySingleValueTestSuite{
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewApiKeyRepository(pool)
			_, err := repo.Create(ctx, defaultApiKey)
			return err
		},
		expectedScanCalls: 1,
		expectedScannedProps: [][]interface{}{
			{&uuid.UUID{}},
		},
	}

	suite.Run(t, &s)
}

func TestApiKeyRepository_Create_ReturnsInputApiKey(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)

	actual, err := repo.Create(context.Background(), defaultApiKey)

	assert.Nil(err)
	assert.Equal(defaultApiKey, actual)
}

func TestApiKeyRepository_Get_DbInteraction(t *testing.T) {
	s := RepositoryPoolTestSuite{
		sqlMode: QueryBased,
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewApiKeyRepository(pool)
			_, err := repo.Get(context.Background(), defaultApiKeyId)
			return err
		},
		expectedSql: `SELECT id, key, api_user, valid_until FROM api_key WHERE id = $1`,
		expectedArguments: []interface{}{
			defaultApiKeyId,
		},
	}

	suite.Run(t, &s)
}

func TestApiKeyRepository_Get_InterpretDbData(t *testing.T) {
	s := RepositorySingleValueTestSuite{
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewApiKeyRepository(pool)
			_, err := repo.Get(ctx, defaultApiKeyId)
			return err
		},
		expectedScanCalls: 1,
		expectedScannedProps: [][]interface{}{
			{
				&uuid.UUID{},
				&uuid.UUID{},
				&uuid.UUID{},
				&time.Time{},
			},
		},
	}

	suite.Run(t, &s)
}

func TestApiKeyRepository_GetForKey_DbInteraction(t *testing.T) {
	s := RepositoryPoolTestSuite{
		sqlMode: QueryBased,
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewApiKeyRepository(pool)
			_, err := repo.GetForKey(context.Background(), defaultApiKeyValue)
			return err
		},
		expectedSql: `SELECT id, key, api_user, valid_until FROM api_key WHERE key = $1`,
		expectedArguments: []interface{}{
			defaultApiKeyValue,
		},
	}

	suite.Run(t, &s)
}

func TestApiKeyRepository_GetForKey_InterpretDbData(t *testing.T) {
	s := RepositorySingleValueTestSuite{
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewApiKeyRepository(pool)
			_, err := repo.GetForKey(ctx, defaultApiKeyValue)
			return err
		},
		expectedScanCalls: 1,
		expectedScannedProps: [][]interface{}{
			{
				&uuid.UUID{},
				&uuid.UUID{},
				&uuid.UUID{},
				&time.Time{},
			},
		},
	}

	suite.Run(t, &s)
}

func TestApiKeyRepository_GetForUser_DbInteraction(t *testing.T) {
	s := RepositoryPoolTestSuite{
		sqlMode: QueryBased,
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewApiKeyRepository(pool)
			_, err := repo.GetForUser(context.Background(), defaultUserId)
			return err
		},
		expectedSql: `SELECT id FROM api_key WHERE api_user = $1`,
		expectedArguments: []interface{}{
			defaultUserId,
		},
	}

	suite.Run(t, &s)
}

func TestApiKeyRepository_GetForUser_InterpretDbData(t *testing.T) {
	s := RepositoryGetAllTestSuite{
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewApiKeyRepository(pool)
			_, err := repo.GetForUser(ctx, defaultUserId)
			return err
		},
		expectedScanCalls: 1,
		expectedScannedProps: [][]interface{}{
			{&uuid.UUID{}},
		},
	}

	suite.Run(t, &s)
}

func TestApiKeyRepository_GetForUserTx_DbInteraction(t *testing.T) {
	s := RepositoryTransactionTestSuite{
		sqlMode: QueryBased,
		testFunc: func(ctx context.Context, tx db.Transaction) error {
			repo := NewApiKeyRepository(&mockConnectionPool{})
			_, err := repo.GetForUserTx(context.Background(), tx, defaultUserId)
			return err
		},
		expectedSql: []string{`SELECT id FROM api_key WHERE api_user = $1`},
		expectedArguments: [][]interface{}{
			{
				defaultUserId,
			},
		},
	}

	suite.Run(t, &s)
}

func TestApiKeyRepository_GetForUserTx_InterpretDbData(t *testing.T) {
	s := RepositoryGetAllTransactionTestSuite{
		testFunc: func(ctx context.Context, tx db.Transaction) error {
			repo := NewApiKeyRepository(&mockConnectionPool{})
			_, err := repo.GetForUserTx(ctx, tx, defaultUserId)
			return err
		},
		expectedScanCalls: 1,
		expectedScannedProps: [][]interface{}{
			{&uuid.UUID{}},
		},
	}

	suite.Run(t, &s)
}

func TestApiKeyRepository_DeleteForUser_DbInteraction(t *testing.T) {
	s := RepositoryTransactionTestSuite{
		sqlMode: ExecBased,
		testFunc: func(ctx context.Context, tx db.Transaction) error {
			repo := NewApiKeyRepository(&mockConnectionPool{})
			return repo.DeleteForUser(context.Background(), tx, defaultUserId)
		},
		expectedSql: []string{
			`DELETE FROM api_key WHERE api_user = $1`,
		},
		expectedArguments: [][]interface{}{
			{defaultUserId},
		},
	}

	suite.Run(t, &s)
}

func TestApiKeyRepository_DeleteForUser_NominalCase(t *testing.T) {
	assert := assert.New(t)

	repo := NewAclRepository()
	mt := &mockTransaction{}

	err := repo.DeleteForUser(context.Background(), mt, defaultUserId)

	assert.Nil(err)
}
