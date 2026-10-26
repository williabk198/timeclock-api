package auth

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/williabk198/timeclock/internal/models"
	"github.com/williabk198/timeclock/internal/utils"
)

func Test_tokenMicroImpl_IssueTokenPair(t *testing.T) {
	testUserID := uuid.New()

	testTokenStore := &mockTokenStore{}
	testTokenStore.On("Add", mock.Anything, mock.MatchedBy(func(tok models.Token) bool {
		return tok.Kind == models.TokenKindAccess
	})).Return(uuid.New(), error(nil))
	testTokenStore.On("Add", mock.Anything, mock.MatchedBy(func(tok models.Token) bool {
		return tok.Kind == models.TokenKindRefresh
	})).Return(uuid.New(), error(nil))

	tmi := tokenMicroImpl{tokenStore: testTokenStore}
	access, refresh, err := tmi.IssueTokenPair(context.Background(), testUserID, "")

	assert.NoError(t, err)
	assert.NotEmpty(t, access.Value)
	assert.NotEmpty(t, refresh.Value)
	assert.NotEqual(t, access.Value, refresh.Value)
	assert.Equal(t, DefaultScope, access.Scope)
	assert.Equal(t, models.TokenKindAccess, access.Kind)
	assert.Equal(t, models.TokenKindRefresh, refresh.Kind)
	assert.True(t, access.ExpiresAt.Before(refresh.ExpiresAt))
	testTokenStore.AssertExpectations(t)
}

func Test_tokenMicroImpl_Introspect(t *testing.T) {
	type args struct {
		ctx      context.Context
		rawToken string
	}

	testUserID := uuid.New()
	now := time.Now().UTC()

	activeRaw := "active-raw-token"
	activeRecord := models.Token{
		ID:        uuid.New(),
		UserID:    testUserID,
		Hash:      utils.HashToken(activeRaw),
		Scope:     "default",
		Kind:      models.TokenKindAccess,
		IssuedAt:  now.Add(-time.Minute),
		ExpiresAt: now.Add(time.Minute),
	}

	expiredRaw := "expired-raw-token"
	expiredRecord := models.Token{
		ID:        uuid.New(),
		UserID:    testUserID,
		Hash:      utils.HashToken(expiredRaw),
		Scope:     "default",
		Kind:      models.TokenKindAccess,
		IssuedAt:  now.Add(-time.Hour),
		ExpiresAt: now.Add(-time.Minute),
	}

	tests := []struct {
		name       string
		setupMock  func(*mockTokenStore)
		args       args
		wantActive bool
		assertion  assert.ErrorAssertionFunc
	}{
		{
			name: "Active",
			setupMock: func(mts *mockTokenStore) {
				mts.On("GetByHash", mock.Anything, utils.HashToken(activeRaw)).Return(activeRecord, error(nil))
			},
			args:       args{ctx: context.Background(), rawToken: activeRaw},
			wantActive: true,
			assertion:  assert.NoError,
		},
		{
			name: "Expired",
			setupMock: func(mts *mockTokenStore) {
				mts.On("GetByHash", mock.Anything, utils.HashToken(expiredRaw)).Return(expiredRecord, error(nil))
			},
			args:       args{ctx: context.Background(), rawToken: expiredRaw},
			wantActive: false,
			assertion:  assert.NoError,
		},
		{
			name: "Unknown",
			setupMock: func(mts *mockTokenStore) {
				mts.On("GetByHash", mock.Anything, utils.HashToken("unknown-raw-token")).Return(models.Token{}, sql.ErrNoRows)
			},
			args:       args{ctx: context.Background(), rawToken: "unknown-raw-token"},
			wantActive: false,
			assertion:  assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testTokenStore := &mockTokenStore{}
			tt.setupMock(testTokenStore)

			tmi := tokenMicroImpl{tokenStore: testTokenStore}
			_, active, err := tmi.Introspect(tt.args.ctx, tt.args.rawToken)
			tt.assertion(t, err)
			assert.Equal(t, tt.wantActive, active)
			testTokenStore.AssertExpectations(t)
		})
	}
}

func Test_tokenMicroImpl_Refresh(t *testing.T) {
	testUserID := uuid.New()
	now := time.Now().UTC()

	t.Run("Success", func(t *testing.T) {
		rawRefresh := "valid-refresh-token"
		refreshID := uuid.New()
		refreshRecord := models.Token{
			ID:        refreshID,
			UserID:    testUserID,
			Hash:      utils.HashToken(rawRefresh),
			Scope:     "default",
			Kind:      models.TokenKindRefresh,
			IssuedAt:  now.Add(-time.Minute),
			ExpiresAt: now.Add(time.Hour),
		}

		testTokenStore := &mockTokenStore{}
		testTokenStore.On("GetByHash", mock.Anything, utils.HashToken(rawRefresh)).Return(refreshRecord, error(nil))
		testTokenStore.On("Revoke", mock.Anything, refreshID).Return(error(nil))
		testTokenStore.On("Add", mock.Anything, mock.AnythingOfType("models.Token")).Return(uuid.New(), error(nil))

		tmi := tokenMicroImpl{tokenStore: testTokenStore}
		access, refresh, err := tmi.Refresh(context.Background(), rawRefresh)

		assert.NoError(t, err)
		assert.NotEmpty(t, access.Value)
		assert.NotEmpty(t, refresh.Value)
		testTokenStore.AssertExpectations(t)
	})

	t.Run("WrongTokenKind", func(t *testing.T) {
		rawAccess := "an-access-token-not-a-refresh-token"
		accessRecord := models.Token{
			ID:        uuid.New(),
			UserID:    testUserID,
			Hash:      utils.HashToken(rawAccess),
			Kind:      models.TokenKindAccess,
			IssuedAt:  now.Add(-time.Minute),
			ExpiresAt: now.Add(time.Hour),
		}

		testTokenStore := &mockTokenStore{}
		testTokenStore.On("GetByHash", mock.Anything, utils.HashToken(rawAccess)).Return(accessRecord, error(nil))

		tmi := tokenMicroImpl{tokenStore: testTokenStore}
		_, _, err := tmi.Refresh(context.Background(), rawAccess)

		assert.ErrorIs(t, err, ErrInvalidToken)
		testTokenStore.AssertExpectations(t)
	})

	t.Run("UnknownToken", func(t *testing.T) {
		testTokenStore := &mockTokenStore{}
		testTokenStore.On("GetByHash", mock.Anything, utils.HashToken("unknown-token")).Return(models.Token{}, sql.ErrNoRows)

		tmi := tokenMicroImpl{tokenStore: testTokenStore}
		_, _, err := tmi.Refresh(context.Background(), "unknown-token")

		assert.ErrorIs(t, err, ErrInvalidToken)
		testTokenStore.AssertExpectations(t)
	})
}

func Test_tokenMicroImpl_Revoke(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		testTokenStore := &mockTokenStore{}
		testTokenStore.On("RevokeByHash", mock.Anything, utils.HashToken("some-token")).Return(error(nil))

		tmi := tokenMicroImpl{tokenStore: testTokenStore}
		err := tmi.Revoke(context.Background(), "some-token")

		assert.NoError(t, err)
		testTokenStore.AssertExpectations(t)
	})

	t.Run("EmptyToken", func(t *testing.T) {
		tmi := tokenMicroImpl{tokenStore: &mockTokenStore{}}
		err := tmi.Revoke(context.Background(), "")

		assert.ErrorIs(t, err, ErrInvalidRequest)
	})
}
