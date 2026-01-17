package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type EmployeeStatus int

const (
	EmployeeStatusActive EmployeeStatus = iota + 1
	EmployeeStatusInactive
	EmployeeStatusGone

	employeeStatusStrActive   string = "active"
	employeeStatusStrInactive string = "inactive"
	employeeStatusStrGone     string = "gone"
)

func (es EmployeeStatus) String() string {
	switch es {
	case EmployeeStatusActive:
		return employeeStatusStrActive
	case EmployeeStatusInactive:
		return employeeStatusStrInactive
	case EmployeeStatusGone:
		return employeeStatusStrGone
	}

	return ""
}

func ParseEmployeeStatus(input string) (EmployeeStatus, error) {
	switch input {
	case employeeStatusStrActive:
		return EmployeeStatusActive, nil
	case employeeStatusStrInactive:
		return EmployeeStatusInactive, nil
	case employeeStatusStrGone:
		return EmployeeStatusGone, nil
	}

	return -1, fmt.Errorf("")
}

type PayCadence string

const (
	PayCadenceHourly PayCadence = "hour"
	PayCadenceYearly PayCadence = "year"
)

type EmployeePay struct {
	Currency string // This value should hold the abbreviation(NOT THE SYMBOL) of the currency (e.g USD, CAD, JPY, etc...)
	Rate     float32
	Cadence  PayCadence
}

func (ep EmployeePay) String() string {
	return fmt.Sprintf("%.2f %s/%s", ep.Rate, ep.Currency, ep.Cadence)
}

func (ep EmployeePay) MarshalQuery() (string, error) {
	return ep.String(), nil
}

type Employee struct {
	ID          uuid.UUID `jagsqlb:"id;omit"`
	PersonID    uuid.UUID `jagsqlb:"person_id;omit-update"`
	ReportsToID uuid.UUID `jagsqlb:"reports_to_eid"` // The employeeID of the individual that the employee reports to. Use uuid.Nil to indicate the employee reports to nobody
	Title       string    `jagsqlb:"title"`
}

type EmployeeMetadata struct {
	EmployeeID uuid.UUID      `jagsqlb:"eid"`
	Pay        EmployeePay    `jagsqlb:"pay"`
	HireDate   time.Time      `jagsqlb:"hire_date"`
	StartDate  time.Time      `jagsqlb:"start_date"`
	SickTime   float32        `jagsqlb:"sick_time"`
	TimeOff    float32        `jagsqlb:"time_off"`
	Exempt     bool           `jagsqlb:"exempt"`
	Status     EmployeeStatus `jagsqlb:"status"`
}
