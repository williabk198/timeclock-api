package datastores

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/williabk198/jagsqlb"
	"github.com/williabk198/jagsqlb/condition"
	"github.com/williabk198/timeclock/internal/models"
)

type PunchStore interface {
	SqlDatastore[models.Punch, uuid.UUID]
}

type sqlPunchStore struct {
	dbConn     *sql.DB
	sqlBuilder jagsqlb.SqlBuilder
	tableName  string
}

// Add implements [PunchStore].
func (s sqlPunchStore) Add(ctx context.Context, item models.Punch) (id uuid.UUID, err error) {
	query, params, err := s.sqlBuilder.Insert(s.tableName).Data(item).Returning("id").Build()
	if err != nil {
		return uuid.Nil, err
	}

	row := s.dbConn.QueryRowContext(ctx, query, params...)
	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

// Delete implements [PunchStore].
func (s sqlPunchStore) Delete(ctx context.Context, id uuid.UUID) (item models.Punch, err error) {
	query, params, err := s.sqlBuilder.Delete(s.tableName).Where(condition.Equals("id", id)).
		Returning("employee_id", "kind", "reported_time", "timestamp").Build()
	if err != nil {
		return models.Punch{}, err
	}

	row := s.dbConn.QueryRowContext(ctx, query, params...)
	if err := row.Scan(&item.EmployeeID, &item.Kind, &item.ReportedTime, &item.Timestamp); err != nil {
		return models.Punch{}, err
	}

	return item, nil
}

// GetAllPaginated implements [PunchStore].
func (s sqlPunchStore) GetAllPaginated(ctx context.Context, offset uint, limit uint) (items []models.Punch, err error) {
	query, params, err := s.sqlBuilder.Select(s.tableName, "*").Offset(offset).Limit(limit).Build()
	if err != nil {
		return nil, err
	}

	rows, err := s.dbConn.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return s.punchesFromRows(rows)
}

// GetSpecific implements [PunchStore].
func (s sqlPunchStore) GetSpecific(ctx context.Context, id uuid.UUID) (item models.Punch, err error) {
	query, params, err := s.sqlBuilder.Select(s.tableName, "*").Where(condition.Equals("id", id)).Build()
	if err != nil {
		return models.Punch{}, err
	}

	row := s.dbConn.QueryRowContext(ctx, query, params...)
	return s.punchFromRow(row)
}

// Update implements [PunchStore].
func (s sqlPunchStore) Update(ctx context.Context, id uuid.UUID, item models.Punch) (err error) {
	query, params, err := s.sqlBuilder.Update(s.tableName).SetStruct(item).Where(condition.Equals("id", id)).Build()
	if err != nil {
		return err
	}

	if _, err := s.dbConn.ExecContext(ctx, query, params...); err != nil {
		return err
	}

	return nil
}

func (s sqlPunchStore) punchFromRow(row *sql.Row) (models.Punch, error) {
	var item models.Punch
	if err := row.Scan(&item.ID, &item.EmployeeID, &item.Kind, &item.ReportedTime, &item.Timestamp); err != nil {
		return models.Punch{}, err
	}
	return item, nil
}

func (s sqlPunchStore) punchesFromRows(rows *sql.Rows) ([]models.Punch, error) {
	results := make([]models.Punch, 0)
	for rows.Next() {
		var item models.Punch
		if err := rows.Scan(&item.ID, &item.EmployeeID, &item.Kind, &item.ReportedTime, &item.Timestamp); err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func NewSqlPunchStore(dbConn *sql.DB) PunchStore {
	return sqlPunchStore{
		dbConn:     dbConn,
		sqlBuilder: jagsqlb.NewSqlBuilder(),
		tableName:  "punches",
	}
}
