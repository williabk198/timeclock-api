package endpoints

import (
	"context"
	"database/sql"

	"github.com/williabk198/timeclock/internal/datastores"
	"github.com/williabk198/timeclock/internal/services/admin"
)

type Endpoints interface {
	Person() PersonEndpoints
	Contact() ContactEndpoints
}

type PersonEndpoints interface {
	Add(ctx context.Context, person PersonData) (PersonData, error)
	Delete(ctx context.Context, id string) (PersonData, error)
	GetSpecific(ctx context.Context, id string) (PersonData, error)
	GetAll(ctx context.Context, reqData GetPaginatedRequestData) ([]PersonData, error)
	Update(ctx context.Context, urd UpdateRequestData[PersonData]) (PersonData, error)
}

type ContactEndpoints interface {
	AddContactAddressForPerson(ctx context.Context, reqData AddSubItemRequestData[PersonAddressData]) (PersonAddressData, error)
	AddContactEmailForPerson(ctx context.Context, reqData AddSubItemRequestData[PersonEmailData]) (PersonEmailData, error)
	AddContactPhoneForPerson(ctx context.Context, reqData AddSubItemRequestData[PersonPhoneData]) (PersonPhoneData, error)
	DeleteContactAddressForPerson(ctx context.Context, reqData DeleteContactRequestData) (PersonAddressData, error)
	DeleteContactEmailForPerson(ctx context.Context, reqData DeleteContactRequestData) (PersonEmailData, error)
	DeleteContactPhoneForPerson(ctx context.Context, reqData DeleteContactRequestData) (PersonPhoneData, error)
	GetPersonContacts(ctx context.Context, personID string) (PersonContactData, error)
	GetPersonContactAddresses(ctx context.Context, personID string) ([]PersonAddressData, error)
	GetPersonContactEmails(ctx context.Context, personID string) ([]PersonEmailData, error)
	GetPersonContactPhones(ctx context.Context, personID string) ([]PersonPhoneData, error)
	UpdatePersonContactAddress(ctx context.Context, reqData UpdateContactRequestData[PersonAddressData]) (PersonAddressData, error)
	UpdatePersonContactEmail(ctx context.Context, reqData UpdateContactRequestData[PersonEmailData]) (PersonEmailData, error)
	UpdatePersonContactPhone(ctx context.Context, reqData UpdateContactRequestData[PersonPhoneData]) (PersonPhoneData, error)
}

type adminEndpoints struct {
	adminService admin.Service
}

// Person implements Endpoints.
func (a adminEndpoints) Person() PersonEndpoints {
	return adminPersonEndpoints{
		personMicro: a.adminService.Person(),
	}
}

func (a adminEndpoints) Contact() ContactEndpoints {
	return adminContactEndpoints{
		contactMicro: a.adminService.Contact(),
	}
}

func NewAdminEndpointHandlers(dbSession *sql.DB) Endpoints {
	return adminEndpoints{
		adminService: admin.NewService(datastores.NewPersonStore(dbSession), datastores.NewContactStore(dbSession)),
	}
}

type AddSubItemRequestData[T any] struct {
	ParentID string
	Data     T
}

type UpdateRequestData[T any] struct {
	ID   string
	Data T
}

type DeleteContactRequestData struct {
	PerosnID  string
	ContactID string
}

type UpdateContactRequestData[T ContactConstraint] struct {
	PersonID  string
	ContactID string
	Data      T
}

type GetPaginatedRequestData struct {
	Offset uint
	Limit  uint
}

type ContactConstraint interface {
	PersonAddressData | PersonEmailData | PersonPhoneData
}
