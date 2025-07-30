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
				t.Errorf("sql expections were not met: %v", err)
			}
		})
	}
}

func Test_personStore_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name      string
		ps        personStore
		args      args
		wantItem  models.Person
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotItem, err := tt.ps.Delete(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.wantItem, gotItem)
		})
	}
}

func Test_personStore_GetAllPaginated(t *testing.T) {
	type args struct {
		ctx    context.Context
		offset uint
		limit  uint
	}
	tests := []struct {
		name      string
		ps        personStore
		args      args
		wantItems []models.Person
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotItems, err := tt.ps.GetAllPaginated(tt.args.ctx, tt.args.offset, tt.args.limit)
			tt.assertion(t, err)
			assert.Equal(t, tt.wantItems, gotItems)
		})
	}
}

func Test_personStore_GetSpecific(t *testing.T) {
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name      string
		ps        personStore
		args      args
		wantItem  models.Person
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotItem, err := tt.ps.GetSpecific(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.wantItem, gotItem)
		})
	}
}

func Test_personStore_Update(t *testing.T) {
	type args struct {
		ctx  context.Context
		id   uuid.UUID
		item models.Person
	}
	tests := []struct {
		name      string
		ps        personStore
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, tt.ps.Update(tt.args.ctx, tt.args.id, tt.args.item))
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
