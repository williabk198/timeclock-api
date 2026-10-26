package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/williabk198/timeclock/internal/datastores"
	"github.com/williabk198/timeclock/internal/models"
)

// maxFailedLoginAttempts is the number of consecutive failed login attempts allowed
// before a user's account is automatically locked.
const maxFailedLoginAttempts = 5

type UserMicro interface {
	// Authenticate verifies the given email/password pair against the stored user
	// record and returns the matching user if, and only if, the credentials are valid
	// and the account is enabled.
	Authenticate(ctx context.Context, email, password string) (models.User, error)
	// Register creates a new set of login credentials for an existing person. The
	// password is hashed before being persisted; it is never stored in plain text.
	Register(ctx context.Context, personID uuid.UUID, email, password string) (uuid.UUID, error)
	// Scopes returns the full, space-separated scope string for a user: the "default"
	// scope every authenticated user carries, plus whatever has been explicitly granted
	// to them. This is what gets baked into a token at issuance time.
	Scopes(ctx context.Context, userID uuid.UUID) (string, error)
	// GrantScope grants a user a scope (e.g. "user:admin"), which will be included on
	// tokens issued for them from that point on.
	GrantScope(ctx context.Context, userID uuid.UUID, scope string) error
	// RevokeScope revokes a previously granted scope from a user. Tokens already issued
	// to the user keep whatever scope they were issued with until they expire or are
	// refreshed.
	RevokeScope(ctx context.Context, userID uuid.UUID, scope string) error
}

type userMicroImpl struct {
	userStore datastores.UserDatastore
}

// Authenticate implements UserMicro.
func (umi userMicroImpl) Authenticate(ctx context.Context, email, password string) (models.User, error) {
	user, err := umi.userStore.GetByEmail(ctx, email)
	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, ErrInvalidCredentials
	} else if err != nil {
		return models.User{}, fmt.Errorf("failed to fetch user by email: %w", err)
	}

	metadata, err := umi.userStore.GetMetadata(ctx, user.ID)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to fetch user metadata: %w", err)
	}

	switch metadata.Status {
	case models.UserStatusLocked:
		return models.User{}, ErrAccountLocked
	case models.UserStatusDisabled:
		return models.User{}, ErrAccountDisabled
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPass), []byte(password)); err != nil {
		if recordErr := umi.recordFailedLogin(ctx, user.ID, metadata); recordErr != nil {
			return models.User{}, fmt.Errorf("failed to record failed login attempt: %w", recordErr)
		}
		return models.User{}, ErrInvalidCredentials
	}

	if err := umi.userStore.RecordSuccessfulLogin(ctx, user.ID, time.Now().UTC()); err != nil {
		return models.User{}, fmt.Errorf("failed to record successful login: %w", err)
	}

	return user, nil
}

func (umi userMicroImpl) recordFailedLogin(ctx context.Context, userID uuid.UUID, metadata models.UserMetadata) error {
	metadata.FailedLoginCount++
	if metadata.FailedLoginCount >= maxFailedLoginAttempts {
		metadata.Status = models.UserStatusLocked
	}

	return umi.userStore.UpdateMetadata(ctx, userID, metadata)
}

// Register implements UserMicro.
func (umi userMicroImpl) Register(ctx context.Context, personID uuid.UUID, email, password string) (uuid.UUID, error) {
	if personID == uuid.Nil || email == "" || password == "" {
		return uuid.Nil, fmt.Errorf("%w: personID, email, and password are required", ErrInvalidRequest)
	}

	_, err := umi.userStore.GetByEmail(ctx, email)
	switch {
	case err == nil:
		return uuid.Nil, ErrEmailInUse
	case !errors.Is(err, sql.ErrNoRows):
		return uuid.Nil, fmt.Errorf("failed to check for existing user: %w", err)
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := models.User{
		PersonID:   personID,
		Email:      email,
		HashedPass: string(hashedPass),
	}
	metadata := models.UserMetadata{
		CreatedAt:         time.Now().UTC(),
		LastPasswordReset: time.Now().UTC(),
		Status:            models.UserStatusEnabled,
	}

	id, err := umi.userStore.Add(ctx, user, metadata)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to add user to database: %w", err)
	}

	return id, nil
}

// Scopes implements UserMicro.
func (umi userMicroImpl) Scopes(ctx context.Context, userID uuid.UUID) (string, error) {
	granted, err := umi.userStore.GetGrantedScopes(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch granted scopes: %w", err)
	}

	scopes := append([]string{DefaultScope}, granted...)
	return strings.Join(scopes, " "), nil
}

// GrantScope implements UserMicro.
func (umi userMicroImpl) GrantScope(ctx context.Context, userID uuid.UUID, scope string) error {
	if scope == "" {
		return fmt.Errorf("%w: scope is required", ErrInvalidRequest)
	}

	return umi.userStore.GrantScope(ctx, userID, scope)
}

// RevokeScope implements UserMicro.
func (umi userMicroImpl) RevokeScope(ctx context.Context, userID uuid.UUID, scope string) error {
	return umi.userStore.RevokeScope(ctx, userID, scope)
}
