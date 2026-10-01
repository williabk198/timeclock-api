package endpoints

import (
	"database/sql"

	"github.com/williabk198/timeclock/internal/datastores"
	"github.com/williabk198/timeclock/internal/services/core"
)

type Endpoints interface {
	Punch() PunchEndpoints
}

type coreEndpoints struct {
	coreService core.Service
}

// Punch implements [Endpoints].
func (c coreEndpoints) Punch() PunchEndpoints {
	return corePunchEndpoints{}
}

func NewCoreEndpoints(dbSession *sql.DB) Endpoints {
	return coreEndpoints{
		coreService: core.NewCoreServices(
			datastores.NewSqlPunchStore(dbSession),
		),
	}
}
