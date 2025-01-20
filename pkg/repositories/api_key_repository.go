package repositories

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type ApiKeyRepository interface {
	Create(ctx context.Context, apiKey persistence.ApiKey) (persistence.ApiKey, error)
	Get(ctx context.Context, id uuid.UUID) (persistence.ApiKey, error)
	GetForKey(ctx context.Context, apiKey uuid.UUID) (persistence.ApiKey, error)
	GetForUser(ctx context.Context, user uuid.UUID) (persistence.ApiKey, error)
	DeleteForUser(ctx context.Context, tx db.Transaction, user uuid.UUID) error
}

type apiKeyRepositoryImpl struct {
	conn db.Connection
}

func NewApiKeyRepository(conn db.Connection) ApiKeyRepository {
	return &apiKeyRepositoryImpl{
		conn: conn,
	}
}

const createApiKeySqlTemplate = `
INSERT INTO api_key (id, key, api_user, valid_until)
	VALUES($1, $2, $3, $4)
	ON CONFLICT (api_user) DO UPDATE
	SET
		valid_until = excluded.valid_until
	WHERE
		api_key.api_user = excluded.api_user
	RETURNING
		api_key.id,
		api_key.key
`

func (r *apiKeyRepositoryImpl) Create(ctx context.Context, apiKey persistence.ApiKey) (persistence.ApiKey, error) {
	type apiKeyDetails struct {
		Id  uuid.UUID
		Key uuid.UUID
	}
	keyDetails, err := db.QueryOne[apiKeyDetails](ctx, r.conn, createApiKeySqlTemplate, apiKey.Id, apiKey.Key, apiKey.ApiUser, apiKey.ValidUntil)
	if err != nil {
		return persistence.ApiKey{}, err
	}

	apiKey.Id = keyDetails.Id
	apiKey.Key = keyDetails.Key
	return apiKey, nil
}

const getApiKeySqlTemplate = `
SELECT
	id, key, api_user, valid_until
FROM
	api_key
WHERE
	id = $1`

func (r *apiKeyRepositoryImpl) Get(ctx context.Context, id uuid.UUID) (persistence.ApiKey, error) {
	return db.QueryOne[persistence.ApiKey](ctx, r.conn, getApiKeySqlTemplate, id)
}

const getApiKeyForKeySqlTemplate = `
SELECT
	id, key, api_user, valid_until
FROM
	api_key
WHERE
	key = $1`

func (r *apiKeyRepositoryImpl) GetForKey(ctx context.Context, apiKey uuid.UUID) (persistence.ApiKey, error) {
	return db.QueryOne[persistence.ApiKey](ctx, r.conn, getApiKeyForKeySqlTemplate, apiKey)
}

const getApiKeyForUserSqlTemplate = `
SELECT
	id, key, api_user, valid_until
FROM
	api_key
WHERE
	api_user = $1`

func (r *apiKeyRepositoryImpl) GetForUser(ctx context.Context, user uuid.UUID) (persistence.ApiKey, error) {
	return db.QueryOne[persistence.ApiKey](ctx, r.conn, getApiKeyForUserSqlTemplate, user)
}

const deleteApiKeyForUserSqlTemplate = `
DELETE FROM
	api_key
WHERE
	api_user = $1`

func (r *apiKeyRepositoryImpl) DeleteForUser(ctx context.Context, tx db.Transaction, user uuid.UUID) error {
	_, err := tx.Exec(ctx, deleteApiKeyForUserSqlTemplate, user)
	return err
}
