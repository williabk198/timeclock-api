package auth

import (
	"github.com/williabk198/timeclock/internal/datastores"
)

type Service interface {
	User() UserMicro
	Token() TokenMicro
}

type authService struct {
	userStore  datastores.UserDatastore
	tokenStore datastores.TokenDatastore
}

// User implements Service.
func (a authService) User() UserMicro {
	return userMicroImpl{
		userStore: a.userStore,
	}
}

// Token implements Service.
func (a authService) Token() TokenMicro {
	return tokenMicroImpl{
		tokenStore: a.tokenStore,
	}
}

func NewService(userStore datastores.UserDatastore, tokenStore datastores.TokenDatastore) Service {
	return authService{
		userStore:  userStore,
		tokenStore: tokenStore,
	}
}
