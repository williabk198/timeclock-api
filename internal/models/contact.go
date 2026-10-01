package models

import (
	"fmt"

	"github.com/google/uuid"
)

type AddressType string

const (
	AddressTypePhysical AddressType = "physical"
	AddressTypeMailing  AddressType = "mailing"
)

type PhoneType string

const (
	PhoneTypeHome PhoneType = "home"
	PhoneTypeCell PhoneType = "cell"
)

type ContactAddress struct {
	ID         uuid.UUID   `jagsqlb:"id;omit"`
	PersonID   uuid.UUID   `jagsqlb:"person_id;omit-update"`
	Street1    string      `jagsqlb:"street1" json:"street1"`
	Street2    string      `jagsqlb:"street2" json:"street2"`
	Locality   string      `jagsqlb:"locality" json:"locality"`
	Region     string      `jagsqlb:"region" json:"region"` // Locality is the State, Province, Prefecture, etc... that the City is located in
	PostalCode string      `jagsqlb:"postal_code" json:"postalCode"`
	Country    string      `jagsqlb:"country" json:"country"`
	Type       AddressType `jagsqlb:"kind" json:"type"` // Type holds what kind of address this value represents: Mailing, Physical, Billing, etc...
	Primary    bool        `jagsqlb:"primary" json:"primary"`
}

type ContactEmail struct {
	ID       uuid.UUID `jagsqlb:"id;omit"`
	PersonID uuid.UUID `jagsqlb:"person_id;omit-update"`
	Username string    `jagsqlb:"username"`
	Provider string    `jagsqlb:"provider"`
	Primary  bool      `jagsqlb:"primary" json:"primary"`
}

func (ce ContactEmail) String() string {
	return fmt.Sprintf("%s@%s", ce.Username, ce.Provider)
}

type ContactPhone struct {
	ID          uuid.UUID `jagsqlb:"id;omit"`
	PersonID    uuid.UUID `jagsqlb:"person_id;omit-update"`
	CountryCode int       `jagsqlb:"country_code"`
	PhoneNumber string    `jagsqlb:"phone_number"`
	Type        PhoneType `jagsqlb:"kind" json:"type"` // Type holds what kind of phone number this value represents: Home, Cell, Work, etc...
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
