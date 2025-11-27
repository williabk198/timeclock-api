package admin

import (
	"context"

	"github.com/google/uuid"
	"github.com/williabk198/timeclock/internal/datastores"
	"github.com/williabk198/timeclock/internal/models"
)

type EmployeeMicro interface {
	Add(ctx context.Context, person models.Employee) (uuid.UUID, error)
	Delete(ctx context.Context, id uuid.UUID) (models.Employee, error)
	GetAll(ctx context.Context, offset uint, limit uint) ([]models.Employee, error)
	GetSpecific(ctx context.Context, id uuid.UUID) (models.Employee, error)
	Update(ctx context.Context, id uuid.UUID, newVal models.Employee) error
}

type employeeMicroImpl struct {
	employeeStore datastores.EmployeeDatastore
}

// Add implements EmployeeMicro.
func (e employeeMicroImpl) Add(ctx context.Context, person models.Employee) (uuid.UUID, error) {
	panic("unimplemented")
}

// Delete implements EmployeeMicro.
func (e employeeMicroImpl) Delete(ctx context.Context, id uuid.UUID) (models.Employee, error) {
	panic("unimplemented")
}

// GetAll implements EmployeeMicro.
func (e employeeMicroImpl) GetAll(ctx context.Context, offset uint, limit uint) ([]models.Employee, error) {
	panic("unimplemented")
}

// GetSpecific implements EmployeeMicro.
func (e employeeMicroImpl) GetSpecific(ctx context.Context, id uuid.UUID) (models.Employee, error) {
	panic("unimplemented")
}

// Update implements EmployeeMicro.
func (e employeeMicroImpl) Update(ctx context.Context, id uuid.UUID, newVal models.Employee) error {
	panic("unimplemented")
}
