package endpoints

import (
	"context"
	"fmt"
	"time"

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
	UpdateExemptStatus(ctx context.Context, urd UpdateRequestData[bool]) (bool, error)
	UpdatePay(ctx context.Context, urd UpdateRequestData[models.EmployeePay]) (models.EmployeePay, error)
	UpdateSickTimeHours(ctx context.Context, urd UpdateRequestData[float64]) (float64, error)
	UpdateStatus(ctx context.Context, urd UpdateRequestData[int]) (int, error)
	UpdateTimeOffHours(ctx context.Context, urd UpdateRequestData[float64]) (float64, error)
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
	}, models.EmployeeMetadata{
		Pay:       employee.Metadata.Pay,
		HireDate:  time.UnixMilli(employee.Metadata.HireDate),
		StartDate: time.UnixMilli(employee.Metadata.StartDate),
		SickTime:  employee.Metadata.SickTime,
		TimeOff:   employee.Metadata.TimeOff,
		Exempt:    employee.Metadata.Exempt,
		Status:    models.EmployeeStatus(employee.Metadata.Status),
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

	err = a.employeeMicro.UpdateEmployee(ctx, employeeID, models.Employee{
		ReportsToID: reportsToID,
		Title:       urd.Data.Title,
	})
	if err != nil {
		return EmployeeData{}, fmt.Errorf("failed to update employee for %s: %w", employeeID, err)
	}

	return urd.Data, nil
}

func (a adminEmployeeEndpoints) UpdateExemptStatus(ctx context.Context, urd UpdateRequestData[bool]) (bool, error) {
	employeeID, err := uuid.Parse(urd.ID)
	if err != nil {
		return false, fmt.Errorf("failed to parse employee ID: %w", err)
	}

	err = a.employeeMicro.UpdateExemptStatus(ctx, employeeID, urd.Data)
	if err != nil {
		return false, fmt.Errorf("failed to update exempt status for %s: %w", employeeID, err)
	}

	return urd.Data, nil
}

func (a adminEmployeeEndpoints) UpdatePay(ctx context.Context, urd UpdateRequestData[models.EmployeePay]) (models.EmployeePay, error) {
	employeeID, err := uuid.Parse(urd.ID)
	if err != nil {
		return models.EmployeePay{}, fmt.Errorf("failed to parse employee ID: %w", err)
	}

	err = a.employeeMicro.UpdatePay(ctx, employeeID, urd.Data)
	if err != nil {
		return models.EmployeePay{}, fmt.Errorf("failed to update employee pay for %s: %w", employeeID, err)
	}

	return urd.Data, nil
}

func (a adminEmployeeEndpoints) UpdateSickTimeHours(ctx context.Context, urd UpdateRequestData[float64]) (float64, error) {
	employeeID, err := uuid.Parse(urd.ID)
	if err != nil {
		return -1.0, fmt.Errorf("failed to parse employee ID: %w", err)
	}

	err = a.employeeMicro.UpdateSickTime(ctx, employeeID, urd.Data)
	if err != nil {
		return -1.0, fmt.Errorf("failed to update sick time hours for %s: %w", employeeID, err)
	}

	return urd.Data, nil
}

func (a adminEmployeeEndpoints) UpdateStatus(ctx context.Context, urd UpdateRequestData[int]) (int, error) {
	employeeID, err := uuid.Parse(urd.ID)
	if err != nil {
		return -1, fmt.Errorf("failed to parse employee ID: %w", err)
	}

	rawStatus := urd.Data
	if rawStatus < int(models.EmployeeStatusActive) && rawStatus > int(models.EmployeeStatusInactive) {
		return -1, fmt.Errorf("%d is an invalid status value", rawStatus)
	}

	err = a.employeeMicro.UpdateStatus(ctx, employeeID, models.EmployeeStatus(rawStatus))
	if err != nil {
		return -1, fmt.Errorf("failed to update active status for %s: %w", employeeID, err)
	}

	return urd.Data, nil
}

func (a adminEmployeeEndpoints) UpdateTimeOffHours(ctx context.Context, urd UpdateRequestData[float64]) (float64, error) {
	employeeID, err := uuid.Parse(urd.ID)
	if err != nil {
		return -1.0, fmt.Errorf("failed to parse employee ID: %w", err)
	}

	err = a.employeeMicro.UpdateTimeOff(ctx, employeeID, urd.Data)
	if err != nil {
		return -1.0, fmt.Errorf("failed to update time off hours for %s: %w", employeeID, err)
	}

	return urd.Data, nil
}
