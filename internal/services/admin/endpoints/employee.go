package endpoints

import (
	"context"

	"github.com/williabk198/timeclock/internal/services/admin"
)

type EmployeeEndpoints interface {
	Add(ctx context.Context, person EmployeeData) (EmployeeData, error)
	Delete(ctx context.Context, id string) (EmployeeData, error)
	GetSpecific(ctx context.Context, id string) (EmployeeData, error)
	GetAll(ctx context.Context, reqData GetPaginatedRequestData) ([]EmployeeData, error)
	Update(ctx context.Context, urd UpdateRequestData[EmployeeData]) (EmployeeData, error)
}

type adminEmployeeEndpoints struct {
	employeeMicro admin.EmployeeMicro
}

// Add implements EmployeeEndpoints.
func (a adminEmployeeEndpoints) Add(ctx context.Context, person EmployeeData) (EmployeeData, error) {
	panic("unimplemented")
}

// Delete implements EmployeeEndpoints.
func (a adminEmployeeEndpoints) Delete(ctx context.Context, id string) (EmployeeData, error) {
	panic("unimplemented")
}

// GetAll implements EmployeeEndpoints.
func (a adminEmployeeEndpoints) GetAll(ctx context.Context, reqData GetPaginatedRequestData) ([]EmployeeData, error) {
	panic("unimplemented")
}

// GetSpecific implements EmployeeEndpoints.
func (a adminEmployeeEndpoints) GetSpecific(ctx context.Context, id string) (EmployeeData, error) {
	panic("unimplemented")
}

// Update implements EmployeeEndpoints.
func (a adminEmployeeEndpoints) Update(ctx context.Context, urd UpdateRequestData[EmployeeData]) (EmployeeData, error) {
	panic("unimplemented")
}
