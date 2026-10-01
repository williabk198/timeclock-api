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
	UpdateEmployee(ctx context.Context, id uuid.UUID, newVal models.Employee) error
	UpdateExemptStatus(ctx context.Context, id uuid.UUID, newVal bool) error
	UpdatePay(ctx context.Context, id uuid.UUID, newVal models.EmployeePay) error
	UpdateSickTime(ctx context.Context, id uuid.UUID, newVal float64) error
	UpdateStatus(ctx context.Context, id uuid.UUID, newVal models.EmployeeStatus) error
	UpdateTimeOff(ctx context.Context, id uuid.UUID, newVal float64) error
}

type employeeMicroImpl struct {
	employeeStore datastores.EmployeeDatastore
}

// Add implements EmployeeMicro.
func (emi employeeMicroImpl) Add(ctx context.Context, employee models.Employee, metadata models.EmployeeMetadata) (uuid.UUID, error) {
	return emi.employeeStore.Add(ctx, employee, metadata)
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

// UpdateEmployee implements EmployeeMicro.
func (emi employeeMicroImpl) UpdateEmployee(ctx context.Context, id uuid.UUID, newVal models.Employee) error {
	return emi.employeeStore.UpdateEmployee(ctx, id, newVal)
}

// UpdateExemptStatus implements [EmployeeMicro].
func (emi employeeMicroImpl) UpdateExemptStatus(ctx context.Context, id uuid.UUID, newVal bool) error {
	return emi.employeeStore.UpdateExemptStatus(ctx, id, newVal)
}

// UpdatePay implements [EmployeeMicro].
func (emi employeeMicroImpl) UpdatePay(ctx context.Context, id uuid.UUID, newVal models.EmployeePay) error {
	return emi.employeeStore.UpdatePay(ctx, id, newVal)
}

// UpdateSickTime implements [EmployeeMicro].
func (emi employeeMicroImpl) UpdateSickTime(ctx context.Context, id uuid.UUID, newVal float64) error {
	return emi.employeeStore.UpdateSickTime(ctx, id, newVal)
}

// UpdateStatus implements [EmployeeMicro].
func (emi employeeMicroImpl) UpdateStatus(ctx context.Context, id uuid.UUID, newVal models.EmployeeStatus) error {
	return emi.employeeStore.UpdateStatus(ctx, id, newVal)
}

// UpdateTimeOff implements [EmployeeMicro].
func (emi employeeMicroImpl) UpdateTimeOff(ctx context.Context, id uuid.UUID, newVal float64) error {
	return emi.employeeStore.UpdateTimeOff(ctx, id, newVal)
}
