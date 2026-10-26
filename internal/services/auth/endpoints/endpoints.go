package endpoints

import (
	"database/sql"

	"github.com/williabk198/timeclock/internal/datastores"
	"github.com/williabk198/timeclock/internal/services/auth"
)

type Endpoints interface {
	AuthEndpoints
}

func NewAuthEndpointHandlers(dbSession *sql.DB) Endpoints {
	authService := auth.NewService(
		datastores.NewUserStore(dbSession),
		datastores.NewTokenStore(dbSession),
	)

	return authEndpoints{
		userMicro:  authService.User(),
		tokenMicro: authService.Token(),
	}
}
