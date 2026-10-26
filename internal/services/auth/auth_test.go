package auth

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/williabk198/timeclock/internal/datastores"
	"github.com/williabk198/timeclock/internal/models"
)

func TestNewService(t *testing.T) {
	type args struct {
		userStore  datastores.UserDatastore
		tokenStore datastores.TokenDatastore
	}

	testUserStore := &mockUserStore{}
	testTokenStore := &mockTokenStore{}

	tests := []struct {
		name string
		args args
		want Service
	}{
		{
			name: "Success",
			args: args{
				userStore:  testUserStore,
				tokenStore: testTokenStore,
			},
			want: authService{
				userStore:  testUserStore,
				tokenStore: testTokenStore,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewService(tt.args.userStore, tt.args.tokenStore))
		})
	}
}

type mockUserStore struct {
	mock.Mock
}

func (mus *mockUserStore) Add(ctx context.Context, user models.User, metadata models.UserMetadata) (uuid.UUID, error) {
	args := mus.Called(ctx, user, metadata)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (mus *mockUserStore) GetByEmail(ctx context.Context, email string) (models.User, error) {
	args := mus.Called(ctx, email)
	return args.Get(0).(models.User), args.Error(1)
}

func (mus *mockUserStore) GetSpecific(ctx context.Context, id uuid.UUID) (models.User, error) {
	args := mus.Called(ctx, id)
	return args.Get(0).(models.User), args.Error(1)
}

func (mus *mockUserStore) GetMetadata(ctx context.Context, userID uuid.UUID) (models.UserMetadata, error) {
	args := mus.Called(ctx, userID)
	return args.Get(0).(models.UserMetadata), args.Error(1)
}

func (mus *mockUserStore) UpdateMetadata(ctx context.Context, userID uuid.UUID, metadata models.UserMetadata) error {
	args := mus.Called(ctx, userID, metadata)
	return args.Error(0)
}

func (mus *mockUserStore) RecordSuccessfulLogin(ctx context.Context, userID uuid.UUID, loginTime time.Time) error {
	args := mus.Called(ctx, userID, loginTime)
	return args.Error(0)
}

func (mus *mockUserStore) GetGrantedScopes(ctx context.Context, userID uuid.UUID) ([]string, error) {
	args := mus.Called(ctx, userID)
	return args.Get(0).([]string), args.Error(1)
}

func (mus *mockUserStore) GrantScope(ctx context.Context, userID uuid.UUID, scope string) error {
	args := mus.Called(ctx, userID, scope)
	return args.Error(0)
}

func (mus *mockUserStore) RevokeScope(ctx context.Context, userID uuid.UUID, scope string) error {
	args := mus.Called(ctx, userID, scope)
	return args.Error(0)
}

type mockTokenStore struct {
	mock.Mock
}

func (mts *mockTokenStore) Add(ctx context.Context, token models.Token) (uuid.UUID, error) {
	args := mts.Called(ctx, token)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (mts *mockTokenStore) GetByHash(ctx context.Context, hash string) (models.Token, error) {
	args := mts.Called(ctx, hash)
	return args.Get(0).(models.Token), args.Error(1)
}

func (mts *mockTokenStore) Revoke(ctx context.Context, id uuid.UUID) error {
	args := mts.Called(ctx, id)
	return args.Error(0)
}

func (mts *mockTokenStore) RevokeByHash(ctx context.Context, hash string) error {
	args := mts.Called(ctx, hash)
	return args.Error(0)
}
