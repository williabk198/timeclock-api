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

func Test_employeeMetadataSqlStore_Add(t *testing.T) {
	type args struct {
		ctx  context.Context
		data models.EmployeeMetadata
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
	testBadEmployeeID := uuid.New()
	testHireDate := time.Date(2017, 5, 4, 0, 0, 0, 0, time.UTC)
	testStartDate := time.Date(2017, 6, 5, 13, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		emss      employeeMetadataSqlStore
		args      args
		wantQuery *wantQuery
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			emss: employeeMetadataSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "metadata.employees",
			},
			args: args{
				ctx: context.Background(),
				data: models.EmployeeMetadata{
					EmployeeID: testEmployeeID,
					Pay: models.EmployeePay{
						Currency: "USD",
						Rate:     19.0,
						Cadence:  models.PayCadenceHourly,
					},
					HireDate:  testHireDate,
					StartDate: testStartDate,
					SickTime:  16.0,
					TimeOff:   20.0,
					Exempt:    false,
					Status:    models.EmployeeStatusActive,
				},
			},
			wantQuery: &wantQuery{
				rawQuery:  `INSERT INTO "metadata"."employees" ("eid", "pay", "hire_date", "start_date", "sick_time", "time_off", "exempt", "status") VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`,
				arguments: []driver.Value{testEmployeeID.String(), "19.00 USD/hour", testHireDate, testStartDate, 16.0, 20.0, false, 1},
				result:    driver.ResultNoRows,
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; SQL Builder",
			emss: employeeMetadataSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  ".invalid",
			},
			args: args{
				ctx:  context.Background(),
				data: models.EmployeeMetadata{},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Query Execution",
			emss: employeeMetadataSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "metadata.employees",
			},
			args: args{
				ctx: context.Background(),
				data: models.EmployeeMetadata{
					EmployeeID: testBadEmployeeID,
					Pay: models.EmployeePay{
						Currency: "CAD",
						Rate:     100_000.0,
						Cadence:  models.PayCadenceYearly,
					},
					HireDate:  testHireDate,
					StartDate: testStartDate,
					SickTime:  16.0,
					TimeOff:   20.0,
					Exempt:    true,
					Status:    models.EmployeeStatusActive,
				},
			},
			wantQuery: &wantQuery{
				rawQuery:  `INSERT INTO "metadata"."employees" ("eid", "pay", "hire_date", "start_date", "sick_time", "time_off", "exempt", "status") VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`,
				arguments: []driver.Value{testBadEmployeeID.String(), "100000.00 CAD/year", testHireDate, testStartDate, 16.0, 20.0, true, 1},
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

			tt.assertion(t, tt.emss.Add(tt.args.ctx, tt.args.data))

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_employeeMetadataSqlStore_Get(t *testing.T) {
	type args struct {
		ctx        context.Context
		employeeID uuid.UUID
	}
	type wantQuery struct {
		rawQuery  string
		arguments []driver.Value
		result    *sqlmock.Rows
		returnErr error
	}

	columns := []string{"eid", "pay", "hire_date", "start_date", "sick_time", "time_off", "exempt", "status"}
	testSqlBuilder := jagsqlb.NewSqlBuilder()
	mockSession, mockDb, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	testEmployeeID := uuid.New()
	testNotFoundID := uuid.New()
	testHireDate := time.Date(2017, 5, 4, 0, 0, 0, 0, time.UTC)
	testStartDate := time.Date(2017, 6, 5, 13, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		emss      employeeMetadataSqlStore
		args      args
		wantQuery *wantQuery
		want      models.EmployeeMetadata
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			emss: employeeMetadataSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "metadata.employees",
			},
			args: args{
				ctx:        context.Background(),
				employeeID: testEmployeeID,
			},
			wantQuery: &wantQuery{
				rawQuery:  `SELECT * FROM "metadata"."employees" WHERE "id" = $1;`,
				arguments: []driver.Value{testEmployeeID.String()},
				result: sqlmock.NewRows(columns).AddRow(
					testEmployeeID.String(), "74000 USD/year", testHireDate, testStartDate, 20.0, 40.0, true, 1,
				),
			},
			want:      models.EmployeeMetadata{},
			assertion: assert.NoError,
		},
		{
			name: "Error; SQL Builder",
			emss: employeeMetadataSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  ".invalid",
			},
			args: args{
				ctx:        context.Background(),
				employeeID: testEmployeeID,
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Query Execution",
			emss: employeeMetadataSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "metadata.employees",
			},
			args: args{
				ctx:        context.Background(),
				employeeID: testEmployeeID,
			},
			wantQuery: &wantQuery{
				rawQuery:  `SELECT * FROM "metadata"."employees" WHERE "id" = $1;`,
				arguments: []driver.Value{testNotFoundID.String()},
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

			got, err := tt.emss.Get(tt.args.ctx, tt.args.employeeID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_employeeMetadataSqlStore_UpdateExemptStatus(t *testing.T) {
	type args struct {
		ctx          context.Context
		employeeID   uuid.UUID
		newExemptVal bool
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
	testNotFoundID := uuid.New()

	tests := []struct {
		name      string
		emss      employeeMetadataSqlStore
		args      args
		wantQuery *wantQuery
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			emss: employeeMetadataSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "metadata.employees",
			},
			args: args{
				ctx:          context.Background(),
				employeeID:   testEmployeeID,
				newExemptVal: false,
			},
			wantQuery: &wantQuery{
				rawQuery:  `UPDATE "metadata"."employees" SET "exempt" = $1 WHERE "id" = $2;`,
				arguments: []driver.Value{false, testEmployeeID.String()},
				result:    sqlmock.NewResult(0, 1),
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; SQL Builder",
			emss: employeeMetadataSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  ".invalid",
			},
			args: args{
				ctx:          context.Background(),
				employeeID:   testEmployeeID,
				newExemptVal: false,
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Query Execution",
			emss: employeeMetadataSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "metadata.employees",
			},
			args: args{
				ctx:          context.Background(),
				employeeID:   testNotFoundID,
				newExemptVal: false,
			},
			wantQuery: &wantQuery{
				rawQuery:  `UPDATE "metadata"."employees" SET "exempt" = $1 WHERE "id" = $2;`,
				arguments: []driver.Value{false, testNotFoundID.String()},
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

			tt.assertion(t, tt.emss.UpdateExemptStatus(tt.args.ctx, tt.args.employeeID, tt.args.newExemptVal))

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_employeeMetadataSqlStore_UpdatePay(t *testing.T) {
	type args struct {
		ctx        context.Context
		employeeID uuid.UUID
		newPayInfo models.EmployeePay
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
	testNotFoundID := uuid.New()

	tests := []struct {
		name      string
		emss      employeeMetadataSqlStore
		args      args
		wantQuery *wantQuery
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			emss: employeeMetadataSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "metadata.employees",
			},
			args: args{
				ctx:        context.Background(),
				employeeID: testEmployeeID,
				newPayInfo: models.EmployeePay{
					Currency: "USD",
					Rate:     25.0,
					Cadence:  models.PayCadenceHourly,
				},
			},
			wantQuery: &wantQuery{
				rawQuery:  `UPDATE "metadata"."employees" SET pay = $1 WHERE "id" = $2;`,
				arguments: []driver.Value{"25 USD/hour", testEmployeeID.String()},
				result:    sqlmock.NewResult(0, 1),
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; SQL Builder",
			emss: employeeMetadataSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  ".invalid",
			},
			args: args{
				ctx:        context.Background(),
				employeeID: testEmployeeID,
				newPayInfo: models.EmployeePay{},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Query Execution",
			emss: employeeMetadataSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "metadata.employees",
			},
			args: args{
				ctx:        context.Background(),
				employeeID: testNotFoundID,
				newPayInfo: models.EmployeePay{
					Currency: "USD",
					Rate:     25.0,
					Cadence:  models.PayCadenceHourly,
				},
			},
			wantQuery: &wantQuery{
				rawQuery:  `UPDATE "metadata"."employees" SET pay = $1 WHERE "id" = $2;`,
				arguments: []driver.Value{"25 USD/hour", testNotFoundID.String()},
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

			tt.assertion(t, tt.emss.UpdatePay(tt.args.ctx, tt.args.employeeID, tt.args.newPayInfo))

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_employeeMetadataSqlStore_UpdateStatus(t *testing.T) {
	type args struct {
		ctx        context.Context
		employeeID uuid.UUID
		newStatus  models.EmployeeStatus
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
	testNotFoundID := uuid.New()

	tests := []struct {
		name      string
		emss      employeeMetadataSqlStore
		args      args
		wantQuery *wantQuery
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			emss: employeeMetadataSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "metadata.employees",
			},
			args: args{
				ctx:        context.Background(),
				employeeID: testEmployeeID,
				newStatus:  models.EmployeeStatusInactive,
			},
			wantQuery: &wantQuery{
				rawQuery:  `UPDATE "metadata"."employees" SET "status" = $1 WHERE "id" = $2;`,
				arguments: []driver.Value{3, testEmployeeID.String()},
				result:    sqlmock.NewResult(0, 1),
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; SQL Builder",
			emss: employeeMetadataSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  ".invalid",
			},
			args: args{
				ctx:        context.Background(),
				employeeID: testEmployeeID,
				newStatus:  models.EmployeeStatusGone,
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Query Execution",
			emss: employeeMetadataSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "metadata.employees",
			},
			args: args{
				ctx:        context.Background(),
				employeeID: testNotFoundID,
				newStatus:  models.EmployeeStatusInactive,
			},
			wantQuery: &wantQuery{
				rawQuery:  `UPDATE "metadata"."employees" SET "status" = $1 WHERE "id" = $2;`,
				arguments: []driver.Value{3, testNotFoundID.String()},
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

			tt.assertion(t, tt.emss.UpdateStatus(tt.args.ctx, tt.args.employeeID, tt.args.newStatus))

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_employeeMetadataSqlStore_UpdateSickTime(t *testing.T) {
	type args struct {
		ctx        context.Context
		employeeID uuid.UUID
		newVal     float64
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
	testBadEmployeeID := uuid.New()

	tests := []struct {
		name      string
		emss      employeeMetadataSqlStore
		args      args
		wantQuery *wantQuery
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			emss: employeeMetadataSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  ".invalid",
			},
			args: args{
				ctx:        context.Background(),
				employeeID: testEmployeeID,
				newVal:     17.5,
			},
			wantQuery: &wantQuery{
				rawQuery:  `UPDATE "metadata"."employees" SET "sick_time" = $1 WHERE "eid" = $2;`,
				arguments: []driver.Value{17.5, testEmployeeID.String()},
				result:    sqlmock.NewResult(0, 1),
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; SQL Builder",
			emss: employeeMetadataSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  ".invalid",
			},
			args: args{
				ctx:        context.Background(),
				employeeID: testEmployeeID,
				newVal:     10,
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Query Execution",
			emss: employeeMetadataSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "metadata.employees",
			},
			args: args{
				ctx:        context.Background(),
				employeeID: testBadEmployeeID,
				newVal:     0,
			},
			wantQuery: &wantQuery{
				rawQuery:  `UPDATE "metadata"."employees" SET "sick_time" = $1 WHERE "eid" = $2;`,
				arguments: []driver.Value{0.0, testBadEmployeeID.String()},
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

			err := tt.emss.AdjustSickTime(tt.args.ctx, tt.args.employeeID, tt.args.newVal)
			tt.assertion(t, err)

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_employeeMetadataSqlStore_UpdateTimeOff(t *testing.T) {
	type args struct {
		ctx        context.Context
		employeeID uuid.UUID
		newVal     float64
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
	testNotFoundID := uuid.New()

	tests := []struct {
		name      string
		emss      employeeMetadataSqlStore
		args      args
		wantQuery *wantQuery
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			emss: employeeMetadataSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "metadata.employees",
			},
			args: args{
				ctx:        context.Background(),
				employeeID: testEmployeeID,
				newVal:     36.5,
			},
			wantQuery: &wantQuery{
				rawQuery:  `UPDATE "metadata"."employees" SET "time_off" = $1 WHERE "id" = $2;`,
				arguments: []driver.Value{36.5, testEmployeeID.String()},
				result:    sqlmock.NewResult(0, 1),
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; SQL Builder",
			emss: employeeMetadataSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  ".invalid",
			},
			args: args{
				ctx:        context.Background(),
				employeeID: testEmployeeID,
				newVal:     36.5,
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Query Execution",
			emss: employeeMetadataSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "metadata.employees",
			},
			args: args{
				ctx:        context.Background(),
				employeeID: testNotFoundID,
				newVal:     36.5,
			},
			wantQuery: &wantQuery{
				rawQuery:  `UPDATE "metadata"."employees" SET "time_off" = $1 WHERE "id" = $2;`,
				arguments: []driver.Value{36.5, testNotFoundID.String()},
				returnErr: assert.AnError,
			},
			assertion: assert.NoError,
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

			err := tt.emss.AdjustTimeOff(tt.args.ctx, tt.args.employeeID, tt.args.newVal)
			tt.assertion(t, err)

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}
