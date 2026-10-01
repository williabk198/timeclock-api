package datastores

import (
	"context"
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

func Test_sqlPunchStore_Add(t *testing.T) {
	type wantQuery struct {
		rawQuery  string
		arguments []driver.Value
		result    *sqlmock.Rows
		returnErr error
	}

	const sqlQuery = `INSERT INTO "punches" ("employee_id", "kind", "reported_time", "timestamp") VALUES ($1, $2, $3, $4) RETURNING "id";`
	testSqlBuilder := jagsqlb.NewSqlBuilder()
	mockSession, mockDb, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	testTableName := "punches"
	testEmployeeID := uuid.New()
	testBadEmployeeID := uuid.New()
	testPunchID := uuid.New()

	tests := []struct {
		name      string // description of this test case
		s         sqlPunchStore
		item      models.Punch
		want      uuid.UUID
		wantQuery *wantQuery
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			s: sqlPunchStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  testTableName,
			},
			item: models.Punch{
				EmployeeID:   testEmployeeID,
				Kind:         models.PunchKindIn,
				ReportedTime: time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC),
				Timestamp:    time.Date(2026, 9, 8, 12, 07, 33, 0, time.UTC),
			},
			want: testPunchID,
			wantQuery: &wantQuery{
				rawQuery: sqlQuery,
				arguments: []driver.Value{
					testEmployeeID.String(),
					int(models.PunchKindIn),
					time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC),
					time.Date(2026, 9, 8, 12, 07, 33, 0, time.UTC),
				},
				result: sqlmock.NewRows([]string{"id"}).AddRow(testPunchID),
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; Query Builder",
			s: sqlPunchStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  ".invalid",
			},
			item:      models.Punch{},
			assertion: assert.Error,
		},
		{
			name: "Error; Query Execution",
			s: sqlPunchStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "punches",
			},
			item: models.Punch{
				EmployeeID:   testBadEmployeeID,
				Kind:         models.PunchKindIn,
				ReportedTime: time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC),
				Timestamp:    time.Date(2026, 9, 8, 12, 07, 33, 0, time.UTC),
			},
			wantQuery: &wantQuery{
				rawQuery: sqlQuery,
				arguments: []driver.Value{
					testBadEmployeeID.String(),
					int(models.PunchKindIn),
					time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC),
					time.Date(2026, 9, 8, 12, 07, 33, 0, time.UTC),
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

			got, gotErr := tt.s.Add(context.Background(), tt.item)
			tt.assertion(t, gotErr)
			assert.Equal(t, tt.want, got)

			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_sqlPunchStore_Delete(t *testing.T) {
	type wantQuery struct {
		rawQuery  string
		arguments []driver.Value
		result    *sqlmock.Rows
		returnErr error
	}

	const sqlQuery = `DELETE FROM "punches" WHERE "id" = $1 RETURNING "employee_id", "kind", "reported_time", "timestamp";`
	tableRows := []string{"employee_id", "kind", "reported_time", "timestamp"}
	testSqlBuilder := jagsqlb.NewSqlBuilder()
	mockSession, mockDb, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	testTableName := "punches"
	testPunchID := uuid.New()
	badPunchID := uuid.New()
	testEmployeeID := uuid.New()

	testPunch := models.Punch{
		EmployeeID:   testEmployeeID,
		Kind:         models.PunchKindIn,
		ReportedTime: time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC),
		Timestamp:    time.Date(2026, 9, 8, 12, 07, 33, 0, time.UTC),
	}

	tests := []struct {
		name      string // description of this test case
		s         sqlPunchStore
		id        uuid.UUID
		want      models.Punch
		wantQuery *wantQuery
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			s: sqlPunchStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  testTableName,
			},
			id:   testPunchID,
			want: testPunch,
			wantQuery: &wantQuery{
				rawQuery:  sqlQuery,
				arguments: []driver.Value{testPunchID.String()},
				result: sqlmock.NewRows(tableRows).AddRow(
					testEmployeeID.String(),
					int(models.PunchKindIn),
					time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC),
					time.Date(2026, 9, 8, 12, 07, 33, 0, time.UTC),
				),
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; Query Builder",
			s: sqlPunchStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  ".invalid",
			},
			id:        testPunchID,
			assertion: assert.Error,
		},
		{
			name: "Error; Query Execution",
			s: sqlPunchStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "punches",
			},
			id: badPunchID,
			wantQuery: &wantQuery{
				rawQuery:  sqlQuery,
				arguments: []driver.Value{badPunchID.String()},
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

			got, gotErr := tt.s.Delete(context.Background(), tt.id)
			tt.assertion(t, gotErr)
			assert.Equal(t, tt.want, got)

			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_sqlPunchStore_GetAllPaginated(t *testing.T) {
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

	tableRows := []string{"id", "employee_id", "kind", "reported_time", "timestamp"}
	testSqlBuilder := jagsqlb.NewSqlBuilder()
	mockSession, mockDb, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	testEmployeeID := uuid.New()
	testPunches := []models.Punch{
		{
			ID:           uuid.New(),
			EmployeeID:   testEmployeeID,
			Kind:         models.PunchKindIn,
			ReportedTime: time.Date(2026, 9, 8, 8, 0, 0, 0, time.UTC),
			Timestamp:    time.Date(2026, 9, 8, 8, 4, 12, 0, time.UTC),
		},
		{
			ID:           uuid.New(),
			EmployeeID:   testEmployeeID,
			Kind:         models.PunchKindMealOut,
			ReportedTime: time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC),
			Timestamp:    time.Date(2026, 9, 8, 12, 1, 33, 0, time.UTC),
		},
		{
			ID:           uuid.New(),
			EmployeeID:   testEmployeeID,
			Kind:         models.PunchKindMealIn,
			ReportedTime: time.Date(2026, 9, 8, 12, 30, 0, 0, time.UTC),
			Timestamp:    time.Date(2026, 9, 8, 12, 28, 47, 0, time.UTC),
		},
		{
			ID:           uuid.New(),
			EmployeeID:   testEmployeeID,
			Kind:         models.PunchKindOut,
			ReportedTime: time.Date(2026, 9, 8, 17, 0, 0, 0, time.UTC),
			Timestamp:    time.Date(2026, 9, 8, 17, 6, 9, 0, time.UTC),
		},
	}

	tests := []struct {
		name      string // description of this test case
		s         sqlPunchStore
		args      args
		wantQuery *wantQuery
		wantItems []models.Punch
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success; Zero Offset",
			s: sqlPunchStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "punches",
			},
			args: args{
				ctx:   context.Background(),
				limit: 5,
			},
			wantQuery: &wantQuery{
				rawQuery:  `SELECT * FROM "punches" OFFSET 0 LIMIT 5`,
				arguments: []driver.Value{},
				result: sqlmock.NewRows(tableRows).AddRows(
					[]driver.Value{
						testPunches[0].ID.String(), testEmployeeID.String(), int(models.PunchKindIn), time.Date(2026, 9, 8, 8, 0, 0, 0, time.UTC), time.Date(2026, 9, 8, 8, 4, 12, 0, time.UTC),
					},
					[]driver.Value{
						testPunches[1].ID.String(), testEmployeeID.String(), int(models.PunchKindMealOut), time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC), time.Date(2026, 9, 8, 12, 1, 33, 0, time.UTC),
					},
					[]driver.Value{
						testPunches[2].ID.String(), testEmployeeID.String(), int(models.PunchKindMealIn), time.Date(2026, 9, 8, 12, 30, 0, 0, time.UTC), time.Date(2026, 9, 8, 12, 28, 47, 0, time.UTC),
					},
					[]driver.Value{
						testPunches[3].ID.String(), testEmployeeID.String(), int(models.PunchKindOut), time.Date(2026, 9, 8, 17, 0, 0, 0, time.UTC), time.Date(2026, 9, 8, 17, 6, 9, 0, time.UTC),
					},
				),
			},
			wantItems: testPunches,
			assertion: assert.NoError,
		},
		{
			name: "Success; Zero Limit",
			s: sqlPunchStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "punches",
			},
			args: args{
				ctx:    context.Background(),
				offset: 1,
				limit:  0,
			},
			wantQuery: &wantQuery{
				rawQuery:  `SELECT * FROM "punches" OFFSET 1 LIMIT 0;`,
				arguments: []driver.Value{},
				result:    sqlmock.NewRows(tableRows),
			},
			wantItems: []models.Punch{},
			assertion: assert.NoError,
		},
		{
			name: "Success with Non-Zero Limit and Offset",
			s: sqlPunchStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "punches",
			},
			args: args{
				ctx:    context.Background(),
				offset: 1,
				limit:  2,
			},
			wantQuery: &wantQuery{
				rawQuery:  `SELECT * FROM "punches" OFFSET 1 LIMIT 2`,
				arguments: []driver.Value{},
				result: sqlmock.NewRows(tableRows).AddRows(
					[]driver.Value{
						testPunches[1].ID.String(), testEmployeeID.String(), int(models.PunchKindMealOut), time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC), time.Date(2026, 9, 8, 12, 1, 33, 0, time.UTC),
					},
					[]driver.Value{
						testPunches[2].ID.String(), testEmployeeID.String(), int(models.PunchKindMealIn), time.Date(2026, 9, 8, 12, 30, 0, 0, time.UTC), time.Date(2026, 9, 8, 12, 28, 47, 0, time.UTC),
					},
				),
			},
			wantItems: testPunches[1:3],
			assertion: assert.NoError,
		},
		{
			name: "Error; Query Builder",
			s: sqlPunchStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  ".invalid",
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
			s: sqlPunchStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "punches",
			},
			args: args{
				ctx:    context.Background(),
				offset: 0,
				limit:  0,
			},
			wantQuery: &wantQuery{
				rawQuery:  `SELECT * FROM "punches" OFFSET 0 LIMIT 0`,
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

			gotItems, err := tt.s.GetAllPaginated(tt.args.ctx, tt.args.offset, tt.args.limit)
			tt.assertion(t, err)
			assert.Equal(t, tt.wantItems, gotItems)

			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_sqlPunchStore_GetSpecific(t *testing.T) {
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

	tableRows := []string{"id", "employee_id", "kind", "reported_time", "timestamp"}
	testSqlBuilder := jagsqlb.NewSqlBuilder()
	mockSession, mockDb, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	testNotFoundID := uuid.New()
	testEmployeeID := uuid.New()
	testPunchID := uuid.New()
	testPunch := models.Punch{
		ID:           testPunchID,
		EmployeeID:   testEmployeeID,
		Kind:         models.PunchKindIn,
		ReportedTime: time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC),
		Timestamp:    time.Date(2026, 9, 8, 12, 07, 33, 0, time.UTC),
	}

	tests := []struct {
		name      string // description of this test case
		s         sqlPunchStore
		args      args
		wantQuery *wantQuery
		wantItem  models.Punch
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			s: sqlPunchStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "punches",
			},
			args: args{
				ctx: context.Background(),
				id:  testPunchID,
			},
			wantQuery: &wantQuery{
				rawQuery:  `SELECT * FROM "punches" WHERE "id" = $1`,
				arguments: []driver.Value{testPunchID.String()},
				result: sqlmock.NewRows(tableRows).AddRow(
					testPunchID.String(), testEmployeeID.String(), int(models.PunchKindIn), time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC), time.Date(2026, 9, 8, 12, 07, 33, 0, time.UTC),
				),
			},
			wantItem:  testPunch,
			assertion: assert.NoError,
		},
		{
			name: "Error; Query Builder",
			s: sqlPunchStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  ".invalid",
			},
			args: args{
				ctx: context.Background(),
				id:  testPunchID,
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Query Execution",
			s: sqlPunchStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "punches",
			},
			args: args{
				ctx: context.Background(),
				id:  testNotFoundID,
			},
			wantQuery: &wantQuery{
				rawQuery:  `SELECT * FROM "punches" WHERE "id" = $1`,
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

			gotItem, err := tt.s.GetSpecific(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.wantItem, gotItem)

			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_sqlPunchStore_Update(t *testing.T) {
	type args struct {
		ctx  context.Context
		id   uuid.UUID
		item models.Punch
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

	testPunchID := uuid.New()
	testEmployeeID := uuid.New()

	tests := []struct {
		name      string // description of this test case
		s         sqlPunchStore
		args      args
		wantQuery *wantQuery
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			s: sqlPunchStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "punches",
			},
			args: args{
				ctx: context.Background(),
				id:  testPunchID,
				item: models.Punch{
					EmployeeID:   testEmployeeID,
					Kind:         models.PunchKindIn,
					ReportedTime: time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC),
					Timestamp:    time.Date(2026, 9, 8, 12, 07, 33, 0, time.UTC),
				},
			},
			wantQuery: &wantQuery{
				rawQuery: `UPDATE "punches" SET "employee_id"=$1, "kind"=$2, "reported_time"=$3, "timestamp"=$4 WHERE "id" = $5;`,
				arguments: []driver.Value{
					testEmployeeID.String(),
					int(models.PunchKindIn),
					time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC),
					time.Date(2026, 9, 8, 12, 07, 33, 0, time.UTC),
					testPunchID.String(),
				},
				result: sqlmock.NewResult(0, 1),
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; SQL Builder",
			s: sqlPunchStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  ".invalid",
			},
			args: args{
				ctx: context.Background(),
				id:  testPunchID,
				item: models.Punch{
					EmployeeID:   testEmployeeID,
					Kind:         models.PunchKindIn,
					ReportedTime: time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC),
					Timestamp:    time.Date(2026, 9, 8, 12, 07, 33, 0, time.UTC),
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; SQL Execution",
			s: sqlPunchStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "punches",
			},
			args: args{
				ctx: context.Background(),
				id:  testPunchID,
				item: models.Punch{
					EmployeeID:   testEmployeeID,
					Kind:         models.PunchKindIn,
					ReportedTime: time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC),
					Timestamp:    time.Date(2026, 9, 8, 12, 07, 33, 0, time.UTC),
				},
			},
			wantQuery: &wantQuery{
				rawQuery: `UPDATE "punches" SET "employee_id"=$1, "kind"=$2, "reported_time"=$3, "timestamp"=$4 WHERE "id" = $5;`,
				arguments: []driver.Value{
					testEmployeeID.String(),
					int(models.PunchKindIn),
					time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC),
					time.Date(2026, 9, 8, 12, 07, 33, 0, time.UTC),
					testPunchID.String(),
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

			tt.assertion(t, tt.s.Update(tt.args.ctx, tt.args.id, tt.args.item))

			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}
