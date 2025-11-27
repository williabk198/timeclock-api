package endpoints

import (
	"database/sql"

	"github.com/williabk198/timeclock/internal/datastores"
	"github.com/williabk198/timeclock/internal/services/admin"
)

type Endpoints interface {
	Contact() ContactEndpoints
	Employee() EmployeeEndpoints
	Person() PersonEndpoints
}

type adminEndpoints struct {
	adminService admin.Service
}

func (a adminEndpoints) Contact() ContactEndpoints {
	return adminContactEndpoints{
		contactMicro: a.adminService.Contact(),
	}
}

// Person implements Endpoints.
func (a adminEndpoints) Person() PersonEndpoints {
	return adminPersonEndpoints{
		personMicro: a.adminService.Person(),
	}
}

// Employee implements Endpoints.
func (a adminEndpoints) Employee() EmployeeEndpoints {
	return adminEmployeeEndpoints{
		employeeMicro: a.adminService.Employee(),
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
