package models

import (
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

// TODO Create String function for EmployeeStatus

// TODO? Create Parser to turn a string into an EmployeeStatus

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

// TODO: Create `Marshal` and `Unmarshal` functions for EmployeePay as it should be stored as one property in the DB.
// For example if the values here are "USD", 16.00, and "per hour" respectively, then the value in the DB should be
// "16.00 USD/hour"

type Employee struct {
	ID          uuid.UUID `db:"id"`
	PersonID    uuid.UUID `db:"person_id"`
	ReportsToID uuid.UUID `db:"reports_to_id"` // The employeeID of the individual that the employee reports to. Use uuid.Nil to indicate the employee reports to nobody
	Title       string
	Pay         EmployeePay
	Status      EmployeeStatus
}

type EmployeeMetadata struct {
	HireDate time.Time `db:"hire_date"`
	SickTime float32   `db:"sick_time"`
	TimeOff  float32   `db:"time_off"`
	Exempt   bool
}
