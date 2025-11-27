package datastores

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/williabk198/jagsqlb"
	"github.com/williabk198/jagsqlb/condition"
	"github.com/williabk198/timeclock/internal/models"
)

type EmployeeDatastore interface {
	SqlDatastore[models.Employee, uuid.UUID]
}

type employeeSqlStore struct {
	dbConn     *sql.DB
	sqlBuilder jagsqlb.SqlBuilder
	tableName  string
}

// Add implements EmployeeDatastore.
func (ess employeeSqlStore) Add(ctx context.Context, item models.Employee) (id uuid.UUID, err error) {
	query, params, err := ess.sqlBuilder.Insert(ess.tableName).Data(item).Returning("id").Build()
	if err != nil {
		return uuid.Nil, err
	}

	row := ess.dbConn.QueryRowContext(ctx, query, params...)
	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

// Delete implements EmployeeDatastore.
func (ess employeeSqlStore) Delete(ctx context.Context, id uuid.UUID) (item models.Employee, err error) {
	query, params, err := ess.sqlBuilder.Delete(ess.tableName).Where(condition.Equals("id", id)).Returning("*").Build()
	if err != nil {
		return models.Employee{}, err
	}

	row := ess.dbConn.QueryRowContext(ctx, query, params...)
	return ess.employeeFromRow(row)
}

// GetAllPaginated implements EmployeeDatastore.
func (ess employeeSqlStore) GetAllPaginated(ctx context.Context, offset uint, limit uint) (items []models.Employee, err error) {
	query, params, err := ess.sqlBuilder.Select(ess.tableName, "*").Offset(offset).Limit(limit).Build()
	if err != nil {
		return nil, err
	}

	rows, err := ess.dbConn.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return ess.employeesFromRows(rows)
}

// GetSpecific implements EmployeeDatastore.
func (ess employeeSqlStore) GetSpecific(ctx context.Context, id uuid.UUID) (item models.Employee, err error) {
	query, params, err := ess.sqlBuilder.Select(ess.tableName, "*").Where(condition.Equals("id", id)).Build()
	if err != nil {
		return models.Employee{}, err
	}
	row := ess.dbConn.QueryRowContext(ctx, query, params...)
	return ess.employeeFromRow(row)
}

// Update implements EmployeeDatastore.
func (ess employeeSqlStore) Update(ctx context.Context, id uuid.UUID, item models.Employee) (err error) {
	query, params, err := ess.sqlBuilder.Update(ess.tableName).SetStruct(item).Where(condition.Equals("id", id)).Build()
	if err != nil {
		return err
	}

	_, err = ess.dbConn.ExecContext(ctx, query, params...)
	if err != nil {
		return err
	}

	return nil
}

func (ess employeeSqlStore) employeeFromRow(row *sql.Row) (models.Employee, error) {
	var item models.Employee
	if err := row.Scan(
		&item.ID, &item.PersonID, &item.ReportsToID, &item.Title,
	); err != nil {
		return item, err
	}
	return item, nil
}

func (ess employeeSqlStore) employeesFromRows(rows *sql.Rows) ([]models.Employee, error) {
	result := make([]models.Employee, 0)
	for rows.Next() {
		var item models.Employee
		if err := rows.Scan(
			&item.ID, &item.PersonID, &item.ReportsToID, &item.Title,
		); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func NewEmployeeStore(dbConn *sql.DB) EmployeeDatastore {
	return employeeSqlStore{
		dbConn:     dbConn,
		sqlBuilder: jagsqlb.NewSqlBuilder(),
		tableName:  "employees",
	}
}
