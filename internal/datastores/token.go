package datastores

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/williabk198/jagsqlb"
	"github.com/williabk198/jagsqlb/condition"
	"github.com/williabk198/timeclock/internal/models"
)

// TokenDatastore handles persistence of opaque token records issued by the auth service.
// Only the hash of a token is ever stored or queried by; the raw token value never
// touches the database.
type TokenDatastore interface {
	Add(ctx context.Context, token models.Token) (id uuid.UUID, err error)
	GetByHash(ctx context.Context, hash string) (models.Token, error)
	Revoke(ctx context.Context, id uuid.UUID) error
	RevokeByHash(ctx context.Context, hash string) error
}

type tokenSqlStore struct {
	dbConn     *sql.DB
	sqlBuilder jagsqlb.SqlBuilder
	tableName  string
}

// Add implements TokenDatastore.
func (tss tokenSqlStore) Add(ctx context.Context, token models.Token) (id uuid.UUID, err error) {
	query, params, err := tss.sqlBuilder.Insert(tss.tableName).Data(token).Returning("id").Build()
	if err != nil {
		return uuid.Nil, err
	}

	row := tss.dbConn.QueryRowContext(ctx, query, params...)
	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

// GetByHash implements TokenDatastore.
func (tss tokenSqlStore) GetByHash(ctx context.Context, hash string) (models.Token, error) {
	query, params, err := tss.sqlBuilder.Select(tss.tableName, "*").Where(condition.Equals("token_hash", hash)).Build()
	if err != nil {
		return models.Token{}, err
	}

	row := tss.dbConn.QueryRowContext(ctx, query, params...)
	return tss.tokenFromRow(row)
}

// Revoke implements TokenDatastore.
func (tss tokenSqlStore) Revoke(ctx context.Context, id uuid.UUID) error {
	query, params, err := tss.sqlBuilder.Update(tss.tableName).SetMap(map[string]any{"revoked": true}).Where(condition.Equals("id", id)).Build()
	if err != nil {
		return err
	}

	_, err = tss.dbConn.ExecContext(ctx, query, params...)
	return err
}

// RevokeByHash implements TokenDatastore.
func (tss tokenSqlStore) RevokeByHash(ctx context.Context, hash string) error {
	query, params, err := tss.sqlBuilder.Update(tss.tableName).SetMap(map[string]any{"revoked": true}).Where(condition.Equals("token_hash", hash)).Build()
	if err != nil {
		return err
	}

	_, err = tss.dbConn.ExecContext(ctx, query, params...)
	return err
}

func (tss tokenSqlStore) tokenFromRow(row *sql.Row) (models.Token, error) {
	var item models.Token
	if err := row.Scan(
		&item.ID, &item.UserID, &item.Hash, &item.Scope,
		&item.Kind, &item.IssuedAt, &item.ExpiresAt, &item.Revoked,
	); err != nil {
		return models.Token{}, err
	}

	return item, nil
}

func NewTokenStore(dbConn *sql.DB) TokenDatastore {
	return tokenSqlStore{
		dbConn:     dbConn,
		sqlBuilder: jagsqlb.NewSqlBuilder(),
		tableName:  "auth.tokens",
	}
}
