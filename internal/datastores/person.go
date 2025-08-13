package datastores

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/williabk198/jagsqlb"
	"github.com/williabk198/jagsqlb/condition"
	"github.com/williabk198/timeclock/internal/models"
	"github.com/williabk198/timeclock/internal/utils"
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
	query, params, err := ps.sqlBuilder.Delete(ps.tableName).Where(condition.Equals("id", id)).Returning("*").Build()
	if err != nil {
		return models.Person{}, err
	}

	var rawPronounVal string
	var person models.Person
	row := ps.dbConn.QueryRowContext(ctx, query, params...)
	if err := row.Scan(
		&id, &person.Name.GivenName, &person.Name.FamilyName, &person.Name.FamilyNameFirst,
		&person.DateOfBirth, &person.Gender, &rawPronounVal,
	); err != nil {
		return models.Person{}, err
	}

	person.Pronouns, err = utils.ParsePronouns(rawPronounVal)
	if err != nil {
		return models.Person{}, nil
	}

	return person, nil
}

// GetAllPaginated implements Store.
func (ps personStore) GetAllPaginated(ctx context.Context, offset uint, limit uint) (items []models.Person, err error) {
	panic("unimplemented")
}

// GetSpecific implements Store.
func (ps personStore) GetSpecific(ctx context.Context, id uuid.UUID) (item models.Person, err error) {
	query, params, err := ps.sqlBuilder.Select(ps.tableName, "*").Where(condition.Equals("id", id)).Build()
	if err != nil {
		return models.Person{}, err
	}

	var rawPronounVal string
	row := ps.dbConn.QueryRowContext(ctx, query, params...)
	if err := row.Scan(
		&id, &item.Name.GivenName, &item.Name.FamilyName, &item.Name.FamilyNameFirst,
		&item.DateOfBirth, &item.Gender, &rawPronounVal,
	); err != nil {
		return models.Person{}, err
	}

	pronouns, err := utils.ParsePronouns(rawPronounVal)
	if err != nil {
		return models.Person{}, err
	}
	item.Pronouns = pronouns

	return item, nil
}

// Update implements Store.
func (ps personStore) Update(ctx context.Context, id uuid.UUID, item models.Person) (err error) {
	query, params, err := ps.sqlBuilder.Update(ps.tableName).SetStruct(item).Where(condition.Equals("id", id)).Build()
	if err != nil {
		return err
	}

	if _, err := ps.dbConn.ExecContext(ctx, query, params...); err != nil {
		return err
	}

	return nil
}

func NewPersonStore(dbConn *sql.DB) PersonStore {
	return personStore{
		dbConn:     dbConn,
		sqlBuilder: jagsqlb.NewSqlBuilder(),
		tableName:  "person.persons",
	}
}
