package datastores

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/williabk198/jagsqlb"
	"github.com/williabk198/jagsqlb/condition"
	"github.com/williabk198/timeclock/internal/models"
)

// UserDatastore handles persistence of user login credentials and the metadata about
// a user's account (login history, lockout status, etc...). It intentionally does not
// implement the generic SqlDatastore interface since a user's "identity" data (email,
// password) and the person it belongs to are owned by the admin service.
type UserDatastore interface {
	Add(ctx context.Context, user models.User, metadata models.UserMetadata) (id uuid.UUID, err error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	GetSpecific(ctx context.Context, id uuid.UUID) (models.User, error)
	GetMetadata(ctx context.Context, userID uuid.UUID) (models.UserMetadata, error)
	UpdateMetadata(ctx context.Context, userID uuid.UUID, metadata models.UserMetadata) error
	RecordSuccessfulLogin(ctx context.Context, userID uuid.UUID, loginTime time.Time) error
	// GetGrantedScopes returns the scopes that have been explicitly granted to a user.
	// It does not include the "default" scope, which is implicit for every user.
	GetGrantedScopes(ctx context.Context, userID uuid.UUID) ([]string, error)
	// GrantScope grants a user a scope. Granting a scope the user already has returns
	// an error (the underlying table has a unique constraint on (user_id, scope)).
	GrantScope(ctx context.Context, userID uuid.UUID, scope string) error
	// RevokeScope revokes a previously granted scope from a user. Revoking a scope the
	// user doesn't have is a no-op.
	RevokeScope(ctx context.Context, userID uuid.UUID, scope string) error
}

type userSqlStore struct {
	dbConn            *sql.DB
	sqlBuilder        jagsqlb.SqlBuilder
	tableName         string
	metadataTableName string
	scopesTableName   string
}

// Add implements UserDatastore.
func (uss userSqlStore) Add(ctx context.Context, user models.User, metadata models.UserMetadata) (id uuid.UUID, err error) {
	query, params, err := uss.sqlBuilder.Insert(uss.tableName).Data(user).Returning("id").Build()
	if err != nil {
		return uuid.Nil, err
	}

	tx, err := uss.dbConn.BeginTx(ctx, nil)
	if err != nil {
		return uuid.Nil, err
	}

	row := tx.QueryRowContext(ctx, query, params...)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}
	metadata.UserID = id

	query, params, err = uss.sqlBuilder.Insert(uss.metadataTableName).Data(metadata).Build()
	if err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	if _, err := tx.ExecContext(ctx, query, params...); err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	if err := tx.Commit(); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

// GetByEmail implements UserDatastore.
func (uss userSqlStore) GetByEmail(ctx context.Context, email string) (models.User, error) {
	query, params, err := uss.sqlBuilder.Select(uss.tableName, "*").Where(condition.Equals("email", email)).Build()
	if err != nil {
		return models.User{}, err
	}

	row := uss.dbConn.QueryRowContext(ctx, query, params...)
	return uss.userFromRow(row)
}

// GetSpecific implements UserDatastore.
func (uss userSqlStore) GetSpecific(ctx context.Context, id uuid.UUID) (models.User, error) {
	query, params, err := uss.sqlBuilder.Select(uss.tableName, "*").Where(condition.Equals("id", id)).Build()
	if err != nil {
		return models.User{}, err
	}

	row := uss.dbConn.QueryRowContext(ctx, query, params...)
	return uss.userFromRow(row)
}

// GetMetadata implements UserDatastore.
func (uss userSqlStore) GetMetadata(ctx context.Context, userID uuid.UUID) (models.UserMetadata, error) {
	query, params, err := uss.sqlBuilder.Select(uss.metadataTableName, "*").Where(condition.Equals("user_id", userID)).Build()
	if err != nil {
		return models.UserMetadata{}, err
	}

	row := uss.dbConn.QueryRowContext(ctx, query, params...)
	return uss.metadataFromRow(row)
}

// UpdateMetadata implements UserDatastore.
func (uss userSqlStore) UpdateMetadata(ctx context.Context, userID uuid.UUID, metadata models.UserMetadata) error {
	query, params, err := uss.sqlBuilder.Update(uss.metadataTableName).SetStruct(metadata).Where(condition.Equals("user_id", userID)).Build()
	if err != nil {
		return err
	}

	_, err = uss.dbConn.ExecContext(ctx, query, params...)
	return err
}

// RecordSuccessfulLogin implements UserDatastore.
func (uss userSqlStore) RecordSuccessfulLogin(ctx context.Context, userID uuid.UUID, loginTime time.Time) error {
	query, params, err := uss.sqlBuilder.Update(uss.metadataTableName).SetMap(map[string]any{
		"last_login":    loginTime,
		"failure_count": 0,
	}).Where(condition.Equals("user_id", userID)).Build()
	if err != nil {
		return err
	}

	_, err = uss.dbConn.ExecContext(ctx, query, params...)
	return err
}

// GetGrantedScopes implements UserDatastore.
func (uss userSqlStore) GetGrantedScopes(ctx context.Context, userID uuid.UUID) ([]string, error) {
	query, params, err := uss.sqlBuilder.Select(uss.scopesTableName, "scope").Where(condition.Equals("user_id", userID)).Build()
	if err != nil {
		return nil, err
	}

	rows, err := uss.dbConn.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scopes := make([]string, 0)
	for rows.Next() {
		var scope string
		if err := rows.Scan(&scope); err != nil {
			return nil, err
		}
		scopes = append(scopes, scope)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return scopes, nil
}

// GrantScope implements UserDatastore.
func (uss userSqlStore) GrantScope(ctx context.Context, userID uuid.UUID, scope string) error {
	query, params, err := uss.sqlBuilder.Insert(uss.scopesTableName).Columns("user_id", "scope").Values([]any{userID, scope}).Build()
	if err != nil {
		return err
	}

	_, err = uss.dbConn.ExecContext(ctx, query, params...)
	return err
}

// RevokeScope implements UserDatastore.
func (uss userSqlStore) RevokeScope(ctx context.Context, userID uuid.UUID, scope string) error {
	query, params, err := uss.sqlBuilder.Delete(uss.scopesTableName).Where(condition.GroupedAnd(
		condition.Equals("user_id", userID),
		condition.Equals("scope", scope),
	)).Build()
	if err != nil {
		return err
	}

	_, err = uss.dbConn.ExecContext(ctx, query, params...)
	return err
}

func (uss userSqlStore) userFromRow(row *sql.Row) (models.User, error) {
	var item models.User
	if err := row.Scan(&item.ID, &item.PersonID, &item.Email, &item.HashedPass); err != nil {
		return models.User{}, err
	}

	return item, nil
}

func (uss userSqlStore) metadataFromRow(row *sql.Row) (models.UserMetadata, error) {
	var item models.UserMetadata
	if err := row.Scan(
		&item.UserID, &item.CreatedAt, &item.FailedLoginCount,
		&item.LastLogin, &item.LastPasswordReset, &item.Status,
	); err != nil {
		return models.UserMetadata{}, err
	}

	return item, nil
}

func NewUserStore(dbConn *sql.DB) UserDatastore {
	return userSqlStore{
		dbConn:            dbConn,
		sqlBuilder:        jagsqlb.NewSqlBuilder(),
		tableName:         "auth.users",
		metadataTableName: "auth.user_metadata",
		scopesTableName:   "auth.user_scopes",
	}
}
