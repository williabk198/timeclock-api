package datastores

import (
	"context"

	"github.com/google/uuid"
)

type SqlIdentifier interface {
	~int | uuid.UUID
}

type SqlDatastore[T any, U SqlIdentifier] interface {
	Add(ctx context.Context, item T) (id U, err error)
	Delete(ctx context.Context, id U) (item T, err error)
	GetAllPaginated(ctx context.Context, offset uint, limit uint) (items []T, err error)
	GetSpecific(ctx context.Context, id U) (item T, err error)
	Update(ctx context.Context, id U, item T) (err error)
}
