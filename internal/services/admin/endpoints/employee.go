package endpoints

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/williabk198/timeclock/internal/models"
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
func (a adminEmployeeEndpoints) Add(ctx context.Context, employee EmployeeData) (EmployeeData, error) {
	personID, err := uuid.Parse(employee.PersonID)
	if err != nil {
		return EmployeeData{}, fmt.Errorf("failed to parse person ID: %w", err)
	}

	reportsToID, err := uuid.Parse(employee.ReportsToID)
	if err != nil {
		return EmployeeData{}, fmt.Errorf("failed to parse reportsTo ID: %w", err)
	}

	id, err := a.employeeMicro.Add(ctx, models.Employee{
		PersonID:    personID,
		ReportsToID: reportsToID,
		Title:       employee.Title,
	})
	if err != nil {
		return EmployeeData{}, fmt.Errorf("failed to add employee to DB: %w", err)
	}

	employee.ID = id.String()
	return employee, nil
}

// Delete implements EmployeeEndpoints.
func (a adminEmployeeEndpoints) Delete(ctx context.Context, id string) (EmployeeData, error) {
	employeeID, err := uuid.Parse(id)
	if err != nil {
		return EmployeeData{}, fmt.Errorf("failed to parse employee ID: %w", err)
	}

	employee, err := a.employeeMicro.Delete(ctx, employeeID)
	if err != nil {
		return EmployeeData{}, fmt.Errorf("failed to remove employee: %w", err)
	}

	return EmployeeData{
		PersonID:    employee.PersonID.String(),
		ReportsToID: employee.ReportsToID.String(),
		Title:       employee.Title,
	}, nil
}

// GetAll implements EmployeeEndpoints.
func (a adminEmployeeEndpoints) GetAll(ctx context.Context, reqData GetPaginatedRequestData) ([]EmployeeData, error) {
	employees, err := a.employeeMicro.GetAll(ctx, reqData.Offset, reqData.Limit)
	if err != nil {
		return []EmployeeData{}, fmt.Errorf("failed to fetch employee records: %w", err)
	}

	results := make([]EmployeeData, len(employees))
	for i, e := range employees {
		results[i] = EmployeeData{
			ID:          e.ID.String(),
			PersonID:    e.PersonID.String(),
			ReportsToID: e.ReportsToID.String(),
			Title:       e.Title,
		}
	}

	return results, nil
}

// GetSpecific implements EmployeeEndpoints.
func (a adminEmployeeEndpoints) GetSpecific(ctx context.Context, id string) (EmployeeData, error) {
	employeeID, err := uuid.Parse(id)
	if err != nil {
		return EmployeeData{}, fmt.Errorf("failed to parse employee ID: %w", err)
	}

	employee, err := a.employeeMicro.GetSpecific(ctx, employeeID)
	if err != nil {
		return EmployeeData{}, fmt.Errorf("failed to fetch person from DB: %w", err)
	}

	return EmployeeData{
		ID:          employee.ID.String(),
		PersonID:    employee.PersonID.String(),
		ReportsToID: employee.ReportsToID.String(),
		Title:       employee.Title,
	}, nil
}

// Update implements EmployeeEndpoints.
func (a adminEmployeeEndpoints) Update(ctx context.Context, urd UpdateRequestData[EmployeeData]) (EmployeeData, error) {
	employeeID, err := uuid.Parse(urd.ID)
	if err != nil {
		return EmployeeData{}, fmt.Errorf("failed to parse employee ID: %w", err)
	}

	reportsToID, err := uuid.Parse(urd.Data.ReportsToID)
	if err != nil {
		return EmployeeData{}, fmt.Errorf("failed to parse reportsTo ID: %w", err)
	}

	err = a.employeeMicro.Update(ctx, employeeID, models.Employee{
		ReportsToID: reportsToID,
		Title:       urd.Data.Title,
	})
	if err != nil {
		return EmployeeData{}, fmt.Errorf("failed to update employee: %w", err)
	}

	return urd.Data, nil
}
