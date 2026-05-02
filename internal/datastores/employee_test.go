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
		ctx      context.Context
		employee models.Employee
		metadata models.EmployeeMetadata
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
		name                 string
		e                    employeeSqlStore
		args                 args
		queryExpectationFunc queryExpectationsFunc
		wantId               uuid.UUID
		assertion            assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			e: employeeSqlStore{
				dbConn:            mockSession,
				sqlBuilder:        testSqlBuilder,
				tableName:         "employees",
				metadataTableName: "metadata.employees",
			},
			args: args{
				ctx: context.Background(),
				employee: models.Employee{
					PersonID:    testPersonID,
					ReportsToID: testReportsToID,
					Title:       "QA Tester",
				},
				metadata: models.EmployeeMetadata{
					Pay:       models.EmployeePay{Currency: "USD", Rate: 25.0, Cadence: models.PayCadenceHourly},
					HireDate:  time.Date(2017, 5, 8, 0, 0, 0, 0, time.UTC),
					StartDate: time.Date(2017, 6, 5, 0, 0, 0, 0, time.UTC),
					SickTime:  40.0,
					TimeOff:   40.0,
					Status:    models.EmployeeStatusActive,
				},
			},
			queryExpectationFunc: func(s sqlmock.Sqlmock) {
				s.ExpectBegin()
				s.ExpectQuery(
					regexp.QuoteMeta(`INSERT INTO "employees" ("person_id", "reports_to_eid", "title") VALUES ($1, $2, $3) RETURNING "id";`)).WithArgs(
					testPersonID.String(), testReportsToID.String(), "QA Tester",
				).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newEmployeeID.String()))

				s.ExpectExec(regexp.QuoteMeta(
					`INSERT INTO "metadata"."employees" ("eid", "pay", "hire_date", "start_date", "sick_time_hrs", "time_off_hrs", "exempt", "status") VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`,
				)).WithArgs(
					newEmployeeID, "25.00 USD/hour", time.Date(2017, 5, 8, 0, 0, 0, 0, time.UTC), time.Date(2017, 6, 5, 0, 0, 0, 0, time.UTC), 40.0, 40.0, false, 1,
				).WillReturnResult(sqlmock.NewResult(0, 1))
				s.ExpectCommit()
			},
			wantId:    newEmployeeID,
			assertion: assert.NoError,
		},
		{
			name: "Error; Employee SQL Builder",
			e: employeeSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  ".invalid",
			},
			args:      args{ctx: context.Background(), employee: models.Employee{}},
			assertion: assert.Error,
		},
		{
			name: "Error; Employee Query Execution",
			e: employeeSqlStore{
				dbConn:     mockSession,
				sqlBuilder: testSqlBuilder,
				tableName:  "employees",
			},
			args: args{
				ctx: context.Background(),
				employee: models.Employee{
					PersonID:    newEmployeeID,
					ReportsToID: testReportsToID,
					Title:       "invalid",
				},
			},
			queryExpectationFunc: func(s sqlmock.Sqlmock) {
				s.ExpectBegin()
				s.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO "employees" ("person_id", "reports_to_eid", "title") VALUES ($1, $2, $3) RETURNING "id";`,
				)).WithArgs(
					newEmployeeID.String(), testReportsToID.String(), "invalid",
				).WillReturnError(assert.AnError)
				s.ExpectRollback()
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Metadata Query Builder",
			e: employeeSqlStore{
				dbConn:            mockSession,
				sqlBuilder:        testSqlBuilder,
				tableName:         "employees",
				metadataTableName: ".invalid",
			},
			args: args{
				ctx: context.Background(),
				employee: models.Employee{
					PersonID:    testPersonID,
					ReportsToID: testReportsToID,
					Title:       "QA Tester",
				},
			},
			queryExpectationFunc: func(s sqlmock.Sqlmock) {
				s.ExpectBegin()
				s.ExpectQuery(
					regexp.QuoteMeta(`INSERT INTO "employees" ("person_id", "reports_to_eid", "title") VALUES ($1, $2, $3) RETURNING "id";`)).WithArgs(
					testPersonID.String(), testReportsToID.String(), "QA Tester",
				).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newEmployeeID.String()))
				s.ExpectRollback()
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Metadata Query Execution",
			e: employeeSqlStore{
				dbConn:            mockSession,
				sqlBuilder:        testSqlBuilder,
				tableName:         "employees",
				metadataTableName: "metadata.employees",
			},
			args: args{
				ctx: context.Background(),
				employee: models.Employee{
					PersonID:    testPersonID,
					ReportsToID: testReportsToID,
					Title:       "QA Tester",
				},
				metadata: models.EmployeeMetadata{
					Pay:       models.EmployeePay{Currency: "USD", Rate: 52000, Cadence: models.PayCadenceYearly},
					HireDate:  time.Date(2017, 5, 8, 0, 0, 0, 0, time.UTC),
					StartDate: time.Date(2017, 6, 5, 0, 0, 0, 0, time.UTC),
					SickTime:  -1.0,
					TimeOff:   -1.0,
					Exempt:    true,
					Status:    models.EmployeeStatusActive,
				},
			},
			queryExpectationFunc: func(s sqlmock.Sqlmock) {
				s.ExpectBegin()
				s.ExpectQuery(
					regexp.QuoteMeta(`INSERT INTO "employees" ("person_id", "reports_to_eid", "title") VALUES ($1, $2, $3) RETURNING "id";`)).WithArgs(
					testPersonID.String(), testReportsToID.String(), "QA Tester",
				).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newEmployeeID.String()))

				s.ExpectExec(regexp.QuoteMeta(
					`INSERT INTO "metadata"."employees" ("eid", "pay", "hire_date", "start_date", "sick_time_hrs", "time_off_hrs", "exempt", "status") VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`,
				)).WithArgs(
					newEmployeeID, "52000.00 USD/year", time.Date(2017, 5, 8, 0, 0, 0, 0, time.UTC), time.Date(2017, 6, 5, 0, 0, 0, 0, time.UTC), -1.0, -1.0, true, 1,
				).WillReturnError(assert.AnError)
				s.ExpectRollback()
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.queryExpectationFunc != nil {
				tt.queryExpectationFunc.SetExpectactions(mockDb)
			}
			gotId, err := tt.e.Add(tt.args.ctx, tt.args.employee, tt.args.metadata)
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

			tt.assertion(t, tt.e.UpdateEmployee(tt.args.ctx, tt.args.id, tt.args.item))

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
				dbConn:            mockSession,
				sqlBuilder:        testSqlBuilder,
				tableName:         "employees",
				metadataTableName: "metadata.employees",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewEmployeeStore(tt.args.dbConn))
		})
	}
}

func Test_employeeSqlStore_UpdateExemptStatus(t *testing.T) {
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
		ess       employeeSqlStore
		args      args
		wantQuery *wantQuery
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ess: employeeSqlStore{
				dbConn:            mockSession,
				sqlBuilder:        testSqlBuilder,
				metadataTableName: "metadata.employees",
			},
			args: args{
				ctx:          context.Background(),
				employeeID:   testEmployeeID,
				newExemptVal: false,
			},
			wantQuery: &wantQuery{
				rawQuery:  `UPDATE "metadata"."employees" SET "exempt"=$1 WHERE "eid" = $2;`,
				arguments: []driver.Value{false, testEmployeeID.String()},
				result:    sqlmock.NewResult(0, 1),
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; SQL Builder",
			ess: employeeSqlStore{
				dbConn:            mockSession,
				sqlBuilder:        testSqlBuilder,
				metadataTableName: ".invalid",
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
			ess: employeeSqlStore{
				dbConn:            mockSession,
				sqlBuilder:        testSqlBuilder,
				metadataTableName: "metadata.employees",
			},
			args: args{
				ctx:          context.Background(),
				employeeID:   testNotFoundID,
				newExemptVal: false,
			},
			wantQuery: &wantQuery{
				rawQuery:  `UPDATE "metadata"."employees" SET "exempt"=$1 WHERE "eid" = $2;`,
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

			tt.assertion(t, tt.ess.UpdateExemptStatus(tt.args.ctx, tt.args.employeeID, tt.args.newExemptVal))

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_employeeSqlStore_UpdatePay(t *testing.T) {
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
		ess       employeeSqlStore
		args      args
		wantQuery *wantQuery
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ess: employeeSqlStore{
				dbConn:            mockSession,
				sqlBuilder:        testSqlBuilder,
				metadataTableName: "metadata.employees",
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
				rawQuery:  `UPDATE "metadata"."employees" SET "pay"=$1 WHERE "eid" = $2;`,
				arguments: []driver.Value{"25.00 USD/hour", testEmployeeID.String()},
				result:    sqlmock.NewResult(0, 1),
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; SQL Builder",
			ess: employeeSqlStore{
				dbConn:            mockSession,
				sqlBuilder:        testSqlBuilder,
				metadataTableName: ".invalid",
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
			ess: employeeSqlStore{
				dbConn:            mockSession,
				sqlBuilder:        testSqlBuilder,
				metadataTableName: "metadata.employees",
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
				rawQuery:  `UPDATE "metadata"."employees" SET "pay"=$1 WHERE "eid" = $2;`,
				arguments: []driver.Value{"25.00 USD/hour", testNotFoundID.String()},
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

			tt.assertion(t, tt.ess.UpdatePay(tt.args.ctx, tt.args.employeeID, tt.args.newPayInfo))

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_employeeSqlStore_UpdateStatus(t *testing.T) {
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
		ess       employeeSqlStore
		args      args
		wantQuery *wantQuery
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ess: employeeSqlStore{
				dbConn:            mockSession,
				sqlBuilder:        testSqlBuilder,
				metadataTableName: "metadata.employees",
			},
			args: args{
				ctx:        context.Background(),
				employeeID: testEmployeeID,
				newStatus:  models.EmployeeStatusInactive,
			},
			wantQuery: &wantQuery{
				rawQuery:  `UPDATE "metadata"."employees" SET "status"=$1 WHERE "eid" = $2;`,
				arguments: []driver.Value{2, testEmployeeID.String()},
				result:    sqlmock.NewResult(0, 1),
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; SQL Builder",
			ess: employeeSqlStore{
				dbConn:            mockSession,
				sqlBuilder:        testSqlBuilder,
				metadataTableName: ".invalid",
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
			ess: employeeSqlStore{
				dbConn:            mockSession,
				sqlBuilder:        testSqlBuilder,
				metadataTableName: "metadata.employees",
			},
			args: args{
				ctx:        context.Background(),
				employeeID: testNotFoundID,
				newStatus:  models.EmployeeStatusInactive,
			},
			wantQuery: &wantQuery{
				rawQuery:  `UPDATE "metadata"."employees" SET "status"=$1 WHERE "eid" = $2;`,
				arguments: []driver.Value{2, testNotFoundID.String()},
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

			tt.assertion(t, tt.ess.UpdateStatus(tt.args.ctx, tt.args.employeeID, tt.args.newStatus))

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_employeeSqlStore_UpdateSickTime(t *testing.T) {
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
		ess       employeeSqlStore
		args      args
		wantQuery *wantQuery
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ess: employeeSqlStore{
				dbConn:            mockSession,
				sqlBuilder:        testSqlBuilder,
				metadataTableName: "metadata.employees",
			},
			args: args{
				ctx:        context.Background(),
				employeeID: testEmployeeID,
				newVal:     17.5,
			},
			wantQuery: &wantQuery{
				rawQuery:  `UPDATE "metadata"."employees" SET "sick_time_hrs"=$1 WHERE "eid" = $2;`,
				arguments: []driver.Value{17.5, testEmployeeID.String()},
				result:    sqlmock.NewResult(0, 1),
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; SQL Builder",
			ess: employeeSqlStore{
				dbConn:            mockSession,
				sqlBuilder:        testSqlBuilder,
				metadataTableName: ".invalid",
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
			ess: employeeSqlStore{
				dbConn:            mockSession,
				sqlBuilder:        testSqlBuilder,
				metadataTableName: "metadata.employees",
			},
			args: args{
				ctx:        context.Background(),
				employeeID: testBadEmployeeID,
				newVal:     0.0,
			},
			wantQuery: &wantQuery{
				rawQuery:  `UPDATE "metadata"."employees" SET "sick_time_hrs"=$1 WHERE "eid" = $2;`,
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

			err := tt.ess.UpdateSickTime(tt.args.ctx, tt.args.employeeID, tt.args.newVal)
			tt.assertion(t, err)

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}

func Test_employeeSqlStore_UpdateTimeOff(t *testing.T) {
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
		ess       employeeSqlStore
		args      args
		wantQuery *wantQuery
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ess: employeeSqlStore{
				dbConn:            mockSession,
				sqlBuilder:        testSqlBuilder,
				metadataTableName: "metadata.employees",
			},
			args: args{
				ctx:        context.Background(),
				employeeID: testEmployeeID,
				newVal:     36.5,
			},
			wantQuery: &wantQuery{
				rawQuery:  `UPDATE "metadata"."employees" SET "time_off_hrs"=$1 WHERE "eid" = $2;`,
				arguments: []driver.Value{36.5, testEmployeeID.String()},
				result:    sqlmock.NewResult(0, 1),
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; SQL Builder",
			ess: employeeSqlStore{
				dbConn:            mockSession,
				sqlBuilder:        testSqlBuilder,
				metadataTableName: ".invalid",
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
			ess: employeeSqlStore{
				dbConn:            mockSession,
				sqlBuilder:        testSqlBuilder,
				metadataTableName: "metadata.employees",
			},
			args: args{
				ctx:        context.Background(),
				employeeID: testNotFoundID,
				newVal:     36.5,
			},
			wantQuery: &wantQuery{
				rawQuery:  `UPDATE "metadata"."employees" SET "time_off_hrs"=$1 WHERE "eid" = $2;`,
				arguments: []driver.Value{36.5, testNotFoundID.String()},
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

			err := tt.ess.UpdateTimeOff(tt.args.ctx, tt.args.employeeID, tt.args.newVal)
			tt.assertion(t, err)

			// Check to see if the expectations of the query were met
			if err := mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("sql expectations were not met: %v", err)
			}
		})
	}
}
