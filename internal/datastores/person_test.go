package datastores

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/williabk198/jagsqlb"
	"github.com/williabk198/timeclock/internal/models"
)

func Test_personStore_Add(t *testing.T) {
	type args struct {
		ctx  context.Context
		item models.Person
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

	newPersonID := uuid.New()

	tests := []struct {
		name      string
		ps        personStore
		args      args
		wantQuery *wantQuery
		wantId    uuid.UUID
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ps: personStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "person.persons",
			},
			args: args{
				ctx: context.Background(),
				item: models.Person{
					Name: models.Name{
						GivenName:       "Testy",
						FamilyName:      "McTesterson",
						FamilyNameFirst: models.FirstNameGiven,
					},
					DateOfBirth: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
					Gender:      models.GenderNonBinary,
					Pronouns:    models.Pronouns{Subject: "they", Object: "them"},
				},
			},
			wantQuery: &wantQuery{
				rawQuery: `INSERT INTO "person"."persons" ("given_name", "family_name", "first_name", "dob", "gender", "pronouns") VALUES ($1, $2, $3, $4, $5, $6) RETURNING "id"`,
				arguments: []driver.Value{
					"Testy", "McTesterson", "given", time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), "non-binary", "they/them",
				},
				result: mockDb.NewRows([]string{"id"}).AddRow(newPersonID),
			},
			wantId:    newPersonID,
			assertion: assert.NoError,
		},
		{
			name: "Error; Query Builder",
			ps: personStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  ".invalid",
			},
			args:      args{ctx: context.Background(), item: models.Person{}},
			assertion: assert.Error,
		},
		{
			name: "Error; Query Execution",
			ps: personStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "person.persons",
			},
			args: args{
				ctx: context.Background(),
				item: models.Person{
					Name: models.Name{
						GivenName:       "Testy",
						FamilyName:      "McTesterson",
						FamilyNameFirst: models.FirstNameGiven,
					},
					DateOfBirth: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
					Gender:      models.GenderNonBinary,
					Pronouns:    models.Pronouns{Subject: "they", Object: "them"},
				},
			},
			wantQuery: &wantQuery{
				rawQuery: `INSERT INTO "person"."persons" ("given_name", "family_name", "first_name", "dob", "gender", "pronouns") VALUES ($1, $2, $3, $4, $5, $6) RETURNING "id"`,
				arguments: []driver.Value{
					"Testy", "McTesterson", "given", time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), "non-binary", "they/them",
				},
				result:    mockDb.NewRows(nil),
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

			gotId, err := tt.ps.Add(tt.args.ctx, tt.args.item)
			tt.assertion(t, err)
			assert.Equal(t, tt.wantId, gotId)

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_personStore_AddSpecificContactEmail(t *testing.T) {
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
		ps        personStore
		args      args
		wantQuery *wantQuery
		wantID    uuid.UUID
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ps: personStore{
				dbConn:       mockSession,
				sqlBuilder:   testSqlBuilder,
				tableNameMap: map[string]string{"contact": "person.contacts"},
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
				rawQuery:  `INSERT INTO "person"."contacts"("perosn_id", "username", "provider", "primary") VALUES ($1, $2, $3, $4) RETURNING "id";`,
				arguments: []driver.Value{testPersonID.String(), "test", "example.com", true},
				result:    sqlmock.NewRows([]string{"id"}).AddRow(newEmailID),
			},
			wantID:    newEmailID,
			assertion: assert.NoError,
		},
		{
			name: "Error; Query Builder",
			ps: personStore{
				dbConn:       mockSession,
				sqlBuilder:   testSqlBuilder,
				tableNameMap: map[string]string{"contact": ".bad_val"},
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
			ps: personStore{
				dbConn:       mockSession,
				sqlBuilder:   testSqlBuilder,
				tableNameMap: map[string]string{"contact": "person.contacts"},
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
				rawQuery:  `INSERT INTO "person"."contacts"("perosn_id", "username", "provider", "primary") VALUES ($1, $2, $3, $4) RETURNING "id";`,
				arguments: []driver.Value{testPersonID.String(), "test", "example.com", true},
				result:    sqlmock.NewRows(nil),
				returnErr: assert.AnError,
			},
			wantID:    uuid.Nil,
			assertion: assert.NoError,
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

			got, err := tt.ps.AddSpecificContactEmail(tt.args.ctx, tt.args.email)
			tt.assertion(t, err)
			assert.Equal(t, tt.wantID, got)

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_personStore_Delete(t *testing.T) {
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

	tableRows := []string{"id", "given_name", "family_name", "first_name", "dob", "gender", "pronouns"}
	testSqlBuilder := jagsqlb.NewSqlBuilder()
	mockSession, mockDb, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	testPersonID := uuid.New()
	testPerson := models.Person{
		ID: testPersonID,
		Name: models.Name{
			GivenName:       "Testy",
			FamilyName:      "McTesterson",
			FamilyNameFirst: models.FirstNameGiven,
		},
		DateOfBirth: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		Gender:      models.GenderNonBinary,
		Pronouns: models.Pronouns{
			Subject: "they",
			Object:  "them",
		},
	}

	tests := []struct {
		name      string
		ps        personStore
		args      args
		wantQuery *wantQuery
		wantItem  models.Person
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ps: personStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "person.persons",
			},
			args: args{
				ctx: context.Background(),
				id:  testPersonID,
			},
			wantQuery: &wantQuery{
				rawQuery:  `DELETE FROM "person"."persons" WHERE "id" = $1 RETURNING *`,
				arguments: []driver.Value{testPersonID},
				result: sqlmock.NewRows(tableRows).AddRow(
					testPersonID.String(), "Testy", "McTesterson", "given", time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), "non-binary", "they/them",
				),
			},
			wantItem:  testPerson,
			assertion: assert.NoError,
		},
		{
			name: "Error; Query Builder",
			ps: personStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  ".invalid_name",
			},
			args: args{
				ctx: context.Background(),
				id:  testPersonID,
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Query Execution",
			ps: personStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "person.persons",
			},
			args: args{
				ctx: context.Background(),
				id:  testPersonID,
			},
			wantQuery: &wantQuery{
				rawQuery:  `DELETE FROM "person"."persons" WHERE "id" = $1 RETURNING *`,
				arguments: []driver.Value{testPersonID},
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

			gotItem, err := tt.ps.Delete(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.wantItem, gotItem)

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_personStore_GetAllPaginated(t *testing.T) {
	type args struct {
		ctx    context.Context
		offset uint
		limit  uint
	}
	type wantQuery struct {
		rawQuery  string
		arguments []driver.Value
		result    *sqlmock.Rows
		returnErr error
	}

	tableRows := []string{"id", "given_name", "family_name", "first_name", "dob", "gender", "pronouns"}
	testSqlBuilder := jagsqlb.NewSqlBuilder()
	mockSession, mockDb, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	testPersons := []models.Person{
		{
			ID:          uuid.New(),
			Name:        models.Name{GivenName: "Testy", FamilyName: "McTesterson", FamilyNameFirst: models.FirstNameGiven},
			DateOfBirth: time.Unix(0, 0),
			Gender:      models.GenderNonBinary,
			Pronouns:    models.Pronouns{Subject: "they", Object: "them"},
		},
		{
			ID:          uuid.New(),
			Name:        models.Name{GivenName: "Tetsuya", FamilyName: "Takahashi", FamilyNameFirst: models.FirstNameFamily},
			DateOfBirth: time.Date(1966, 11, 18, 0, 0, 0, 0, time.UTC),
			Gender:      models.GenderMale,
			Pronouns:    models.Pronouns{Subject: "he", Object: "him"},
		},
		{
			ID:          uuid.New(),
			Name:        models.Name{GivenName: "Brandon", FamilyName: "Williams", FamilyNameFirst: models.FirstNameGiven},
			DateOfBirth: time.Date(1992, 1, 27, 0, 0, 0, 0, time.UTC),
			Gender:      models.GenderMale,
			Pronouns:    models.Pronouns{Subject: "he", Object: "him"},
		},
		{
			ID:          uuid.New(),
			Name:        models.Name{GivenName: "Testita", FamilyName: "Tester", FamilyNameFirst: models.FirstNameGiven},
			DateOfBirth: time.Date(1950, 1, 1, 0, 0, 0, 0, time.UTC),
			Gender:      models.GenderFemale,
			Pronouns:    models.Pronouns{Subject: "she", Object: "her"},
		},
	}

	tests := []struct {
		name      string
		ps        personStore
		args      args
		wantQuery *wantQuery
		wantItems []models.Person
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success; Zero Offset",
			ps: personStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "person.persons",
			},
			args: args{
				ctx:   context.Background(),
				limit: 5,
			},
			wantQuery: &wantQuery{
				rawQuery:  `SELECT * FROM "person"."persons" OFFSET 0 LIMIT 5`,
				arguments: []driver.Value{},
				result: sqlmock.NewRows(tableRows).AddRows(
					[]driver.Value{
						testPersons[0].ID.String(), "Testy", "McTesterson", "given", time.Unix(0, 0), "non-binary", "they/them",
					},
					[]driver.Value{
						testPersons[1].ID.String(), "Tetsuya", "Takahashi", "family", time.Date(1966, 11, 18, 0, 0, 0, 0, time.UTC), "male", "he/him",
					},
					[]driver.Value{
						testPersons[2].ID.String(), "Brandon", "Williams", "given", time.Date(1992, 1, 27, 0, 0, 0, 0, time.UTC), "male", "he/him",
					},
					[]driver.Value{
						testPersons[3].ID.String(), "Testita", "Tester", "given", time.Date(1950, 1, 1, 0, 0, 0, 0, time.UTC), "female", "she/her",
					},
				),
			},
			wantItems: testPersons,
			assertion: assert.NoError,
		},
		{
			name: "Success; Zero Limit",
			ps: personStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "person.persons",
			},
			args: args{
				ctx:    context.Background(),
				offset: 1,
				limit:  0,
			},
			wantQuery: &wantQuery{
				rawQuery:  `SELECT * FROM "person"."persons" OFFSET 1 LIMIT 0;`,
				arguments: []driver.Value{},
				result:    sqlmock.NewRows(tableRows),
			},
			wantItems: []models.Person{},
			assertion: assert.NoError,
		},
		{
			name: "Success with Non-Zero Limit and Offset",
			ps: personStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "person.persons",
			},
			args: args{
				ctx:    context.Background(),
				offset: 1,
				limit:  2,
			},
			wantQuery: &wantQuery{
				rawQuery:  `SELECT * FROM "person"."persons" OFFSET 1 LIMIT 2`,
				arguments: []driver.Value{},
				result: sqlmock.NewRows(tableRows).AddRows(
					[]driver.Value{
						testPersons[1].ID.String(), "Tetsuya", "Takahashi", "family", time.Date(1966, 11, 18, 0, 0, 0, 0, time.UTC), "male", "he/him",
					},
					[]driver.Value{
						testPersons[2].ID.String(), "Brandon", "Williams", "given", time.Date(1992, 1, 27, 0, 0, 0, 0, time.UTC), "male", "he/him",
					},
				),
			},
			wantItems: testPersons[1:3],
			assertion: assert.NoError,
		},
		{
			name: "Error; Query Builder",
			ps: personStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  ".invalid_name",
			},
			args: args{
				ctx:    context.Background(),
				offset: 0,
				limit:  0,
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Query Execution",
			ps: personStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "person.persons",
			},
			args: args{
				ctx:    context.Background(),
				offset: 0,
				limit:  0,
			},
			wantQuery: &wantQuery{
				rawQuery:  `SELECT * FROM "person"."persons" OFFSET 0 LIMIT 0`,
				arguments: []driver.Value{},
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

			gotItems, err := tt.ps.GetAllPaginated(tt.args.ctx, tt.args.offset, tt.args.limit)
			tt.assertion(t, err)
			assert.Equal(t, tt.wantItems, gotItems)

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_personStore_GetSpecific(t *testing.T) {
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

	tableRows := []string{"id", "given_name", "family_name", "first_name", "dob", "gender", "pronouns"}
	testSqlBuilder := jagsqlb.NewSqlBuilder()
	mockSession, mockDb, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	testNotFoundID := uuid.New()
	testPersonID := uuid.New()
	testPerson := models.Person{
		ID: testPersonID,
		Name: models.Name{
			GivenName:       "Testy",
			FamilyName:      "McTesterson",
			FamilyNameFirst: models.FirstNameGiven,
		},
		DateOfBirth: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		Gender:      models.GenderNonBinary,
		Pronouns: models.Pronouns{
			Subject: "they",
			Object:  "them",
		},
	}

	tests := []struct {
		name      string
		ps        personStore
		args      args
		wantQuery *wantQuery
		wantItem  models.Person
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ps: personStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "person.persons",
			},
			args: args{
				ctx: context.Background(),
				id:  testPersonID,
			},
			wantQuery: &wantQuery{
				rawQuery:  `SELECT * FROM "person"."persons" WHERE "id" = $1`,
				arguments: []driver.Value{testPersonID},
				result: sqlmock.NewRows(tableRows).AddRow(
					testPersonID.String(), "Testy", "McTesterson", "given", time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), "non-binary", "they/them",
				),
			},
			wantItem:  testPerson,
			assertion: assert.NoError,
		},
		{
			name: "Error; Query Builder",
			ps: personStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  ".invalid",
			},
			args: args{
				ctx: context.Background(),
				id:  testPersonID,
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Query Execution",
			ps: personStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "person.persons",
			},
			args: args{
				ctx: context.Background(),
				id:  testNotFoundID,
			},
			wantQuery: &wantQuery{
				rawQuery:  `SELECT * FROM "person"."persons" WHERE "id" = $1`,
				arguments: []driver.Value{testNotFoundID},
				result:    sqlmock.NewRows(nil),
				returnErr: sql.ErrNoRows,
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Pronoun Parser",
			ps: personStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "person.persons",
			},
			args: args{
				ctx: context.Background(),
				id:  testPersonID,
			},
			wantQuery: &wantQuery{
				rawQuery:  `SELECT * FROM "person"."persons" WHERE "id" = $1`,
				arguments: []driver.Value{testPersonID},
				result: sqlmock.NewRows(tableRows).AddRow(
					testPersonID.String(), "Testy", "McTesterson", "given", time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), "non-binary", "malformed//data",
				),
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

			gotItem, err := tt.ps.GetSpecific(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.wantItem, gotItem)

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_personStore_Update(t *testing.T) {
	type args struct {
		ctx  context.Context
		id   uuid.UUID
		item models.Person
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

	tests := []struct {
		name      string
		ps        personStore
		args      args
		wantQuery *wantQuery
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ps: personStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "person.persons",
			},
			args: args{
				ctx: context.Background(),
				id:  testPersonID,
				item: models.Person{
					Name: models.Name{
						GivenName:       "Testy",
						FamilyName:      "McTesterson",
						FamilyNameFirst: models.FirstNameGiven,
					},
					DateOfBirth: time.Date(1992, 1, 27, 11, 57, 0, 0, time.FixedZone("EST", -18000)),
					Gender:      models.GenderMale,
					Pronouns: models.Pronouns{
						Subject: "he",
						Object:  "him",
					},
				},
			},
			wantQuery: &wantQuery{
				rawQuery: `UPDATE "person"."persons" SET "given_name"=$1, "family_name"=$2, "first_name"=$3, "dob"=$4, "gender"=$5, "pronouns"=$6 WHERE "id" = $7;`,
				arguments: []driver.Value{
					"Testy",
					"McTesterson",
					models.FirstNameGiven,
					time.Date(1992, 1, 27, 11, 57, 0, 0, time.FixedZone("EST", -18000)),
					models.GenderMale,
					"he/him",
					testPersonID.String(),
				},
				result: sqlmock.NewResult(0, 1),
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; SQL Builder",
			ps: personStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  ".invalid_name",
			},
			args: args{
				ctx: context.Background(),
				id:  testPersonID,
				item: models.Person{
					Name: models.Name{
						GivenName:       "Testy",
						FamilyName:      "McTesterson",
						FamilyNameFirst: models.FirstNameGiven,
					},
					DateOfBirth: time.Date(1992, 1, 27, 11, 57, 0, 0, time.FixedZone("EST", -18000)),
					Gender:      models.GenderMale,
					Pronouns: models.Pronouns{
						Subject: "he",
						Object:  "him",
					},
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; SQL Execution",
			ps: personStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "person.persons",
			},
			args: args{
				ctx: context.Background(),
				id:  testPersonID,
				item: models.Person{
					Name: models.Name{
						GivenName:       "Testy",
						FamilyName:      "McTesterson",
						FamilyNameFirst: models.FirstNameGiven,
					},
					DateOfBirth: time.Date(1992, 1, 27, 11, 57, 0, 0, time.FixedZone("EST", -18000)),
					Gender:      models.GenderMale,
					Pronouns: models.Pronouns{
						Subject: "he",
						Object:  "him",
					},
				},
			},
			wantQuery: &wantQuery{
				rawQuery: `UPDATE "person"."persons" SET "given_name"=$1, "family_name"=$2, "first_name"=$3, "dob"=$4, "gender"=$5, "pronouns"=$6 WHERE "id" = $7;`,
				arguments: []driver.Value{
					"Testy",
					"McTesterson",
					models.FirstNameGiven,
					time.Date(1992, 1, 27, 11, 57, 0, 0, time.FixedZone("EST", -18000)),
					models.GenderMale,
					"he/him",
					testPersonID,
				},
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

			tt.assertion(t, tt.ps.Update(tt.args.ctx, tt.args.id, tt.args.item))

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_personStore_GetSpecificContactAddresses(t *testing.T) {
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

	testPersonStore := personStore{
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
		ps        personStore
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
			ps: personStore{
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

			got, err := tt.ps.GetSpecificContactAddresses(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_personStore_GetSpecificContactEmails(t *testing.T) {
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

	testPersonStore := personStore{
		dbConn:     mockSession,
		sqlBuilder: testSqlBuilder,
		tableNameMap: map[string]string{
			"emails": "person.emails",
		},
	}

	tests := []struct {
		name      string
		ps        personStore
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
			ps: personStore{
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

			got, err := tt.ps.GetSpecificContactEmails(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_personStore_GetSpecificContactPhones(t *testing.T) {
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

	testPersonStore := personStore{
		dbConn:     mockSession,
		sqlBuilder: testSqlBuilder,
		tableNameMap: map[string]string{
			"phones": "person.phones",
		},
	}

	tests := []struct {
		name      string
		ps        personStore
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
			ps: personStore{
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

			got, err := tt.ps.GetSpecificContactPhones(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func TestNewAdminStore(t *testing.T) {
	type args struct {
		dbConn *sql.DB
	}
	tests := []struct {
		name string
		args args
		want PersonStore
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewPersonStore(tt.args.dbConn))
		})
	}
}
