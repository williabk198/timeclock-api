package admin

import (
	"context"

	"github.com/google/uuid"
	"github.com/williabk198/timeclock/internal/datastores"
	"github.com/williabk198/timeclock/internal/models"
)

type PersonMicro interface {
	Add(ctx context.Context, person models.Person) (uuid.UUID, error)
	Delete(ctx context.Context, id uuid.UUID) (models.Person, error)
	GetAll(ctx context.Context, offset uint, limit uint) ([]models.Person, error)
	GetSpecific(ctx context.Context, id uuid.UUID) (models.Person, error)
	Update(ctx context.Context, id uuid.UUID, newVal models.Person) error
}

type personMicroImpl struct {
	personStore datastores.PersonStore
}

func (pmi personMicroImpl) Add(ctx context.Context, person models.Person) (uuid.UUID, error) {
	return pmi.personStore.Add(ctx, person)
}

func (pmi personMicroImpl) Delete(ctx context.Context, id uuid.UUID) (models.Person, error) {
	return pmi.personStore.Delete(ctx, id)
}

func (pmi personMicroImpl) GetAll(ctx context.Context, offset uint, limit uint) ([]models.Person, error) {
	return pmi.personStore.GetAllPaginated(ctx, offset, limit)
}

func (pmi personMicroImpl) GetSpecific(ctx context.Context, id uuid.UUID) (models.Person, error) {
	return pmi.personStore.GetSpecific(ctx, id)
}

func (pmi personMicroImpl) Update(ctx context.Context, id uuid.UUID, newVal models.Person) error {
	return pmi.personStore.Update(ctx, id, newVal)
}
