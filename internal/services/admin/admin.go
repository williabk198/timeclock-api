package admin

import (
	"context"

	"github.com/google/uuid"
	"github.com/williabk198/timeclock/internal/datastores"
	"github.com/williabk198/timeclock/internal/models"
)

type Service interface {
	AddPerson(context.Context, models.Person) (uuid.UUID, error)
	DeletePerson(context.Context, uuid.UUID) (models.Person, error)
	GetAllPersons(ctx context.Context, offset uint, limit uint) ([]models.Person, error)
	GetPerson(context.Context, uuid.UUID) (models.Person, error)
	UpdatePerson(context.Context, uuid.UUID, models.Person) error
}

type adminService struct {
	personStore datastores.PersonStore
}

// AddPunch implements Service.
func (as adminService) AddPerson(ctx context.Context, person models.Person) (uuid.UUID, error) {
	return as.personStore.Add(ctx, person)
}

// DeletePerson implements Service.
func (as adminService) DeletePerson(ctx context.Context, id uuid.UUID) (models.Person, error) {
	return as.personStore.Delete(ctx, id)
}

// GetAllPersons implements Service.
func (as adminService) GetAllPersons(ctx context.Context, offset uint, limit uint) ([]models.Person, error) {
	return as.personStore.GetAllPaginated(ctx, offset, limit)
}

// GetPerson implements Service.
func (as adminService) GetPerson(ctx context.Context, id uuid.UUID) (models.Person, error) {
	return as.personStore.GetSpecific(ctx, id)
}

// UpdatePerson implements Service.
func (as adminService) UpdatePerson(ctx context.Context, id uuid.UUID, data models.Person) error {
	return as.personStore.Update(ctx, id, data)
}

func NewService(personStore datastores.PersonStore) Service {
	return adminService{
		personStore: personStore,
	}
}
