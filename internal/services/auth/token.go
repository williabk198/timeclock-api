package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/williabk198/timeclock/internal/datastores"
	"github.com/williabk198/timeclock/internal/models"
	"github.com/williabk198/timeclock/internal/utils"
)

const (
	// AccessTokenTTL is how long an issued access token remains valid.
	AccessTokenTTL = 10 * time.Minute
	// RefreshTokenTTL is how long an issued refresh token remains valid.
	RefreshTokenTTL = 10 * time.Hour
	// DefaultScope is the scope assigned to a token when the caller doesn't specify one.
	DefaultScope = "default"
)

// TokenMicro issues, introspects, refreshes, and revokes the opaque tokens used to
// authenticate requests to the other services. Tokens are opaque, random values; they
// carry no embedded information, so validating one always requires a lookup against the
// datastore via Introspect.
type TokenMicro interface {
	// IssueTokenPair generates and persists a new access and refresh token for the given
	// user and scope. The raw token values are only ever available on the returned
	// models.IssuedToken; they can't be recovered later.
	IssueTokenPair(ctx context.Context, userID uuid.UUID, scope string) (access models.IssuedToken, refresh models.IssuedToken, err error)
	// Introspect reports whether the given raw token value corresponds to an active
	// (i.e. not expired or revoked) token, along with the persisted record for it.
	// An unknown token is reported as inactive, not as an error.
	Introspect(ctx context.Context, rawToken string) (token models.Token, active bool, err error)
	// Refresh exchanges a valid, active refresh token for a new access/refresh token
	// pair. The provided refresh token is revoked as part of the exchange (refresh
	// token rotation), so it can't be replayed.
	Refresh(ctx context.Context, rawRefreshToken string) (access models.IssuedToken, refresh models.IssuedToken, err error)
	// Revoke immediately invalidates the given raw token value, regardless of its kind.
	Revoke(ctx context.Context, rawToken string) error
}

type tokenMicroImpl struct {
	tokenStore datastores.TokenDatastore
}

// IssueTokenPair implements TokenMicro.
func (tmi tokenMicroImpl) IssueTokenPair(ctx context.Context, userID uuid.UUID, scope string) (models.IssuedToken, models.IssuedToken, error) {
	if scope == "" {
		scope = DefaultScope
	}

	access, err := tmi.issue(ctx, userID, scope, models.TokenKindAccess, AccessTokenTTL)
	if err != nil {
		return models.IssuedToken{}, models.IssuedToken{}, fmt.Errorf("failed to issue access token: %w", err)
	}

	refresh, err := tmi.issue(ctx, userID, scope, models.TokenKindRefresh, RefreshTokenTTL)
	if err != nil {
		return models.IssuedToken{}, models.IssuedToken{}, fmt.Errorf("failed to issue refresh token: %w", err)
	}

	return access, refresh, nil
}

func (tmi tokenMicroImpl) issue(ctx context.Context, userID uuid.UUID, scope string, kind models.TokenKind, ttl time.Duration) (models.IssuedToken, error) {
	raw, err := utils.GenerateOpaqueToken()
	if err != nil {
		return models.IssuedToken{}, err
	}

	now := time.Now().UTC()
	record := models.Token{
		UserID:    userID,
		Hash:      utils.HashToken(raw),
		Scope:     scope,
		Kind:      kind,
		IssuedAt:  now,
		ExpiresAt: now.Add(ttl),
	}

	if _, err := tmi.tokenStore.Add(ctx, record); err != nil {
		return models.IssuedToken{}, err
	}

	return models.IssuedToken{
		Value:     raw,
		Scope:     scope,
		Kind:      kind,
		IssuedAt:  record.IssuedAt,
		ExpiresAt: record.ExpiresAt,
	}, nil
}

// Introspect implements TokenMicro.
func (tmi tokenMicroImpl) Introspect(ctx context.Context, rawToken string) (models.Token, bool, error) {
	if rawToken == "" {
		return models.Token{}, false, nil
	}

	record, err := tmi.tokenStore.GetByHash(ctx, utils.HashToken(rawToken))
	if errors.Is(err, sql.ErrNoRows) {
		return models.Token{}, false, nil
	} else if err != nil {
		return models.Token{}, false, fmt.Errorf("failed to fetch token by hash: %w", err)
	}

	return record, record.Active(time.Now().UTC()), nil
}

// Refresh implements TokenMicro.
func (tmi tokenMicroImpl) Refresh(ctx context.Context, rawRefreshToken string) (models.IssuedToken, models.IssuedToken, error) {
	record, active, err := tmi.Introspect(ctx, rawRefreshToken)
	if err != nil {
		return models.IssuedToken{}, models.IssuedToken{}, err
	}
	if !active || record.Kind != models.TokenKindRefresh {
		return models.IssuedToken{}, models.IssuedToken{}, ErrInvalidToken
	}

	// Revoke the old refresh token immediately so it can't be replayed, even if issuing
	// the new pair below were to fail partway through.
	if err := tmi.tokenStore.Revoke(ctx, record.ID); err != nil {
		return models.IssuedToken{}, models.IssuedToken{}, fmt.Errorf("failed to revoke used refresh token: %w", err)
	}

	return tmi.IssueTokenPair(ctx, record.UserID, record.Scope)
}

// Revoke implements TokenMicro.
func (tmi tokenMicroImpl) Revoke(ctx context.Context, rawToken string) error {
	if rawToken == "" {
		return fmt.Errorf("%w: token is required", ErrInvalidRequest)
	}

	return tmi.tokenStore.RevokeByHash(ctx, utils.HashToken(rawToken))
}
