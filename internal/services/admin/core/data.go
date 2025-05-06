package core

import "github.com/williabk198/timeclock/internal/models"

type PersonData struct {
	ID          string      `json:"id"`
	Name        models.Name `json:"name"`
	DateOfBirth int64       `json:"dob"` // UNIX Timestamp
	Gender      string      `json:"gender"`
	Pronouns    string      `json:"pronouns"`
}

type EmployeeData struct {
	ID          string `json:"id"`
	PersonID    string `json:"personID"`
	ReportsToID string `json:"reportsToID"`
	Title       string `json:"title"`
	Exempt      bool   `json:"exempt"`
	Status      int    `json:"status"`
}

type EmployeeMetadata struct {
	EmployeeID string  `json:"employeeID"`
	HireDate   int64   `json:"hireDate"` // UNIX Timestamp
	SickTime   float32 `json:"sickTime"`
	TimeOff    float32 `json:"timeOff"`
}
