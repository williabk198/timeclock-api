package datastores

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/williabk198/jagsqlb"
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
func (e employeeSqlStore) Add(ctx context.Context, item models.Employee) (id uuid.UUID, err error) {
	panic("unimplemented")
}

// Delete implements EmployeeDatastore.
func (e employeeSqlStore) Delete(ctx context.Context, id uuid.UUID) (item models.Employee, err error) {
	panic("unimplemented")
}

// GetAllPaginated implements EmployeeDatastore.
func (e employeeSqlStore) GetAllPaginated(ctx context.Context, offset uint, limit uint) (items []models.Employee, err error) {
	panic("unimplemented")
}

// GetSpecific implements EmployeeDatastore.
func (e employeeSqlStore) GetSpecific(ctx context.Context, id uuid.UUID) (item models.Employee, err error) {
	panic("unimplemented")
}

// Update implements EmployeeDatastore.
func (e employeeSqlStore) Update(ctx context.Context, id uuid.UUID, item models.Employee) (err error) {
	panic("unimplemented")
}

func NewEmployeeStore(dbConn *sql.DB) EmployeeDatastore {
	panic("unimplemented")
}
