package endpoints

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/williabk198/timeclock/internal/models"
	"github.com/williabk198/timeclock/internal/services/auth"
)

type AuthEndpoints interface {
	Register(ctx context.Context, reqData RegisterUserData) (RegisterUserResponse, error)
	Login(ctx context.Context, reqData LoginRequestData) (LoginResponseData, error)
	Introspect(ctx context.Context, reqData IntrospectRequestData) (IntrospectResponseData, error)
	Refresh(ctx context.Context, reqData RefreshRequestData) (LoginResponseData, error)
	Revoke(ctx context.Context, reqData RevokeRequestData) (RevokeResponseData, error)
	GrantScope(ctx context.Context, reqData GrantScopeRequestData) (ScopeResponseData, error)
	RevokeScope(ctx context.Context, reqData RevokeScopeRequestData) (ScopeResponseData, error)
}

type authEndpoints struct {
	userMicro  auth.UserMicro
	tokenMicro auth.TokenMicro
}

// Register implements AuthEndpoints.
func (ae authEndpoints) Register(ctx context.Context, reqData RegisterUserData) (RegisterUserResponse, error) {
	personID, err := uuid.Parse(reqData.PersonID)
	if err != nil {
		return RegisterUserResponse{}, fmt.Errorf("%w: invalid personID: %v", auth.ErrInvalidRequest, err)
	}

	id, err := ae.userMicro.Register(ctx, personID, reqData.Email, reqData.Password)
	if err != nil {
		return RegisterUserResponse{}, err
	}

	return RegisterUserResponse{ID: id.String()}, nil
}

// Login implements AuthEndpoints.
func (ae authEndpoints) Login(ctx context.Context, reqData LoginRequestData) (LoginResponseData, error) {
	user, err := ae.userMicro.Authenticate(ctx, reqData.Email, reqData.Password)
	if err != nil {
		return LoginResponseData{}, err
	}

	scope, err := ae.userMicro.Scopes(ctx, user.ID)
	if err != nil {
		return LoginResponseData{}, fmt.Errorf("failed to resolve user scopes: %w", err)
	}

	access, refresh, err := ae.tokenMicro.IssueTokenPair(ctx, user.ID, scope)
	if err != nil {
		return LoginResponseData{}, fmt.Errorf("failed to issue tokens: %w", err)
	}

	return LoginResponseData{
		AccessToken:  issuedTokenToData(access),
		RefreshToken: issuedTokenToData(refresh),
	}, nil
}

// Introspect implements AuthEndpoints.
func (ae authEndpoints) Introspect(ctx context.Context, reqData IntrospectRequestData) (IntrospectResponseData, error) {
	record, active, err := ae.tokenMicro.Introspect(ctx, reqData.Token)
	if err != nil {
		return IntrospectResponseData{}, err
	}
	if !active {
		return IntrospectResponseData{Active: false}, nil
	}

	return IntrospectResponseData{
		Active:    true,
		Scope:     record.Scope,
		TokenType: record.Kind.String(),
		Subject:   record.UserID.String(),
		IssuedAt:  record.IssuedAt.Unix(),
		ExpiresAt: record.ExpiresAt.Unix(),
	}, nil
}

// Refresh implements AuthEndpoints.
func (ae authEndpoints) Refresh(ctx context.Context, reqData RefreshRequestData) (LoginResponseData, error) {
	access, refresh, err := ae.tokenMicro.Refresh(ctx, reqData.RefreshToken)
	if err != nil {
		return LoginResponseData{}, err
	}

	return LoginResponseData{
		AccessToken:  issuedTokenToData(access),
		RefreshToken: issuedTokenToData(refresh),
	}, nil
}

// Revoke implements AuthEndpoints.
func (ae authEndpoints) Revoke(ctx context.Context, reqData RevokeRequestData) (RevokeResponseData, error) {
	if err := ae.tokenMicro.Revoke(ctx, reqData.Token); err != nil {
		return RevokeResponseData{}, err
	}

	return RevokeResponseData{Revoked: true}, nil
}

// GrantScope implements AuthEndpoints.
func (ae authEndpoints) GrantScope(ctx context.Context, reqData GrantScopeRequestData) (ScopeResponseData, error) {
	userID, err := uuid.Parse(reqData.UserID)
	if err != nil {
		return ScopeResponseData{}, fmt.Errorf("%w: invalid userID: %v", auth.ErrInvalidRequest, err)
	}

	if err := ae.userMicro.GrantScope(ctx, userID, reqData.Scope); err != nil {
		return ScopeResponseData{}, err
	}

	return ScopeResponseData{UserID: reqData.UserID, Scope: reqData.Scope}, nil
}

// RevokeScope implements AuthEndpoints.
func (ae authEndpoints) RevokeScope(ctx context.Context, reqData RevokeScopeRequestData) (ScopeResponseData, error) {
	userID, err := uuid.Parse(reqData.UserID)
	if err != nil {
		return ScopeResponseData{}, fmt.Errorf("%w: invalid userID: %v", auth.ErrInvalidRequest, err)
	}

	if err := ae.userMicro.RevokeScope(ctx, userID, reqData.Scope); err != nil {
		return ScopeResponseData{}, err
	}

	return ScopeResponseData{UserID: reqData.UserID, Scope: reqData.Scope}, nil
}

func issuedTokenToData(t models.IssuedToken) TokenData {
	return TokenData{
		Token:     t.Value,
		Scope:     t.Scope,
		Kind:      t.Kind.String(),
		IssuedAt:  t.IssuedAt.Unix(),
		ExpiresAt: t.ExpiresAt.Unix(),
	}
}
