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
	Add(ctx context.Context, employee models.Employee, metadata models.EmployeeMetadata) (id uuid.UUID, err error)
	Delete(ctx context.Context, id uuid.UUID) (item models.Employee, err error)
	GetAllPaginated(ctx context.Context, offset uint, limit uint) (items []models.Employee, err error)
	GetSpecific(ctx context.Context, id uuid.UUID) (item models.Employee, err error)
	UpdateEmployee(ctx context.Context, id uuid.UUID, item models.Employee) (err error)
	UpdateExemptStatus(ctx context.Context, id uuid.UUID, isExempt bool) (err error)
	UpdatePay(ctx context.Context, id uuid.UUID, payData models.EmployeePay) (err error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.EmployeeStatus) (err error)
	UpdateSickTime(ctx context.Context, id uuid.UUID, newVal float64) (err error)
	UpdateTimeOff(ctx context.Context, id uuid.UUID, newVal float64) (err error)
}

type employeeSqlStore struct {
	dbConn            *sql.DB
	sqlBuilder        jagsqlb.SqlBuilder
	tableName         string
	metadataTableName string
}

// Add implements EmployeeDatastore.
func (ess employeeSqlStore) Add(ctx context.Context, employee models.Employee, metadata models.EmployeeMetadata) (id uuid.UUID, err error) {
	//TODO: Have this function also handle inserting EmployeeMetadata so that it can be easily wrapped in a transaction
	query, params, err := ess.sqlBuilder.Insert(ess.tableName).Data(employee).Returning("id").Build()
	if err != nil {
		return uuid.Nil, err
	}

	tx, err := ess.dbConn.BeginTx(ctx, nil)
	if err != nil {
		return uuid.Nil, err
	}

	row := tx.QueryRowContext(ctx, query, params...)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}
	metadata.EmployeeID = id

	query, params, err = ess.sqlBuilder.Insert(ess.metadataTableName).Data(metadata).Build()
	if err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	_, err = tx.ExecContext(ctx, query, params...)
	if err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	tx.Commit()

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

// UpdateEmployee implements EmployeeDatastore.
func (ess employeeSqlStore) UpdateEmployee(ctx context.Context, id uuid.UUID, item models.Employee) (err error) {
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

// UpdateExemptStatus implements [EmployeeDatastore].
func (ess employeeSqlStore) UpdateExemptStatus(ctx context.Context, id uuid.UUID, isExempt bool) (err error) {
	query, params, err := ess.sqlBuilder.Update(ess.metadataTableName).SetMap(map[string]any{"exempt": isExempt}).Where(condition.Equals("eid", id)).Build()
	if err != nil {
		return err
	}

	_, err = ess.dbConn.ExecContext(ctx, query, params...)
	if err != nil {
		return err
	}

	return nil
}

// UpdatePay implements [EmployeeDatastore].
func (ess employeeSqlStore) UpdatePay(ctx context.Context, id uuid.UUID, newPayData models.EmployeePay) (err error) {
	query, params, err := ess.sqlBuilder.Update(ess.metadataTableName).SetMap(map[string]any{"pay": newPayData.String()}).Where(condition.Equals("eid", id)).Build()
	if err != nil {
		return err
	}

	_, err = ess.dbConn.ExecContext(ctx, query, params...)
	if err != nil {
		return err
	}

	return nil
}

// UpdateSickTime implements [EmployeeDatastore].
func (ess employeeSqlStore) UpdateSickTime(ctx context.Context, id uuid.UUID, newVal float64) (err error) {
	query, params, err := ess.sqlBuilder.Update(ess.metadataTableName).SetMap(map[string]any{"sick_time_hrs": newVal}).Where(condition.Equals("eid", id)).Build()
	if err != nil {
		return err
	}

	_, err = ess.dbConn.ExecContext(ctx, query, params...)
	if err != nil {
		return err
	}

	return nil
}

// UpdateStatus implements [EmployeeDatastore].
func (ess employeeSqlStore) UpdateStatus(ctx context.Context, id uuid.UUID, status models.EmployeeStatus) (err error) {
	query, params, err := ess.sqlBuilder.Update(ess.metadataTableName).SetMap(map[string]any{"status": status}).Where(condition.Equals("eid", id)).Build()
	if err != nil {
		return err
	}

	_, err = ess.dbConn.ExecContext(ctx, query, params...)
	if err != nil {
		return err
	}

	return nil
}

// UpdateTimeOff implements [EmployeeDatastore].
func (ess employeeSqlStore) UpdateTimeOff(ctx context.Context, id uuid.UUID, newVal float64) (err error) {
	query, params, err := ess.sqlBuilder.Update(ess.metadataTableName).SetMap(map[string]any{"time_off_hrs": newVal}).Where(condition.Equals("eid", id)).Build()
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
		dbConn:            dbConn,
		sqlBuilder:        jagsqlb.NewSqlBuilder(),
		tableName:         "employees",
		metadataTableName: "metadata.employees",
	}
}
