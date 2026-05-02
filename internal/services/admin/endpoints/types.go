package endpoints

import "github.com/williabk198/timeclock/internal/models"

type PersonData struct {
	ID          string             `json:"id"`
	Name        models.Name        `json:"name"`
	DateOfBirth int64              `json:"dob"` // UNIX Timestamp
	Gender      string             `json:"gender"`
	Pronouns    string             `json:"pronouns"`
	Contacts    *PersonContactData `json:"contacts,omitempty"`
}

type PersonContactData struct {
	Addresses    []PersonAddressData `json:"addresses"`
	Emails       []PersonEmailData   `json:"emails"`
	PhoneNumbers []PersonPhoneData   `json:"phoneNumbers"`
}

type PersonAddressData struct {
	ID         string `json:"id"`
	Street1    string `json:"street1"`
	Street2    string `json:"street2"`
	Locality   string `json:"locality"` // Locality is the name of the city, town, etc..
	Region     string `json:"region"`   // Region is the name of the State, Province, Prefecture, etc...
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
	Type       string `json:"type"` // Type holds what kind of address this value represents: Mailing, Physical, Billing, etc...
	Primary    bool   `json:"primary"`
}

type PersonEmailData struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Primary bool   `json:"primary"`
}

type PersonPhoneData struct {
	ID          string `json:"id"`
	CountryCode int    `json:"countryCode"`
	PhoneNumber string `json:"phoneNumber"`
	Type        string `json:"type"`
	Primary     bool   `json:"primary"`
}

type EmployeeData struct {
	ID          string            `json:"id"`
	PersonID    string            `json:"personID"`
	ReportsToID string            `json:"reportsToID"`
	Title       string            `json:"title"`
	Metadata    *EmployeeMetadata `json:"metadata,omitempty"`
}

type EmployeeMetadata struct {
	Pay       models.EmployeePay `json:"pay"`
	HireDate  int64              `json:"hireDate"`  // UNIX Timestamp
	StartDate int64              `json:"startDate"` // UNIX Timestamp
	SickTime  float32            `json:"sickTime"`
	TimeOff   float32            `json:"timeOff"`
	Exempt    bool               `json:"exempt"`
	Status    int                `json:"status"`
}
