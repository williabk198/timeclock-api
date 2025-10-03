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
	AddSpecificContactAddress(ctx context.Context, address models.ContactAddress) (uuid.UUID, error)
	AddSpecificContactEmail(ctx context.Context, email models.ContactEmail) (uuid.UUID, error)
	AddSpecificContactPhone(ctx context.Context, phone models.ContactPhone) (uuid.UUID, error)
	GetSpecificContactAddresses(ctx context.Context, id uuid.UUID) ([]models.ContactAddress, error)
	GetSpecificContactEmails(ctx context.Context, id uuid.UUID) ([]models.ContactEmail, error)
	GetSpecificContactPhones(ctx context.Context, id uuid.UUID) ([]models.ContactPhone, error)
}

type personStore struct {
	dbConn     *sql.DB
	sqlBuilder jagsqlb.SqlBuilder
	tableName  string // TODO: Remove

	tableNameMap map[string]string
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

// AddSpecificContactAddress implements PersonStore.
func (ps personStore) AddSpecificContactAddress(ctx context.Context, address models.ContactAddress) (id uuid.UUID, err error) {
	query, params, err := ps.sqlBuilder.Insert(ps.tableNameMap["addresses"]).Data(address).Returning("id").Build()
	if err != nil {
		return uuid.Nil, err
	}

	row := ps.dbConn.QueryRowContext(ctx, query, params...)
	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

// AddSpecificContactEmail implements PersonStore.
func (ps personStore) AddSpecificContactEmail(ctx context.Context, email models.ContactEmail) (id uuid.UUID, err error) {
	query, params, err := ps.sqlBuilder.Insert(ps.tableNameMap["emails"]).Data(email).Returning("id").Build()
	if err != nil {
		return uuid.Nil, err
	}

	row := ps.dbConn.QueryRowContext(ctx, query, params...)
	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

// AddSpecificContactPhone implements PersonStore.
func (ps personStore) AddSpecificContactPhone(ctx context.Context, phone models.ContactPhone) (id uuid.UUID, err error) {
	query, params, err := ps.sqlBuilder.Insert(ps.tableNameMap["phones"]).Data(phone).Returning("id").Build()
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

	row := ps.dbConn.QueryRowContext(ctx, query, params...)
	return ps.personFromRow(row)
}

// GetAllPaginated implements Store.
func (ps personStore) GetAllPaginated(ctx context.Context, offset uint, limit uint) (items []models.Person, err error) {
	query, params, err := ps.sqlBuilder.Select(ps.tableName, "*").Offset(offset).Limit(limit).Build()
	if err != nil {
		return nil, err
	}

	rows, err := ps.dbConn.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return ps.personSliceFromRows(rows)
}

// GetSpecific implements Store.
func (ps personStore) GetSpecific(ctx context.Context, id uuid.UUID) (item models.Person, err error) {
	query, params, err := ps.sqlBuilder.Select(ps.tableName, "*").Where(condition.Equals("id", id)).Build()
	if err != nil {
		return models.Person{}, err
	}
	row := ps.dbConn.QueryRowContext(ctx, query, params...)
	return ps.personFromRow(row)
}

// GetSpecificContactAddresses implements PersonStore.
func (ps personStore) GetSpecificContactAddresses(ctx context.Context, id uuid.UUID) ([]models.ContactAddress, error) {
	query, params, err := ps.sqlBuilder.Select(ps.tableNameMap["addresses"], "*").Where(condition.Equals("person_id", id)).Build()
	if err != nil {
		return nil, err
	}

	rows, err := ps.dbConn.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return ps.personAddressFromRows(rows)
}

// GetSpecificContacts implements PersonStore.
func (ps personStore) GetSpecificContactEmails(ctx context.Context, id uuid.UUID) ([]models.ContactEmail, error) {
	query, params, err := ps.sqlBuilder.Select(ps.tableNameMap["emails"], "*").Where(condition.Equals("person_id", id)).Build()
	if err != nil {
		return nil, err
	}

	rows, err := ps.dbConn.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	return ps.personEmailFromRows(rows)
}

// GetSpecificContactPhones implements PersonStore.
func (ps personStore) GetSpecificContactPhones(ctx context.Context, id uuid.UUID) ([]models.ContactPhone, error) {
	query, params, err := ps.sqlBuilder.Select(ps.tableNameMap["phones"], "*").Where(condition.Equals("person_id", id)).Build()
	if err != nil {
		return nil, err
	}

	rows, err := ps.dbConn.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	return ps.personPhoneFromRows(rows)
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

func (ps personStore) personFromRow(row *sql.Row) (models.Person, error) {
	var rawPronounVal string
	var item models.Person
	if err := row.Scan(
		&item.ID, &item.Name.GivenName, &item.Name.FamilyName, &item.Name.FamilyNameFirst,
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

func (ps personStore) personSliceFromRows(rows *sql.Rows) ([]models.Person, error) {
	result := make([]models.Person, 0)
	for rows.Next() {
		var rawPronounVal string
		var item models.Person
		if err := rows.Scan(
			&item.ID, &item.Name.GivenName, &item.Name.FamilyName, &item.Name.FamilyNameFirst,
			&item.DateOfBirth, &item.Gender, &rawPronounVal,
		); err != nil {
			return nil, err
		}

		pronouns, err := utils.ParsePronouns(rawPronounVal)
		if err != nil {
			return nil, err
		}
		item.Pronouns = pronouns

		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (ps personStore) personAddressFromRows(rows *sql.Rows) ([]models.ContactAddress, error) {
	results := make([]models.ContactAddress, 0)

	for rows.Next() {
		var item models.ContactAddress
		if err := rows.Scan(
			&item.ID, &item.PersonID, &item.Street1, &item.Street2, &item.Locality, &item.Region,
			&item.PostalCode, &item.Country, &item.Type, &item.Primary,
		); err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (ps personStore) personEmailFromRows(rows *sql.Rows) ([]models.ContactEmail, error) {
	results := make([]models.ContactEmail, 0)

	for rows.Next() {
		var item models.ContactEmail
		if err := rows.Scan(
			&item.ID, &item.PersonID, &item.Username, &item.Provider, &item.Primary,
		); err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (ps personStore) personPhoneFromRows(rows *sql.Rows) ([]models.ContactPhone, error) {
	results := make([]models.ContactPhone, 0)

	for rows.Next() {
		var item models.ContactPhone
		if err := rows.Scan(
			&item.ID, &item.PersonID, &item.CountryCode, &item.PhoneNumber, &item.Type, &item.Primary,
		); err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func NewPersonStore(dbConn *sql.DB) PersonStore {
	return personStore{
		dbConn:     dbConn,
		sqlBuilder: jagsqlb.NewSqlBuilder(),
		tableName:  "person.persons", // TODO: Remove
		tableNameMap: map[string]string{
			"persons":   "person.persons",
			"addresses": "person.addresses",
			"emails":    "person.emails",
			"phones":    "person.phones",
		},
	}
}
