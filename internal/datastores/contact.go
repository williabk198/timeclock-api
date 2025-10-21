package datastores

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/williabk198/jagsqlb"
	"github.com/williabk198/jagsqlb/condition"
	"github.com/williabk198/timeclock/internal/models"
)

type ContactDatastore interface {
	AddPersonAddress(ctx context.Context, address models.ContactAddress) (uuid.UUID, error)
	AddPersonEmail(ctx context.Context, email models.ContactEmail) (uuid.UUID, error)
	AddPersonPhone(ctx context.Context, phone models.ContactPhone) (uuid.UUID, error)
	DeletePersonAddress(ctx context.Context, personID, addressID uuid.UUID) (models.ContactAddress, error)
	DeletePersonEmail(ctx context.Context, personID, emailID uuid.UUID) (models.ContactEmail, error)
	DeletePersonPhone(ctx context.Context, personID, phoneID uuid.UUID) (models.ContactPhone, error)
	GetPersonAddresses(ctx context.Context, id uuid.UUID) ([]models.ContactAddress, error)
	GetPersonEmails(ctx context.Context, id uuid.UUID) ([]models.ContactEmail, error)
	GetPersonPhones(ctx context.Context, id uuid.UUID) ([]models.ContactPhone, error)
	UpdatePersonAddress(ctx context.Context, personID, addressID uuid.UUID, newVal models.ContactAddress) error
	UpdatePersonEmail(ctx context.Context, personID, emailID uuid.UUID, newVal models.ContactEmail) error
	UpdatePersonPhone(ctx context.Context, personID, phoneID uuid.UUID, newVal models.ContactPhone) error
}

type contactStore struct {
	sqlBuilder jagsqlb.SqlBuilder
	dbConn     *sql.DB

	tableNameMap map[string]string
}

// AddSpecificContactAddress implements PersonStore.
func (cs contactStore) AddPersonAddress(ctx context.Context, address models.ContactAddress) (id uuid.UUID, err error) {
	query, params, err := cs.sqlBuilder.Insert(cs.tableNameMap["addresses"]).Data(address).Returning("id").Build()
	if err != nil {
		return uuid.Nil, err
	}

	row := cs.dbConn.QueryRowContext(ctx, query, params...)
	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

// AddSpecificContactEmail implements PersonStore.
func (cs contactStore) AddPersonEmail(ctx context.Context, email models.ContactEmail) (id uuid.UUID, err error) {
	query, params, err := cs.sqlBuilder.Insert(cs.tableNameMap["emails"]).Data(email).Returning("id").Build()
	if err != nil {
		return uuid.Nil, err
	}

	row := cs.dbConn.QueryRowContext(ctx, query, params...)
	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

// AddSpecificContactPhone implements PersonStore.
func (cs contactStore) AddPersonPhone(ctx context.Context, phone models.ContactPhone) (id uuid.UUID, err error) {
	query, params, err := cs.sqlBuilder.Insert(cs.tableNameMap["phones"]).Data(phone).Returning("id").Build()
	if err != nil {
		return uuid.Nil, err
	}

	row := cs.dbConn.QueryRowContext(ctx, query, params...)
	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

// DeletePersonAddress implements ContactDatastore.
func (cs contactStore) DeletePersonAddress(ctx context.Context, personID uuid.UUID, addressID uuid.UUID) (models.ContactAddress, error) {
	panic("unimplemented")
}

// DeletePersonEmail implements ContactDatastore.
func (cs contactStore) DeletePersonEmail(ctx context.Context, personID uuid.UUID, emailID uuid.UUID) (models.ContactEmail, error) {
	panic("unimplemented")
}

// DeletePersonPhone implements ContactDatastore.
func (cs contactStore) DeletePersonPhone(ctx context.Context, personID uuid.UUID, phoneID uuid.UUID) (models.ContactPhone, error) {
	panic("unimplemented")
}

// GetSpecificContactAddresses implements PersonStore.
func (cs contactStore) GetPersonAddresses(ctx context.Context, id uuid.UUID) ([]models.ContactAddress, error) {
	query, params, err := cs.sqlBuilder.Select(cs.tableNameMap["addresses"], "*").Where(condition.Equals("person_id", id)).Build()
	if err != nil {
		return nil, err
	}

	rows, err := cs.dbConn.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return cs.personAddressFromRows(rows)
}

// GetSpecificContacts implements PersonStore.
func (cs contactStore) GetPersonEmails(ctx context.Context, id uuid.UUID) ([]models.ContactEmail, error) {
	query, params, err := cs.sqlBuilder.Select(cs.tableNameMap["emails"], "*").Where(condition.Equals("person_id", id)).Build()
	if err != nil {
		return nil, err
	}

	rows, err := cs.dbConn.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	return cs.personEmailFromRows(rows)
}

// GetSpecificContactPhones implements PersonStore.
func (cs contactStore) GetPersonPhones(ctx context.Context, id uuid.UUID) ([]models.ContactPhone, error) {
	query, params, err := cs.sqlBuilder.Select(cs.tableNameMap["phones"], "*").Where(condition.Equals("person_id", id)).Build()
	if err != nil {
		return nil, err
	}

	rows, err := cs.dbConn.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	return cs.personPhoneFromRows(rows)
}

// UpdatePersonAddress implements ContactDatastore.
func (cs contactStore) UpdatePersonAddress(ctx context.Context, personID uuid.UUID, addressID uuid.UUID, newVal models.ContactAddress) error {
	query, params, err := cs.sqlBuilder.Update(cs.tableNameMap["addresses"]).SetStruct(newVal).Where(
		condition.Equals("id", addressID),
		condition.Equals("person_id", personID),
	).Build()
	if err != nil {
		return err
	}

	_, err = cs.dbConn.ExecContext(ctx, query, params...)
	if err != nil {
		return err
	}

	return nil
}

// UpdatePersonEmail implements ContactDatastore.
func (cs contactStore) UpdatePersonEmail(ctx context.Context, personID uuid.UUID, emailID uuid.UUID, newVal models.ContactEmail) error {
	query, params, err := cs.sqlBuilder.Update(cs.tableNameMap["emails"]).SetStruct(newVal).Where(
		condition.Equals("id", emailID),
		condition.Equals("person_id", personID),
	).Build()
	if err != nil {
		return err
	}

	_, err = cs.dbConn.ExecContext(ctx, query, params...)
	if err != nil {
		return err
	}

	return nil
}

// UpdatePersonPhone implements ContactDatastore.
func (cs contactStore) UpdatePersonPhone(ctx context.Context, personID uuid.UUID, phoneID uuid.UUID, newVal models.ContactPhone) error {
	query, params, err := cs.sqlBuilder.Update(cs.tableNameMap["phones"]).SetStruct(newVal).Where(
		condition.Equals("id", phoneID),
		condition.Equals("person_id", personID),
	).Build()
	if err != nil {
		return err
	}

	_, err = cs.dbConn.ExecContext(ctx, query, params...)
	if err != nil {
		return err
	}

	return nil
}

func (cs contactStore) personAddressFromRows(rows *sql.Rows) ([]models.ContactAddress, error) {
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

func (cs contactStore) personEmailFromRows(rows *sql.Rows) ([]models.ContactEmail, error) {
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

func (cs contactStore) personPhoneFromRows(rows *sql.Rows) ([]models.ContactPhone, error) {
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

func NewContactStore(dbConn *sql.DB) ContactDatastore {
	return contactStore{
		dbConn:     dbConn,
		sqlBuilder: jagsqlb.NewSqlBuilder(),
		tableNameMap: map[string]string{
			"addresses": "person.addresses",
			"emails":    "person.emails",
			"phones":    "person.phones",
		},
	}
}
