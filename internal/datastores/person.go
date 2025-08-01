package datastores

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/williabk198/jagsqlb"
	"github.com/williabk198/timeclock/internal/models"
)

type PersonStore interface {
	SqlDatastore[models.Person, uuid.UUID]
}

type personStore struct {
	dbConn     *sql.DB
	sqlBuilder jagsqlb.SqlBuilder
	tableName  string
}

// Add implements Store.
func (ps personStore) Add(ctx context.Context, item models.Person) (id uuid.UUID, err error) {
	query, params, err := ps.sqlBuilder.Insert(ps.tableName).Data(item).Returning("id").Build()
	if err != nil {
		return uuid.Nil, err
	}

	row := ps.dbConn.QueryRowContext(ctx, query, params...)
	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

// Delete implements Store.
func (ps personStore) Delete(ctx context.Context, id uuid.UUID) (item models.Person, err error) {
	panic("unimplemented")
}

// GetAllPaginated implements Store.
func (ps personStore) GetAllPaginated(ctx context.Context, offset uint, limit uint) (items []models.Person, err error) {
	panic("unimplemented")
}

// GetSpecific implements Store.
func (ps personStore) GetSpecific(ctx context.Context, id uuid.UUID) (item models.Person, err error) {
	panic("unimplemented")
}

// Update implements Store.
func (ps personStore) Update(ctx context.Context, id uuid.UUID, item models.Person) (err error) {
	panic("unimplemented")
}

func NewPersonStore(dbConn *sql.DB) PersonStore {
	return personStore{
		dbConn:     dbConn,
		sqlBuilder: jagsqlb.NewSqlBuilder(),
		tableName:  "person.persons",
	}
}
