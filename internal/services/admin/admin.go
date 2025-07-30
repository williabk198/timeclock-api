package admin

import (
	"context"

	"github.com/google/uuid"
	"github.com/williabk198/timeclock/internal/datastores"
	"github.com/williabk198/timeclock/internal/models"
)

type Service interface {
	AddPerson(context.Context, models.Person) (uuid.UUID, error)
}

type adminService struct {
	personStore datastores.PersonStore
}

// AddPunch implements Service.
func (as adminService) AddPerson(ctx context.Context, person models.Person) (uuid.UUID, error) {
	panic("unimplemented")
}

func NewService(personStore datastores.PersonStore) Service {
	panic("unimplemented")
}
