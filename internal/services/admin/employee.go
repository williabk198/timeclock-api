package admin

import (
	"context"

	"github.com/google/uuid"
	"github.com/williabk198/timeclock/internal/datastores"
	"github.com/williabk198/timeclock/internal/models"
)

type EmployeeMicro interface {
	Add(ctx context.Context, employee models.Employee, metadata models.EmployeeMetadata) (uuid.UUID, error)
	Delete(ctx context.Context, id uuid.UUID) (models.Employee, error)
	GetAll(ctx context.Context, offset uint, limit uint) ([]models.Employee, error)
	GetSpecific(ctx context.Context, id uuid.UUID) (models.Employee, error)
	Update(ctx context.Context, id uuid.UUID, newVal models.Employee) error
}

type employeeMicroImpl struct {
	employeeStore     datastores.EmployeeDatastore
	employeeMetaStore datastores.EmployeeMetadataDatastore
}

// Add implements EmployeeMicro.
func (emi employeeMicroImpl) Add(ctx context.Context, employee models.Employee, metadata models.EmployeeMetadata) (uuid.UUID, error) {
	return emi.employeeStore.Add(ctx, employee)
}

// Delete implements EmployeeMicro.
func (emi employeeMicroImpl) Delete(ctx context.Context, id uuid.UUID) (models.Employee, error) {
	return emi.employeeStore.Delete(ctx, id)
}

// GetAll implements EmployeeMicro.
func (emi employeeMicroImpl) GetAll(ctx context.Context, offset uint, limit uint) ([]models.Employee, error) {
	return emi.employeeStore.GetAllPaginated(ctx, offset, limit)
}

// GetSpecific implements EmployeeMicro.
func (emi employeeMicroImpl) GetSpecific(ctx context.Context, id uuid.UUID) (models.Employee, error) {
	return emi.employeeStore.GetSpecific(ctx, id)
}

// Update implements EmployeeMicro.
func (emi employeeMicroImpl) Update(ctx context.Context, id uuid.UUID, newVal models.Employee) error {
	return emi.employeeStore.Update(ctx, id, newVal)
}
