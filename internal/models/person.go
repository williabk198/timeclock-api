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
