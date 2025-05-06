package models

import (
	"time"

	"github.com/google/uuid"
)

type Person struct {
	ID          uuid.UUID
	Name        Name      `db:",inline"`
	DateOfBirth time.Time `db:"dob"`
	Gender      Gender
	Pronouns    Pronouns
}

type FirstName string

const (
	FirstNameGiven  FirstName = "given"
	FirstNameFamily FirstName = "family"
)

type Name struct {
	GivenName       string    `db:"given_name" json:"givenName"`
	FamilyName      string    `db:"family_name" json:"familyName"`
	FamilyNameFirst FirstName `db:"first_name" json:"firstName"`
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

// TODO: Create `Marshal` and `Unmarshal` functions for Pronouns as it should be stored as one property in the DB.
// For example if the values here are "he" and "him", then the DB entry should be "he/him"
