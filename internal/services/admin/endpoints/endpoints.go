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
	AddContactEmailForPerson(ctx context.Context, reqData AddSubItemRequestData[PersonEmailData]) (PersonEmailData, error)
	GetPersonContacts(ctx context.Context, personID string) (PersonContactData, error)
	GetPersonContactAddresses(ctx context.Context, personID string) ([]PersonAddressData, error)
	GetPersonContactEmails(ctx context.Context, personID string) ([]PersonEmailData, error)
	GetPersonContactPhones(ctx context.Context, personID string) ([]PersonPhoneData, error)
}

type adminEndpoints struct {
	adminService admin.Service
}

// Person implements Endpoints.
func (a adminEndpoints) Person() PersonEndpoints {
	return adminPersonEndpoints{
		adminService: a.adminService,
	}
}

func (a adminEndpoints) Contact() ContactEndpoints {
	return adminContactEndpoints{
		adminService: a.adminService,
	}
}

func NewAdminEndpointHandlers(dbSession *sql.DB) Endpoints {
	return adminEndpoints{
		adminService: admin.NewService(datastores.NewPersonStore(dbSession)),
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

type GetPaginatedRequestData struct {
	Offset uint
	Limit  uint
}
