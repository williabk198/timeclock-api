package datastores

import (
	"context"

	"github.com/google/uuid"
	"github.com/williabk198/timeclock/internal/models"
)

type SqlIdentifier interface {
	~int | uuid.UUID
}

type SqlDatastore[T any, U SqlIdentifier] interface {
	Add(ctx context.Context, item models.Person) (id uuid.UUID, err error)
	Delete(ctx context.Context, id uuid.UUID) (item models.Person, err error)
	GetAllPaginated(ctx context.Context, offset uint, limit uint) (items []models.Person, err error)
	GetSpecific(ctx context.Context, id uuid.UUID) (item models.Person, err error)
	Update(ctx context.Context, id uuid.UUID, item models.Person) (err error)
}
