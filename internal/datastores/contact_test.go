package datastores

import (
	"context"
	"database/sql/driver"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/williabk198/jagsqlb"
	"github.com/williabk198/timeclock/internal/models"
)

func Test_contactStore_AddPersonAddress(t *testing.T) {
	type args struct {
		ctx     context.Context
		address models.ContactAddress
	}
	type wantQuery struct {
		rawQuery  string
		arguments []driver.Value
		result    *sqlmock.Rows
		returnErr error
	}

	testSqlBuilder := jagsqlb.NewSqlBuilder()
	mockSession, mockDb, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	testPersonID := uuid.New()
	newAddressID := uuid.New()
	tests := []struct {
		name      string
		ps        contactStore
		args      args
		wantQuery *wantQuery
		wantID    uuid.UUID
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ps: contactStore{
				dbConn:       mockSession,
				sqlBuilder:   testSqlBuilder,
				tableNameMap: map[string]string{"addresses": "person.addresses"},
			},
			args: args{
				ctx: context.Background(),
				address: models.ContactAddress{
					PersonID:   testPersonID,
					Street1:    "123 Test Dr",
					Street2:    "APT 1",
					Locality:   "Testington",
					Region:     "Testaria",
					PostalCode: "12345-6789",
					Country:    "Testopia",
					Type:       models.AddressTypePhysical,
					Primary:    true,
				},
			},
			wantQuery: &wantQuery{
				rawQuery:  `INSERT INTO "person"."addresses" ("person_id", "street1", "street2", "locality", "region", "postal_code", "country", "kind", "primary") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING "id";`,
				arguments: []driver.Value{testPersonID.String(), "123 Test Dr", "APT 1", "Testington", "Testaria", "12345-6789", "Testopia", "physical", true},
				result:    sqlmock.NewRows([]string{"id"}).AddRow(newAddressID),
			},
			wantID:    newAddressID,
			assertion: assert.NoError,
		},
		{
			name: "Error; Query Builder",
			ps: contactStore{
				dbConn:       mockSession,
				sqlBuilder:   testSqlBuilder,
				tableNameMap: map[string]string{"emails": ".bad_val"},
			},
			args: args{
				ctx:     context.Background(),
				address: models.ContactAddress{},
			},
			wantID:    uuid.Nil,
			assertion: assert.Error,
		},
		{
			name: "Error; Query Execution",
			ps: contactStore{
				dbConn:       mockSession,
				sqlBuilder:   testSqlBuilder,
				tableNameMap: map[string]string{"addresses": "person.addresses"},
			},
			args: args{
				ctx: context.Background(),
				address: models.ContactAddress{
					PersonID:   testPersonID,
					Street1:    "123 Test Dr",
					Street2:    "APT 1",
					Locality:   "Testington",
					Region:     "Testaria",
					PostalCode: "12345-6789",
					Country:    "Testopia",
					Type:       models.AddressTypePhysical,
					Primary:    true,
				},
			},
			wantQuery: &wantQuery{
				rawQuery:  `INSERT INTO "person"."addresses" ("person_id", "street1", "street2", "locality", "region", "postal_code", "country", "kind", "primary") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING "id";`,
				arguments: []driver.Value{testPersonID.String(), "123 Test Dr", "APT 1", "Testington", "Testaria", "12345-6789", "Testopia", "physical", true},
				result:    sqlmock.NewRows(nil),
				returnErr: assert.AnError,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantQuery != nil {
				mockDb.ExpectQuery(regexp.QuoteMeta(tt.wantQuery.rawQuery)).
					WithArgs(tt.wantQuery.arguments...).WillReturnRows(
					tt.wantQuery.result,
				).WillReturnError(tt.wantQuery.returnErr)
			}

			got, err := tt.ps.AddPersonAddress(tt.args.ctx, tt.args.address)
			tt.assertion(t, err)
			assert.Equal(t, tt.wantID, got)

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_contactStore_AddPersonEmail(t *testing.T) {
	type args struct {
		ctx   context.Context
		email models.ContactEmail
	}
	type wantQuery struct {
		rawQuery  string
		arguments []driver.Value
		result    *sqlmock.Rows
		returnErr error
	}

	testSqlBuilder := jagsqlb.NewSqlBuilder()
	mockSession, mockDb, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	testPersonID := uuid.New()
	newEmailID := uuid.New()

	tests := []struct {
		name      string
		ps        contactStore
		args      args
		wantQuery *wantQuery
		wantID    uuid.UUID
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ps: contactStore{
				dbConn:       mockSession,
				sqlBuilder:   testSqlBuilder,
				tableNameMap: map[string]string{"emails": "person.emails"},
			},
			args: args{
				ctx: context.Background(),
				email: models.ContactEmail{
					PersonID: testPersonID,
					Username: "test",
					Provider: "example.com",
					Primary:  true,
				},
			},
			wantQuery: &wantQuery{
				rawQuery:  `INSERT INTO "person"."emails" ("person_id", "username", "provider", "primary") VALUES ($1, $2, $3, $4) RETURNING "id";`,
				arguments: []driver.Value{testPersonID.String(), "test", "example.com", true},
				result:    sqlmock.NewRows([]string{"id"}).AddRow(newEmailID),
			},
			wantID:    newEmailID,
			assertion: assert.NoError,
		},
		{
			name: "Error; Query Builder",
			ps: contactStore{
				dbConn:       mockSession,
				sqlBuilder:   testSqlBuilder,
				tableNameMap: map[string]string{"emails": ".bad_val"},
			},
			args: args{
				ctx:   context.Background(),
				email: models.ContactEmail{},
			},
			wantID:    uuid.Nil,
			assertion: assert.Error,
		},
		{
			name: "Error; Query Execution",
			ps: contactStore{
				dbConn:       mockSession,
				sqlBuilder:   testSqlBuilder,
				tableNameMap: map[string]string{"emails": "person.emails"},
			},
			args: args{
				ctx: context.Background(),
				email: models.ContactEmail{
					PersonID: testPersonID,
					Username: "test",
					Provider: "example.com",
					Primary:  true,
				},
			},
			wantQuery: &wantQuery{
				rawQuery:  `INSERT INTO "person"."emails" ("person_id", "username", "provider", "primary") VALUES ($1, $2, $3, $4) RETURNING "id";`,
				arguments: []driver.Value{testPersonID.String(), "test", "example.com", true},
				result:    sqlmock.NewRows(nil),
				returnErr: assert.AnError,
			},
			wantID:    uuid.Nil,
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantQuery != nil {
				mockDb.ExpectQuery(regexp.QuoteMeta(tt.wantQuery.rawQuery)).
					WithArgs(tt.wantQuery.arguments...).WillReturnRows(
					tt.wantQuery.result,
				).WillReturnError(tt.wantQuery.returnErr)
			}

			got, err := tt.ps.AddPersonEmail(tt.args.ctx, tt.args.email)
			tt.assertion(t, err)
			assert.Equal(t, tt.wantID, got)

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_contactStore_AddPersonPhone(t *testing.T) {
	type args struct {
		ctx   context.Context
		phone models.ContactPhone
	}
	type wantQuery struct {
		rawQuery  string
		arguments []driver.Value
		result    *sqlmock.Rows
		returnErr error
	}

	testSqlBuilder := jagsqlb.NewSqlBuilder()
	mockSession, mockDb, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	testPersonID := uuid.New()
	newPhoneID := uuid.New()

	tests := []struct {
		name      string
		ps        contactStore
		args      args
		wantQuery *wantQuery
		wantID    uuid.UUID
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ps: contactStore{
				dbConn:       mockSession,
				sqlBuilder:   testSqlBuilder,
				tableNameMap: map[string]string{"phones": "person.phones"},
			},
			args: args{
				ctx: context.Background(),
				phone: models.ContactPhone{
					PersonID:    testPersonID,
					CountryCode: 1,
					PhoneNumber: "(555) 555-5555",
					Type:        models.PhoneTypeHome,
					Primary:     true,
				},
			},
			wantQuery: &wantQuery{
				rawQuery:  `INSERT INTO "person"."phones" ("person_id", "country_code", "phone_number", "kind", "primary") VALUES ($1, $2, $3, $4, $5) RETURNING "id";`,
				arguments: []driver.Value{testPersonID.String(), 1, "(555) 555-5555", "home", true},
				result:    sqlmock.NewRows([]string{"id"}).AddRow(newPhoneID),
			},
			wantID:    newPhoneID,
			assertion: assert.NoError,
		},
		{
			name: "Error; Query Builder",
			ps: contactStore{
				dbConn:       mockSession,
				sqlBuilder:   testSqlBuilder,
				tableNameMap: map[string]string{"emails": ".bad_val"},
			},
			args: args{
				ctx:   context.Background(),
				phone: models.ContactPhone{},
			},
			wantID:    uuid.Nil,
			assertion: assert.Error,
		},
		{
			name: "Error; Query Execution",
			ps: contactStore{
				dbConn:       mockSession,
				sqlBuilder:   testSqlBuilder,
				tableNameMap: map[string]string{"phones": "person.phones"},
			},
			args: args{
				ctx: context.Background(),
				phone: models.ContactPhone{
					PersonID:    testPersonID,
					CountryCode: 1,
					PhoneNumber: "(555) 555-5555",
					Type:        models.PhoneTypeHome,
					Primary:     true,
				},
			},
			wantQuery: &wantQuery{
				rawQuery:  `INSERT INTO "person"."phones" ("person_id", "country_code", "phone_number", "kind", "primary") VALUES ($1, $2, $3, $4, $5) RETURNING "id";`,
				arguments: []driver.Value{testPersonID.String(), 1, "(555) 555-5555", "home", true},
				result:    sqlmock.NewRows(nil),
				returnErr: assert.AnError,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		if tt.wantQuery != nil {
			mockDb.ExpectQuery(regexp.QuoteMeta(tt.wantQuery.rawQuery)).
				WithArgs(tt.wantQuery.arguments...).WillReturnRows(
				tt.wantQuery.result,
			).WillReturnError(tt.wantQuery.returnErr)
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ps.AddPersonPhone(tt.args.ctx, tt.args.phone)
			tt.assertion(t, err)
			assert.Equal(t, tt.wantID, got)
		})

		// Check to see if the expectations of the query were met
		if err := mockDb.ExpectationsWereMet(); err != nil {
			t.Errorf("sql expectations were not met: %v", err)
		}
	}
}

func Test_contactStore_DeletePersonAddress(t *testing.T) {
	type args struct {
		ctx       context.Context
		personID  uuid.UUID
		addressID uuid.UUID
	}
	type wantQuery struct {
		rawQuery  string
		arguments []driver.Value
		result    *sqlmock.Rows
		returnErr error
	}

	tableRows := []string{"id", "person_id", "street1", "street2", "locality", "region", "postal_code", "country", "kind", "primary"}
	testSqlBuilder := jagsqlb.NewSqlBuilder()
	mockSession, mockDb, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	testAddressID := uuid.New()
	testPersonID := uuid.New()

	testNotFoundAddressID := uuid.New()

	testAddressDB := models.ContactAddress{
		ID:         testAddressID,
		PersonID:   testPersonID,
		Street1:    "123 Test Dr",
		Street2:    "APT 1",
		Locality:   "Testington",
		Region:     "Testaria",
		PostalCode: "12345-6789",
		Country:    "Testopia",
		Type:       models.AddressTypePhysical,
		Primary:    true,
	}

	tests := []struct {
		name      string
		cs        contactStore
		args      args
		wantQuery *wantQuery
		want      models.ContactAddress
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			cs: contactStore{
				sqlBuilder:   testSqlBuilder,
				dbConn:       mockSession,
				tableNameMap: map[string]string{"addresses": "person.addresses"},
			},
			args: args{
				ctx:       context.Background(),
				personID:  testPersonID,
				addressID: testAddressID,
			},
			wantQuery: &wantQuery{
				rawQuery:  `DELETE FROM "person"."addresses" WHERE "id" = $1 AND "person_id" = $2 RETURNING *`,
				arguments: []driver.Value{testAddressID.String(), testPersonID.String()},
				result: sqlmock.NewRows(tableRows).AddRow(
					testAddressID.String(), testPersonID.String(), "123 Test Dr", "APT 1",
					"Testington", "Testaria", "12345-6789", "Testopia", "physical", true,
				),
			},
			want:      testAddressDB,
			assertion: assert.NoError,
		},
		{
			name: "Error; Query Builder",
			cs: contactStore{
				sqlBuilder:   testSqlBuilder,
				dbConn:       mockSession,
				tableNameMap: map[string]string{"addresses": ".bad_val"},
			},
			args: args{
				ctx:       context.Background(),
				personID:  testPersonID,
				addressID: testAddressID,
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Query Execution",
			cs: contactStore{
				sqlBuilder:   testSqlBuilder,
				dbConn:       mockSession,
				tableNameMap: map[string]string{"addresses": "person.addresses"},
			},
			args: args{
				ctx:       context.Background(),
				personID:  testPersonID,
				addressID: testNotFoundAddressID,
			},
			wantQuery: &wantQuery{
				rawQuery:  `DELETE FROM "person"."addresses" WHERE "id" = $1 AND "person_id" = $2 RETURNING *`,
				arguments: []driver.Value{testNotFoundAddressID.String(), testPersonID.String()},
				result:    sqlmock.NewRows(nil),
				returnErr: assert.AnError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantQuery != nil {
				mockDb.ExpectQuery(regexp.QuoteMeta(tt.wantQuery.rawQuery)).
					WithArgs(tt.wantQuery.arguments...).WillReturnRows(
					tt.wantQuery.result,
				).WillReturnError(tt.wantQuery.returnErr)
			}

			got, err := tt.cs.DeletePersonAddress(tt.args.ctx, tt.args.personID, tt.args.addressID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_contactStore_DeletePersonEmail(t *testing.T) {
	type args struct {
		ctx      context.Context
		personID uuid.UUID
		emailID  uuid.UUID
	}
	type wantQuery struct {
		rawQuery  string
		arguments []driver.Value
		result    *sqlmock.Rows
		returnErr error
	}

	tableRows := []string{"id", "person_id", "username", "provider", "primary"}
	testSqlBuilder := jagsqlb.NewSqlBuilder()
	mockSession, mockDb, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	testPersonID := uuid.New()
	testEmailID := uuid.New()

	testNotFoundEmailID := uuid.New()

	testEmailDB := models.ContactEmail{
		ID:       testEmailID,
		PersonID: testPersonID,
		Username: "tester",
		Provider: "test.com",
		Primary:  true,
	}

	tests := []struct {
		name      string
		cs        contactStore
		args      args
		wantQuery *wantQuery
		want      models.ContactEmail
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			cs: contactStore{
				sqlBuilder:   testSqlBuilder,
				dbConn:       mockSession,
				tableNameMap: map[string]string{"emails": "person.emails"},
			},
			args: args{
				ctx:      context.Background(),
				personID: testPersonID,
				emailID:  testEmailID,
			},
			wantQuery: &wantQuery{
				rawQuery:  `DELETE FROM "person"."emails" WHERE "id" = $1 AND "person_id" = $2 RETURNING *`,
				arguments: []driver.Value{testEmailID.String(), testPersonID.String()},
				result:    sqlmock.NewRows(tableRows).AddRow(testEmailID.String(), testPersonID.String(), "tester", "test.com", true),
			},
			want:      testEmailDB,
			assertion: assert.NoError,
		},
		{
			name: "Error; Query Builder",
			cs: contactStore{
				sqlBuilder:   testSqlBuilder,
				dbConn:       mockSession,
				tableNameMap: map[string]string{"emails": ".bad_val"},
			},
			args: args{
				ctx:      context.Background(),
				personID: testPersonID,
				emailID:  testEmailID,
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Query Execution",
			cs: contactStore{
				sqlBuilder:   testSqlBuilder,
				dbConn:       mockSession,
				tableNameMap: map[string]string{"emails": "person.emails"},
			},
			args: args{
				ctx:      context.Background(),
				personID: testPersonID,
				emailID:  testNotFoundEmailID,
			},
			wantQuery: &wantQuery{
				rawQuery:  `DELETE FROM "person"."emails" WHERE "id" = $1 AND "person_id" = $2 RETURNING *`,
				arguments: []driver.Value{testNotFoundEmailID.String(), testPersonID.String()},
				result:    sqlmock.NewRows(nil),
				returnErr: assert.AnError,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantQuery != nil {
				mockDb.ExpectQuery(regexp.QuoteMeta(tt.wantQuery.rawQuery)).
					WithArgs(tt.wantQuery.arguments...).WillReturnRows(
					tt.wantQuery.result,
				).WillReturnError(tt.wantQuery.returnErr)
			}

			got, err := tt.cs.DeletePersonEmail(tt.args.ctx, tt.args.personID, tt.args.emailID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_contactStore_DeletePersonPhone(t *testing.T) {
	type args struct {
		ctx      context.Context
		personID uuid.UUID
		phoneID  uuid.UUID
	}
	type wantQuery struct {
		rawQuery  string
		arguments []driver.Value
		result    *sqlmock.Rows
		returnErr error
	}

	tableRows := []string{"id", "person_id", "country_code", "phone_number", "kind", "primary"}
	testSqlBuilder := jagsqlb.NewSqlBuilder()
	mockSession, mockDb, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	testPersonID := uuid.New()
	testPhoneID := uuid.New()

	testNotFoundPhoneID := uuid.New()

	testPhoneDB := models.ContactPhone{
		ID:          testPhoneID,
		PersonID:    testPersonID,
		CountryCode: 1,
		PhoneNumber: "(555)555-5555",
		Type:        models.PhoneTypeHome,
		Primary:     true,
	}

	tests := []struct {
		name      string
		cs        contactStore
		args      args
		wantQuery *wantQuery
		want      models.ContactPhone
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			cs: contactStore{
				sqlBuilder:   testSqlBuilder,
				dbConn:       mockSession,
				tableNameMap: map[string]string{"phones": "person.phones"},
			},
			args: args{
				ctx:      context.Background(),
				personID: testPersonID,
				phoneID:  testPhoneID,
			},
			wantQuery: &wantQuery{
				rawQuery:  `DELETE FROM "person"."phones" WHERE "id" = $1 AND "person_id" = $2`,
				arguments: []driver.Value{testPhoneID.String(), testPersonID.String()},
				result:    sqlmock.NewRows(tableRows).AddRow(testPhoneID.String(), testPersonID.String(), 1, "(555)555-5555", "home", true),
			},
			want:      testPhoneDB,
			assertion: assert.NoError,
		},
		{
			name: "Error; Query Builder",
			cs: contactStore{
				sqlBuilder:   testSqlBuilder,
				dbConn:       mockSession,
				tableNameMap: map[string]string{"phones": ".bad_val"},
			},
			args: args{
				ctx:      context.Background(),
				personID: testPersonID,
				phoneID:  testPhoneID,
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Query Execution",
			cs: contactStore{
				sqlBuilder:   testSqlBuilder,
				dbConn:       mockSession,
				tableNameMap: map[string]string{"phones": "person.phones"},
			},
			args: args{
				ctx:      context.Background(),
				personID: testPersonID,
				phoneID:  testNotFoundPhoneID,
			},
			wantQuery: &wantQuery{
				rawQuery:  `DELETE FROM "person"."phones" WHERE "id" = $1 AND "person_id" = $2`,
				arguments: []driver.Value{testNotFoundPhoneID.String(), testPersonID.String()},
				result:    sqlmock.NewRows(nil),
				returnErr: assert.AnError,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantQuery != nil {
				mockDb.ExpectQuery(regexp.QuoteMeta(tt.wantQuery.rawQuery)).
					WithArgs(tt.wantQuery.arguments...).WillReturnRows(
					tt.wantQuery.result,
				).WillReturnError(tt.wantQuery.returnErr)
			}

			got, err := tt.cs.DeletePersonPhone(tt.args.ctx, tt.args.personID, tt.args.phoneID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_contactStore_GetPersonAddresses(t *testing.T) {
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}

	type wantQuery struct {
		rawQuery  string
		arguments []driver.Value
		result    *sqlmock.Rows
		returnErr error
	}

	addressColumns := []string{"id", "person_id", "street1", "street2", "locality", "region", "postal_code", "country", "kind", "primary"}
	testSqlBuilder := jagsqlb.NewSqlBuilder()
	mockSession, mockDb, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	testPersonStore := contactStore{
		dbConn:     mockSession,
		sqlBuilder: testSqlBuilder,
		tableNameMap: map[string]string{
			"addresses": "person.addresses",
		},
	}
	testPersonID := uuid.New()
	testNotFoundID := uuid.New()
	testAddressID1 := uuid.New()
	testAddressID2 := uuid.New()

	tests := []struct {
		name      string
		ps        contactStore
		args      args
		wantQuery *wantQuery
		want      []models.ContactAddress
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ps:   testPersonStore,
			args: args{
				ctx: context.Background(),
				id:  testPersonID,
			},
			wantQuery: &wantQuery{
				rawQuery:  `SELECT * FROM "person"."addresses" WHERE "person_id" = $1`,
				arguments: []driver.Value{testPersonID.String()},
				result: sqlmock.NewRows(addressColumns).AddRows(
					[]driver.Value{
						testAddressID1, testPersonID, "123 Test Dr", "APT 1", "Testerville", "Testton", "12345-6789", "U.S.", "physical", true,
					},
					[]driver.Value{
						testAddressID2, testPersonID, "P.O. Box 123", "", "Testerville", "Testton", "12344-5567", "U.S.", "mailing", true,
					},
				),
			},
			want: []models.ContactAddress{
				{
					ID:         testAddressID1,
					PersonID:   testPersonID,
					Street1:    "123 Test Dr",
					Street2:    "APT 1",
					Locality:   "Testerville",
					Region:     "Testton",
					PostalCode: "12345-6789",
					Country:    "U.S.",
					Type:       "physical",
					Primary:    true,
				},
				{
					ID:         testAddressID2,
					PersonID:   testPersonID,
					Street1:    "P.O. Box 123",
					Locality:   "Testerville",
					Region:     "Testton",
					PostalCode: "12344-5567",
					Country:    "U.S.",
					Type:       "mailing",
					Primary:    true,
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; SQL Builder",
			ps: contactStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableNameMap: map[string]string{
					"emails": ".bad_val",
				},
			},
			args: args{
				ctx: context.Background(),
				id:  testPersonID,
			},
			assertion: assert.Error,
		},
		{
			name: "Error; SQL Execution",
			ps:   testPersonStore,
			args: args{
				ctx: context.Background(),
				id:  testNotFoundID,
			},
			wantQuery: &wantQuery{
				rawQuery:  `SELECT * FROM "person"."addresses" WHERE "person_id" = $1`,
				arguments: []driver.Value{testNotFoundID.String()},
				result:    sqlmock.NewRows(nil),
				returnErr: assert.AnError,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantQuery != nil {
				mockDb.ExpectQuery(regexp.QuoteMeta(tt.wantQuery.rawQuery)).
					WithArgs(tt.wantQuery.arguments...).WillReturnRows(
					tt.wantQuery.result,
				).WillReturnError(tt.wantQuery.returnErr)
			}

			got, err := tt.ps.GetPersonAddresses(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_contactStore_GetPersonEmails(t *testing.T) {
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}

	type wantQuery struct {
		rawQuery  string
		arguments []driver.Value
		result    *sqlmock.Rows
		returnErr error
	}

	emailColumns := []string{"id", "person_id", "username", "provider", "primary"}
	testSqlBuilder := jagsqlb.NewSqlBuilder()
	mockSession, mockDb, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	testPersonID := uuid.New()
	testNotFoundID := uuid.New()
	testEmailID1 := uuid.New()
	testEmailID2 := uuid.New()

	testPersonStore := contactStore{
		dbConn:     mockSession,
		sqlBuilder: testSqlBuilder,
		tableNameMap: map[string]string{
			"emails": "person.emails",
		},
	}

	tests := []struct {
		name      string
		ps        contactStore
		args      args
		wantQuery *wantQuery
		want      []models.ContactEmail
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ps:   testPersonStore,
			args: args{
				ctx: context.Background(),
				id:  testPersonID,
			},
			wantQuery: &wantQuery{
				rawQuery:  `SELECT * FROM "person"."emails" WHERE "person_id" = $1`,
				arguments: []driver.Value{testPersonID.String()},
				result: sqlmock.NewRows(emailColumns).AddRows(
					[]driver.Value{
						testEmailID1.String(), testPersonID.String(), "tester", "example.com", true,
					},
					[]driver.Value{
						testEmailID2.String(), testPersonID.String(), "admin", "example.com", false,
					},
				),
			},
			want: []models.ContactEmail{
				{
					ID:       testEmailID1,
					PersonID: testPersonID,
					Username: "tester",
					Provider: "example.com",
					Primary:  true,
				},
				{
					ID:       testEmailID2,
					PersonID: testPersonID,
					Username: "admin",
					Provider: "example.com",
					Primary:  false,
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; Query Builder",
			ps: contactStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableNameMap: map[string]string{
					"emails": ".bad_val",
				},
			},
			args: args{
				ctx: context.Background(),
				id:  testPersonID,
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Query Execution",
			ps:   testPersonStore,
			args: args{
				ctx: context.Background(),
				id:  testNotFoundID,
			},
			wantQuery: &wantQuery{
				rawQuery:  `SELECT * FROM "person"."emails" WHERE "person_id" = $1`,
				arguments: []driver.Value{testNotFoundID.String()},
				result:    sqlmock.NewRows(nil),
				returnErr: assert.AnError,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantQuery != nil {
				mockDb.ExpectQuery(regexp.QuoteMeta(tt.wantQuery.rawQuery)).
					WithArgs(tt.wantQuery.arguments...).WillReturnRows(
					tt.wantQuery.result,
				).WillReturnError(tt.wantQuery.returnErr)
			}

			got, err := tt.ps.GetPersonEmails(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_contactStore_GetPersonPhones(t *testing.T) {
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}

	type wantQuery struct {
		rawQuery  string
		arguments []driver.Value
		result    *sqlmock.Rows
		returnErr error
	}

	phoneColumns := []string{"id", "person_id", "country_code", "phone_number", "kind", "primary"}
	testSqlBuilder := jagsqlb.NewSqlBuilder()
	mockSession, mockDb, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	testPersonID := uuid.New()
	testNotFoundID := uuid.New()
	testPhoneID1 := uuid.New()
	testPhoneID2 := uuid.New()

	testPersonStore := contactStore{
		dbConn:     mockSession,
		sqlBuilder: testSqlBuilder,
		tableNameMap: map[string]string{
			"phones": "person.phones",
		},
	}

	tests := []struct {
		name      string
		ps        contactStore
		args      args
		wantQuery *wantQuery
		want      []models.ContactPhone
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ps:   testPersonStore,
			args: args{
				ctx: context.Background(),
				id:  testPersonID,
			},
			wantQuery: &wantQuery{
				rawQuery:  `SELECT * FROM "person"."phones" WHERE "person_id" = $1`,
				arguments: []driver.Value{testPersonID.String()},
				result: sqlmock.NewRows(phoneColumns).AddRows(
					[]driver.Value{
						testPhoneID1.String(), testPersonID.String(), 1, "(315)555-1234", "home", true,
					},
					[]driver.Value{
						testPhoneID2.String(), testPersonID.String(), 1, "(315)555-9876", "cell", false,
					},
				),
			},
			want: []models.ContactPhone{
				{
					ID:          testPhoneID1,
					PersonID:    testPersonID,
					CountryCode: 1,
					PhoneNumber: "(315)555-1234",
					Type:        "home",
					Primary:     true,
				},
				{
					ID:          testPhoneID2,
					PersonID:    testPersonID,
					CountryCode: 1,
					PhoneNumber: "(315)555-9876",
					Type:        "cell",
					Primary:     false,
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; SQL Builder",
			ps: contactStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableNameMap: map[string]string{
					"phone": ".bad_val",
				},
			},
			args: args{
				ctx: context.Background(),
				id:  testPersonID,
			},
			assertion: assert.Error,
		},
		{
			name: "Error; SQL Execution",
			ps:   testPersonStore,
			args: args{
				ctx: context.Background(),
				id:  testNotFoundID,
			},
			wantQuery: &wantQuery{
				rawQuery:  `SELECT * FROM "person"."phones" WHERE "person_id" = $1`,
				arguments: []driver.Value{testNotFoundID.String()},
				result:    sqlmock.NewRows(nil),
				returnErr: assert.AnError,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantQuery != nil {
				mockDb.ExpectQuery(regexp.QuoteMeta(tt.wantQuery.rawQuery)).
					WithArgs(tt.wantQuery.arguments...).WillReturnRows(
					tt.wantQuery.result,
				).WillReturnError(tt.wantQuery.returnErr)
			}

			got, err := tt.ps.GetPersonPhones(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_contactStore_UpdatePersonAddress(t *testing.T) {
	type args struct {
		ctx       context.Context
		personID  uuid.UUID
		addressID uuid.UUID
		newVal    models.ContactAddress
	}

	type wantQuery struct {
		rawQuery  string
		arguments []driver.Value
		result    driver.Result
		returnErr error
	}

	testSqlBuilder := jagsqlb.NewSqlBuilder()
	mockSession, mockDb, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	testPersonID := uuid.New()
	testAddressID := uuid.New()

	tests := []struct {
		name      string
		cs        contactStore
		args      args
		wantQuery *wantQuery
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			cs: contactStore{
				sqlBuilder:   testSqlBuilder,
				dbConn:       mockSession,
				tableNameMap: map[string]string{"addresses": "person.addresses"},
			},
			args: args{
				ctx:       context.Background(),
				personID:  testPersonID,
				addressID: testAddressID,
				newVal: models.ContactAddress{
					ID:         testAddressID,
					PersonID:   testAddressID,
					Street1:    "987 Testing Ln",
					Street2:    "APT 3",
					Locality:   "Testington",
					Region:     "Testoria",
					PostalCode: "98765-4321",
					Country:    "Testopia",
					Type:       models.AddressTypePhysical,
					Primary:    true,
				},
			},
			wantQuery: &wantQuery{
				rawQuery:  `UPDATE "person"."addresses" SET "street1"=$1, "street2"=$2, "locality"=$3, "region"=$4, "postal_code"=$5, "country"=$6, "kind"=$7, "primary"=$8 WHERE "id" = $9 AND "person_id" = $10;`,
				arguments: []driver.Value{"987 Testing Ln", "APT 3", "Testington", "Testoria", "98765-4321", "Testopia", "physical", true, testAddressID.String(), testPersonID.String()},
				result:    sqlmock.NewResult(0, 1),
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; SQL Builder",
			cs: contactStore{
				sqlBuilder:   testSqlBuilder,
				dbConn:       mockSession,
				tableNameMap: map[string]string{"addresses": ".bad_val"},
			},
			args: args{
				ctx:       context.Background(),
				personID:  testPersonID,
				addressID: testAddressID,
			},
			assertion: assert.Error,
		},
		{
			name: "Error; SQL Execution",
			cs: contactStore{
				sqlBuilder:   testSqlBuilder,
				dbConn:       mockSession,
				tableNameMap: map[string]string{"addresses": "person.addresses"},
			},
			args: args{
				ctx:       context.Background(),
				personID:  testPersonID,
				addressID: testAddressID,
				newVal: models.ContactAddress{
					ID:         testAddressID,
					PersonID:   testAddressID,
					Street1:    "987 Testing Ln",
					Street2:    "APT 3",
					Locality:   "Testington",
					Region:     "Testoria",
					PostalCode: "98765-4321",
					Country:    "Testopia",
					Type:       models.AddressTypePhysical,
					Primary:    true,
				},
			},
			wantQuery: &wantQuery{
				rawQuery:  `UPDATE "person"."addresses" SET "street1"=$1, "street2"=$2, "locality"=$3, "region"=$4, "postal_code"=$5, "country"=$6, "kind"=$7, "primary"=$8 WHERE "id" = $9 AND "person_id" = $10;`,
				arguments: []driver.Value{"987 Testing Ln", "APT 3", "Testington", "Testoria", "98765-4321", "Testopia", "physical", true, testAddressID.String(), testPersonID.String()},
				returnErr: assert.AnError,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantQuery != nil {
				mockDb.ExpectExec(regexp.QuoteMeta(tt.wantQuery.rawQuery)).
					WithArgs(tt.wantQuery.arguments...).
					WillReturnResult(tt.wantQuery.result).
					WillReturnError(tt.wantQuery.returnErr)
			}

			tt.assertion(t, tt.cs.UpdatePersonAddress(tt.args.ctx, tt.args.personID, tt.args.addressID, tt.args.newVal))

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_contactStore_UpdatePersonEmail(t *testing.T) {
	type args struct {
		ctx      context.Context
		personID uuid.UUID
		emailID  uuid.UUID
		newVal   models.ContactEmail
	}

	type wantQuery struct {
		rawQuery  string
		arguments []driver.Value
		result    driver.Result
		returnErr error
	}

	testSqlBuilder := jagsqlb.NewSqlBuilder()
	mockSession, mockDb, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	testPersonID := uuid.New()
	testEmailID := uuid.New()

	tests := []struct {
		name      string
		cs        contactStore
		args      args
		wantQuery *wantQuery
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			cs: contactStore{
				sqlBuilder:   testSqlBuilder,
				dbConn:       mockSession,
				tableNameMap: map[string]string{"emails": "person.emails"},
			},
			args: args{
				ctx:      context.Background(),
				personID: testPersonID,
				emailID:  testEmailID,
				newVal: models.ContactEmail{
					Username: "tester",
					Provider: "test.com",
					Primary:  true,
				},
			},
			wantQuery: &wantQuery{
				rawQuery:  `UPDATE "person"."emails" SET "username"=$1, "provider"=$2, "primary"=$3 WHERE "id" = $4 AND "person_id" = $5;`,
				arguments: []driver.Value{"tester", "test.com", true, testEmailID.String(), testPersonID.String()},
				result:    sqlmock.NewResult(0, 1),
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; SQL Builder",
			cs: contactStore{
				sqlBuilder:   testSqlBuilder,
				dbConn:       mockSession,
				tableNameMap: map[string]string{"emails": ".bad_val"},
			},
			args: args{
				ctx:      context.Background(),
				personID: testPersonID,
				emailID:  testEmailID,
				newVal: models.ContactEmail{
					Username: "tester",
					Provider: "test.com",
					Primary:  true,
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; SQL Execution",
			cs: contactStore{
				sqlBuilder:   testSqlBuilder,
				dbConn:       mockSession,
				tableNameMap: map[string]string{"emails": "person.emails"},
			},
			args: args{
				ctx:      context.Background(),
				personID: testPersonID,
				emailID:  testEmailID,
				newVal: models.ContactEmail{
					Username: "tester",
					Provider: "test.com",
					Primary:  true,
				},
			},
			wantQuery: &wantQuery{
				rawQuery:  `UPDATE "person"."emails" SET "username"=$1, "provider"=$2, "primary"=$3 WHERE "id" = $4 AND "person_id" = $5;`,
				arguments: []driver.Value{"tester", "test.com", true, testEmailID.String(), testPersonID.String()},
				returnErr: assert.AnError,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantQuery != nil {
				mockDb.ExpectExec(regexp.QuoteMeta(tt.wantQuery.rawQuery)).
					WithArgs(tt.wantQuery.arguments...).
					WillReturnResult(tt.wantQuery.result).
					WillReturnError(tt.wantQuery.returnErr)
			}

			tt.assertion(t, tt.cs.UpdatePersonEmail(tt.args.ctx, tt.args.personID, tt.args.emailID, tt.args.newVal))

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_contactStore_UpdatePersonPhone(t *testing.T) {
	type args struct {
		ctx      context.Context
		personID uuid.UUID
		phoneID  uuid.UUID
		newVal   models.ContactPhone
	}

	type wantQuery struct {
		rawQuery  string
		arguments []driver.Value
		result    driver.Result
		returnErr error
	}

	testSqlBuilder := jagsqlb.NewSqlBuilder()
	mockSession, mockDb, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	testPersonID := uuid.New()
	testPhoneID := uuid.New()

	tests := []struct {
		name      string
		cs        contactStore
		args      args
		wantQuery *wantQuery
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			cs: contactStore{
				sqlBuilder:   testSqlBuilder,
				dbConn:       mockSession,
				tableNameMap: map[string]string{"phones": "person.phones"},
			},
			args: args{
				ctx:      context.Background(),
				personID: testPersonID,
				phoneID:  testPhoneID,
				newVal: models.ContactPhone{
					CountryCode: 1,
					PhoneNumber: "(123)456-7890",
					Type:        models.PhoneTypeHome,
					Primary:     false,
				},
			},
			wantQuery: &wantQuery{
				rawQuery:  `UPDATE "person"."phones" SET "country_code"=$1, "phone_number"=$2, "kind"=$3, "primary"=$4 WHERE "id" = $5 AND "person_id" = $6`,
				arguments: []driver.Value{1, "(123)456-7890", "home", false, testPhoneID.String(), testPersonID.String()},
				result:    sqlmock.NewResult(0, 1),
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; SQL Builder",
			cs: contactStore{
				sqlBuilder:   testSqlBuilder,
				dbConn:       mockSession,
				tableNameMap: map[string]string{"addresses": ".bad_val"},
			},
			args: args{
				ctx:      context.Background(),
				personID: testPersonID,
				phoneID:  testPhoneID,
			},
			assertion: assert.Error,
		},
		{
			name: "Error; SQL Execution",
			cs: contactStore{
				sqlBuilder:   testSqlBuilder,
				dbConn:       mockSession,
				tableNameMap: map[string]string{"phones": "person.phones"},
			},
			args: args{
				ctx:      context.Background(),
				personID: testPersonID,
				phoneID:  testPhoneID,
				newVal: models.ContactPhone{
					CountryCode: 1,
					PhoneNumber: "(123)456-7890",
					Type:        models.PhoneTypeHome,
					Primary:     false,
				},
			},
			wantQuery: &wantQuery{
				rawQuery:  `UPDATE "person"."phones" SET "country_code"=$1, "phone_number"=$2, "kind"=$3, "primary"=$4 WHERE "id" = $5 AND "person_id" = $6`,
				arguments: []driver.Value{1, "(123)456-7890", "home", false, testPhoneID.String(), testPersonID.String()},
				returnErr: assert.AnError,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantQuery != nil {
				mockDb.ExpectExec(regexp.QuoteMeta(tt.wantQuery.rawQuery)).
					WithArgs(tt.wantQuery.arguments...).
					WillReturnResult(tt.wantQuery.result).
					WillReturnError(tt.wantQuery.returnErr)
			}

			tt.assertion(t, tt.cs.UpdatePersonPhone(tt.args.ctx, tt.args.personID, tt.args.phoneID, tt.args.newVal))

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}
