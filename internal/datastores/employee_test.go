package datastores

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/williabk198/jagsqlb"
	"github.com/williabk198/timeclock/internal/models"
)

func Test_employeeSqlStore_Add(t *testing.T) {
	type args struct {
		ctx  context.Context
		item models.Employee
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

	newEmployeeID := uuid.New()
	testPersonID := uuid.New()
	testReportsToID := uuid.New()

	tests := []struct {
		name      string
		e         employeeSqlStore
		args      args
		wantQuery *wantQuery
		wantId    uuid.UUID
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			e: employeeSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "employees",
			},
			args: args{
				ctx: context.Background(),
				item: models.Employee{
					PersonID:    testPersonID,
					ReportsToID: testReportsToID,
					Title:       "QA Tester",
				},
			},
			wantQuery: &wantQuery{
				rawQuery:  `INSERT INTO "employees" ("person_id", "reports_to_eid", "title") VALUES ($1, $2, $3) RETURNING "id";`,
				arguments: []driver.Value{testPersonID.String(), testReportsToID.String(), "QA Tester"},
				result:    sqlmock.NewRows([]string{"id"}).AddRow(newEmployeeID.String()),
			},
			wantId:    newEmployeeID,
			assertion: assert.NoError,
		},
		{
			name: "Error; SQL Builder",
			e: employeeSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  ".invalid",
			},
			args:      args{ctx: context.Background(), item: models.Employee{}},
			assertion: assert.Error,
		},
		{
			name: "Error Query Execution",
			e: employeeSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "employees",
			},
			args: args{
				ctx: context.Background(),
				item: models.Employee{
					PersonID:    newEmployeeID,
					ReportsToID: testReportsToID,
					Title:       "invalid",
				},
			},
			wantQuery: &wantQuery{
				rawQuery:  `INSERT INTO "employees" ("person_id", "reports_to_eid", "title") VALUES ($1, $2, $3) RETURNING "id";`,
				arguments: []driver.Value{newEmployeeID.String(), testReportsToID.String(), "invalid"},
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

			gotId, err := tt.e.Add(tt.args.ctx, tt.args.item)
			tt.assertion(t, err)
			assert.Equal(t, tt.wantId, gotId)

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_employeeSqlStore_Delete(t *testing.T) {
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

	tableRows := []string{"id", "person_id", "reports_to_eid", "title"}
	testSqlBuilder := jagsqlb.NewSqlBuilder()
	mockSession, mockDb, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	removedEmployeeID := uuid.New()
	badEmployeeID := uuid.New()
	testPersonID := uuid.New()

	testEmployee := models.Employee{
		ID:          removedEmployeeID,
		PersonID:    testPersonID,
		ReportsToID: uuid.Nil,
		Title:       "QA Tester",
	}

	tests := []struct {
		name      string
		e         employeeSqlStore
		args      args
		wantQuery *wantQuery
		wantItem  models.Employee
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			e: employeeSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "employees",
			},
			args: args{ctx: context.Background(), id: removedEmployeeID},
			wantQuery: &wantQuery{
				rawQuery:  `DELETE FROM "employees" WHERE "id" = $1 RETURNING *;`,
				arguments: []driver.Value{removedEmployeeID.String()},
				result: sqlmock.NewRows(tableRows).AddRow(
					removedEmployeeID.String(),
					testPersonID.String(),
					uuid.Nil.String(),
					"QA Tester",
				),
			},
			wantItem:  testEmployee,
			assertion: assert.NoError,
		},
		{
			name: "Error; SQL Builder",
			e: employeeSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  ".invalid",
			},
			args:      args{ctx: context.Background(), id: uuid.Nil},
			assertion: assert.Error,
		},
		{
			name: "Error; Query Execution",
			e: employeeSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "employees",
			},
			args: args{ctx: context.Background(), id: badEmployeeID},
			wantQuery: &wantQuery{
				rawQuery:  `DELETE FROM "employees" WHERE "id" = $1 RETURNING *;`,
				arguments: []driver.Value{badEmployeeID.String()},
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

			gotItem, err := tt.e.Delete(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.wantItem, gotItem)

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_employeeSqlStore_GetAllPaginated(t *testing.T) {
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

	tableRows := []string{"id", "person_id", "reports_to_eid", "title"}
	testSqlBuilder := jagsqlb.NewSqlBuilder()
	mockSession, mockDb, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	testOwnerEmployeeID := uuid.New()
	testTechPrezEmployeeID := uuid.New()

	testEmployees := []models.Employee{
		{
			ID:          testOwnerEmployeeID,
			PersonID:    uuid.New(),
			ReportsToID: uuid.Nil,
			Title:       "Owner/President",
		},
		{
			ID:          testTechPrezEmployeeID,
			PersonID:    uuid.New(),
			ReportsToID: testOwnerEmployeeID,
			Title:       "President of Technology",
		},
		{
			ID:          uuid.New(),
			PersonID:    uuid.New(),
			ReportsToID: testTechPrezEmployeeID,
			Title:       "Corporate Systems Manager",
		},
		{
			ID:          uuid.New(),
			PersonID:    uuid.New(),
			ReportsToID: testTechPrezEmployeeID,
			Title:       "IT Manager",
		},
		{
			ID:          uuid.New(),
			PersonID:    uuid.New(),
			ReportsToID: testOwnerEmployeeID,
			Title:       "HR Manager",
		},
	}

	tests := []struct {
		name      string
		e         employeeSqlStore
		args      args
		wantQuery *wantQuery
		wantItems []models.Employee
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success; Zero Offset",
			e: employeeSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "employees",
			},
			args: args{ctx: context.Background(), limit: 5},
			wantQuery: &wantQuery{
				rawQuery:  `SELECT * FROM "employees" OFFSET 0 LIMIT 5;`,
				arguments: []driver.Value{},
				result: sqlmock.NewRows(tableRows).AddRows(
					[]driver.Value{testEmployees[0].ID.String(), testEmployees[0].PersonID.String(), uuid.Nil.String(), "Owner/President"},
					[]driver.Value{testEmployees[1].ID.String(), testEmployees[1].PersonID.String(), testOwnerEmployeeID.String(), "President of Technology"},
					[]driver.Value{testEmployees[2].ID.String(), testEmployees[2].PersonID.String(), testTechPrezEmployeeID.String(), "Corporate Systems Manager"},
					[]driver.Value{testEmployees[3].ID.String(), testEmployees[3].PersonID.String(), testTechPrezEmployeeID.String(), "IT Manager"},
					[]driver.Value{testEmployees[4].ID.String(), testEmployees[4].PersonID.String(), testOwnerEmployeeID.String(), "HR Manager"},
				),
			},
			wantItems: testEmployees,
			assertion: assert.NoError,
		},
		{
			name: "Success; Zero Limit",
			e: employeeSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "employees",
			},
			args: args{ctx: context.Background(), offset: 1, limit: 0},
			wantQuery: &wantQuery{
				rawQuery:  `SELECT * FROM "employees" OFFSET 1 LIMIT 0;`,
				arguments: []driver.Value{},
				result:    sqlmock.NewRows(tableRows),
			},
			wantItems: []models.Employee{},
			assertion: assert.NoError,
		},
		{
			name: "Success; Non-Zero Offset & Limit",
			e: employeeSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "employees",
			},
			args: args{ctx: context.Background(), offset: 1, limit: 2},
			wantQuery: &wantQuery{
				rawQuery:  `SELECT * FROM "employees" OFFSET 1 LIMIT 2;`,
				arguments: []driver.Value{},
				result: sqlmock.NewRows(tableRows).AddRows(
					[]driver.Value{testEmployees[1].ID.String(), testEmployees[1].PersonID.String(), testOwnerEmployeeID.String(), "President of Technology"},
					[]driver.Value{testEmployees[2].ID.String(), testEmployees[2].PersonID.String(), testTechPrezEmployeeID.String(), "Corporate Systems Manager"},
				),
			},
			wantItems: testEmployees[1:3],
			assertion: assert.NoError,
		},
		{
			name: "Error; SQL Builder",
			e: employeeSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  ".invalid",
			},
			args:      args{ctx: context.Background(), offset: 0, limit: 0},
			assertion: assert.Error,
		},
		{
			name: "Error; Query Execution",
			e: employeeSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "employees",
			},
			args: args{ctx: context.Background(), offset: 0, limit: 0},
			wantQuery: &wantQuery{
				rawQuery:  `SELECT * FROM "employees" OFFSET 0 LIMIT 0;`,
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

			gotItems, err := tt.e.GetAllPaginated(tt.args.ctx, tt.args.offset, tt.args.limit)
			tt.assertion(t, err)
			assert.Equal(t, tt.wantItems, gotItems)

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_employeeSqlStore_GetSpecific(t *testing.T) {
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

	tableRows := []string{"id", "person_id", "reports_to_eid", "title"}
	testSqlBuilder := jagsqlb.NewSqlBuilder()
	mockSession, mockDb, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	testPersonID := uuid.New()
	testEmployeeID := uuid.New()
	testEmployeeNotFoundID := uuid.New()
	testEmployee := models.Employee{
		ID:          testEmployeeID,
		PersonID:    testPersonID,
		ReportsToID: uuid.Nil,
		Title:       "Owner/President",
	}

	tests := []struct {
		name      string
		e         employeeSqlStore
		args      args
		wantQuery *wantQuery
		wantItem  models.Employee
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			e: employeeSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "employees",
			},
			args: args{ctx: context.Background(), id: testEmployeeID},
			wantQuery: &wantQuery{
				rawQuery:  `SELECT * FROM "employees" WHERE "id" = $1;`,
				arguments: []driver.Value{testEmployeeID.String()},
				result: sqlmock.NewRows(tableRows).AddRow(
					testEmployeeID.String(), testPersonID.String(), uuid.Nil.String(), "Owner/President",
				),
			},
			wantItem:  testEmployee,
			assertion: assert.NoError,
		},
		{
			name: "Error; SQL Builder",
			e: employeeSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  ".invalid",
			},
			args:      args{ctx: context.Background(), id: uuid.Nil},
			assertion: assert.Error,
		},
		{
			name: "Error; Query Execution",
			e: employeeSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "employees",
			},
			args: args{ctx: context.Background(), id: testEmployeeNotFoundID},
			wantQuery: &wantQuery{
				rawQuery:  `SELECT * FROM "employees" WHERE "id" = $1`,
				arguments: []driver.Value{testEmployeeNotFoundID.String()},
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

			gotItem, err := tt.e.GetSpecific(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.wantItem, gotItem)

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_employeeSqlStore_Update(t *testing.T) {
	type args struct {
		ctx  context.Context
		id   uuid.UUID
		item models.Employee
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

	testEmployeeID := uuid.New()
	testPersonID := uuid.New()

	tests := []struct {
		name      string
		e         employeeSqlStore
		args      args
		wantQuery *wantQuery
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			e: employeeSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "employees",
			},
			args: args{
				ctx: context.Background(),
				id:  testEmployeeID,
				item: models.Employee{
					PersonID:    testPersonID,
					ReportsToID: uuid.Nil,
					Title:       "Owner/President",
				},
			},
			wantQuery: &wantQuery{
				rawQuery:  `UPDATE "employees" SET "reports_to_eid"=$1, "title"=$2 WHERE "id" = $3;`,
				arguments: []driver.Value{uuid.Nil.String(), "Owner/President", testEmployeeID.String()},
				result:    sqlmock.NewResult(0, 1),
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; SQL Builder",
			e: employeeSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  ".invalid",
			},
			args: args{
				ctx: context.Background(),
				id:  testEmployeeID,
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Query Execution",
			e: employeeSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "employees",
			},
			args: args{
				ctx: context.Background(),
				id:  testEmployeeID,
				item: models.Employee{
					ID:          testEmployeeID,
					PersonID:    testPersonID,
					ReportsToID: uuid.Nil,
					Title:       "invalid",
				},
			},
			wantQuery: &wantQuery{
				rawQuery:  `UPDATE "employees" SET "reports_to_eid"=$1, "title"=$2 WHERE "id" = $3;`,
				arguments: []driver.Value{uuid.Nil.String(), "invalid", testEmployeeID.String()},
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

			tt.assertion(t, tt.e.Update(tt.args.ctx, tt.args.id, tt.args.item))

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func TestNewEmployeeStore(t *testing.T) {
	type args struct {
		dbConn *sql.DB
	}

	testSqlBuilder := jagsqlb.NewSqlBuilder()
	mockSession, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		args args
		want EmployeeDatastore
	}{
		{
			name: "Success",
			args: args{
				dbConn: mockSession,
			},
			want: employeeSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "employees",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewEmployeeStore(tt.args.dbConn))
		})
	}
}
