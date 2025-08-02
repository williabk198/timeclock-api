package admin

import (
	"context"

	"github.com/google/uuid"
	"github.com/williabk198/timeclock/internal/datastores"
	"github.com/williabk198/timeclock/internal/models"
)

type Service interface {
	AddPerson(context.Context, models.Person) (uuid.UUID, error)
	GetPerson(context.Context, uuid.UUID) (models.Person, error)
}

type adminService struct {
	personStore datastores.PersonStore
}

// AddPunch implements Service.
func (as adminService) AddPerson(ctx context.Context, person models.Person) (uuid.UUID, error) {
	return as.personStore.Add(ctx, person)
}

// GetPerson implements Service.
func (as adminService) GetPerson(ctx context.Context, id uuid.UUID) (models.Person, error) {
	return as.personStore.GetSpecific(ctx, id)
}

func NewService(personStore datastores.PersonStore) Service {
	return adminService{
		personStore: personStore,
	}
}
