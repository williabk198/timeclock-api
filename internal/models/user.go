package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `jagsqlb:"id;omit"`
	PersonID   uuid.UUID `jagsqlb:"person_id;omit-update"`
	Email      string    `jagsqlb:"email"`
	HashedPass string    `jagsqlb:"pass"`
}

type UserStatus int

const (
	UserStatusEnabled UserStatus = iota + 1
	UserStatusDisabled
	UserStatusLocked
)

type UserMetadata struct {
	UserID            uuid.UUID  `jagsqlb:"user_id;omit-update"`
	CreatedAt         time.Time  `jagsqlb:"created_at;omit-update"`
	FailedLoginCount  int        `jagsqlb:"failure_count"`
	LastLogin         time.Time  `jagsqlb:"last_login"`
	LastPasswordReset time.Time  `jagsqlb:"last_pw_reset"`
	Status            UserStatus `jagsqlb:"status"`
}
