package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Person struct {
	ID          uuid.UUID `jagsqlb:"id;omit"`
	Name        Name      `jagsqlb:";inline"`
	DateOfBirth time.Time `jagsqlb:"dob"`
	Gender      Gender    `jagsqlb:"gender"`
	Pronouns    Pronouns  `jagsqlb:"pronouns"`
}

type FirstName string

const (
	FirstNameGiven  FirstName = "given"
	FirstNameFamily FirstName = "family"
)

type Name struct {
	GivenName       string    `jagsqlb:"given_name" json:"givenName"`
	FamilyName      string    `jagsqlb:"family_name" json:"familyName"`
	FamilyNameFirst FirstName `jagsqlb:"first_name" json:"firstName"`
}

type Gender string

const (
	// Define some of the most common genders
	GenderMale      Gender = "male"
	GenderFemale    Gender = "female"
	GenderTrans     Gender = "transgender"
	GenderNonBinary Gender = "non-binary"
)

// Pronouns represents a how a person identifies themselves. Where subject would be he, she, they, etc...,
// and Object would be him, her, their, etc...
type Pronouns struct {
	Subject string
	Object  string
}

func (p Pronouns) String() string {
	return fmt.Sprintf("%s/%s", p.Subject, p.Object)
}

func (p Pronouns) MarshalQuery() (string, error) {
	return p.String(), nil
}

type ContactAddress struct {
	ID         uuid.UUID `jagsqlb:"id;omit"`
	PersonID   uuid.UUID `jagsqlb:"person_id"`
	Street1    string    `jagsqlb:"street1" json:"street1"`
	Street2    string    `jagsqlb:"street2" json:"street2"`
	Locality   string    `jagsqlb:"city" json:"locality"`
	Region     string    `jagsqlb:"locality" json:"region"` // Locality is the State, Province, Prefecture, etc... that the City is located in
	PostalCode string    `jagsqlb:"postal_code" json:"postalCode"`
	Country    string    `jagsqlb:"country" json:"country"`
	Type       string    `jagsqlb:"kind" json:"type"` // Type holds what kind of address this value represents: Mailing, Physical, Billing, etc...
	Primary    bool      `jagsqlb:"primary" json:"primary"`
}

type ContactEmail struct {
	ID       uuid.UUID `jagsqlb:"id;omit"`
	PersonID uuid.UUID `jagsqlb:"person_id"`
	Username string    `jagsqlb:"username"`
	Provider string    `jagsqlb:"provider"`
	Primary  bool      `jagsqlb:"primary" json:"primary"`
}

func (ce ContactEmail) String() string {
	return fmt.Sprintf("%s@%s", ce.Username, ce.Provider)
}

type ContactPhone struct {
	ID          uuid.UUID `jagsqlb:"id;omit"`
	PersonID    uuid.UUID `jagsqlb:"person_id"`
	CountryCode int       `jagsqlb:"country_code"`
	PhoneNumber string    `jagsqlb:"phone_number"`
	Type        string    `jagsqlb:"kind" json:"type"` // Type holds what kind of phone number this value represents: Home, Cell, Work, etc...
	Primary     bool      `jagsqlb:"primary" json:"primary"`
}

func (cp ContactPhone) String() string {
	return fmt.Sprintf("+%d %s", cp.CountryCode, cp.PhoneNumber)
}

type Contacts struct {
	Addresses []ContactAddress
	Email     []ContactEmail
	Phone     []ContactPhone
}
