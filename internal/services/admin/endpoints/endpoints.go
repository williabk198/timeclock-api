package endpoints

import (
	"context"
	"database/sql"

	"github.com/williabk198/timeclock/internal/datastores"
	"github.com/williabk198/timeclock/internal/services/admin"
)

type Endpoints interface {
	Person() PersonEndpoints
}

type PersonEndpoints interface {
	Add(ctx context.Context, person PersonData) (PersonData, error)
	Delete(ctx context.Context, id string) (PersonData, error)
	GetSpecific(ctx context.Context, id string) (PersonData, error)
	Update(ctx context.Context, urd UpdateRequestData[PersonData]) (PersonData, error)
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

func NewAdminEndpointHandlers(dbSession *sql.DB) Endpoints {
	return adminEndpoints{
		adminService: admin.NewService(datastores.NewPersonStore(dbSession)),
	}
}

type UpdateRequestData[T any] struct {
	ID   string
	Data T
}
