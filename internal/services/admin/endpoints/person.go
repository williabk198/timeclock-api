package endpoints

import (
	"context"

	"github.com/williabk198/timeclock/internal/services/admin"
)

type adminPersonEndpoints struct {
	adminService admin.Service
}

// Add implements PersonEndpoints.
func (ape adminPersonEndpoints) Add(ctx context.Context, person PersonData) (PersonData, error) {
	panic("unimplemented")
}
