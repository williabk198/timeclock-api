package endpoints

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/williabk198/timeclock/internal/models"
)

type mockUserMicro struct {
	mock.Mock
}

func (m *mockUserMicro) Authenticate(ctx context.Context, email, password string) (models.User, error) {
	args := m.Called(ctx, email, password)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *mockUserMicro) Register(ctx context.Context, personID uuid.UUID, email, password string) (uuid.UUID, error) {
	args := m.Called(ctx, personID, email, password)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *mockUserMicro) Scopes(ctx context.Context, userID uuid.UUID) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *mockUserMicro) GrantScope(ctx context.Context, userID uuid.UUID, scope string) error {
	args := m.Called(ctx, userID, scope)
	return args.Error(0)
}

func (m *mockUserMicro) RevokeScope(ctx context.Context, userID uuid.UUID, scope string) error {
	args := m.Called(ctx, userID, scope)
	return args.Error(0)
}

type mockTokenMicro struct {
	mock.Mock
}

func (m *mockTokenMicro) IssueTokenPair(ctx context.Context, userID uuid.UUID, scope string) (models.IssuedToken, models.IssuedToken, error) {
	args := m.Called(ctx, userID, scope)
	return args.Get(0).(models.IssuedToken), args.Get(1).(models.IssuedToken), args.Error(2)
}

func (m *mockTokenMicro) Introspect(ctx context.Context, rawToken string) (models.Token, bool, error) {
	args := m.Called(ctx, rawToken)
	return args.Get(0).(models.Token), args.Bool(1), args.Error(2)
}

func (m *mockTokenMicro) Refresh(ctx context.Context, rawRefreshToken string) (models.IssuedToken, models.IssuedToken, error) {
	args := m.Called(ctx, rawRefreshToken)
	return args.Get(0).(models.IssuedToken), args.Get(1).(models.IssuedToken), args.Error(2)
}

func (m *mockTokenMicro) Revoke(ctx context.Context, rawToken string) error {
	args := m.Called(ctx, rawToken)
	return args.Error(0)
}

func Test_authEndpoints_Login(t *testing.T) {
	type args struct {
		ctx     context.Context
		reqData LoginRequestData
	}

	testUserID := uuid.New()
	testUser := models.User{ID: testUserID, Email: "testy@example.com"}
	now := time.Now().UTC()
	accessToken := models.IssuedToken{Value: "access-raw", Scope: "default user:admin", Kind: models.TokenKindAccess, IssuedAt: now, ExpiresAt: now.Add(time.Minute)}
	refreshToken := models.IssuedToken{Value: "refresh-raw", Scope: "default user:admin", Kind: models.TokenKindRefresh, IssuedAt: now, ExpiresAt: now.Add(time.Hour)}

	tests := []struct {
		name      string
		setupMock func(*mockUserMicro, *mockTokenMicro)
		args      args
		want      LoginResponseData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			setupMock: func(um *mockUserMicro, tm *mockTokenMicro) {
				um.On("Authenticate", mock.Anything, "testy@example.com", "correct-password").Return(testUser, error(nil))
				um.On("Scopes", mock.Anything, testUserID).Return("default user:admin", error(nil))
				tm.On("IssueTokenPair", mock.Anything, testUserID, "default user:admin").Return(accessToken, refreshToken, error(nil))
			},
			args: args{
				ctx:     context.Background(),
				reqData: LoginRequestData{Email: "testy@example.com", Password: "correct-password"},
			},
			want: LoginResponseData{
				AccessToken:  issuedTokenToData(accessToken),
				RefreshToken: issuedTokenToData(refreshToken),
			},
			assertion: assert.NoError,
		},
		{
			name: "AuthenticationFails",
			setupMock: func(um *mockUserMicro, tm *mockTokenMicro) {
				um.On("Authenticate", mock.Anything, "bad@example.com", "wrong").Return(models.User{}, assert.AnError)
			},
			args: args{
				ctx:     context.Background(),
				reqData: LoginRequestData{Email: "bad@example.com", Password: "wrong"},
			},
			assertion: assert.Error,
		},
		{
			name: "ScopeResolutionFails",
			setupMock: func(um *mockUserMicro, tm *mockTokenMicro) {
				um.On("Authenticate", mock.Anything, "testy@example.com", "correct-password").Return(testUser, error(nil))
				um.On("Scopes", mock.Anything, testUserID).Return("", assert.AnError)
			},
			args: args{
				ctx:     context.Background(),
				reqData: LoginRequestData{Email: "testy@example.com", Password: "correct-password"},
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testUserMicro := &mockUserMicro{}
			testTokenMicro := &mockTokenMicro{}
			tt.setupMock(testUserMicro, testTokenMicro)

			ae := authEndpoints{userMicro: testUserMicro, tokenMicro: testTokenMicro}
			got, err := ae.Login(tt.args.ctx, tt.args.reqData)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
			testUserMicro.AssertExpectations(t)
			testTokenMicro.AssertExpectations(t)
		})
	}
}

func Test_authEndpoints_GrantScope(t *testing.T) {
	type args struct {
		ctx     context.Context
		reqData GrantScopeRequestData
	}

	testUserID := uuid.New()

	tests := []struct {
		name      string
		setupMock func(*mockUserMicro)
		args      args
		want      ScopeResponseData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			setupMock: func(um *mockUserMicro) {
				um.On("GrantScope", mock.Anything, testUserID, "user:admin").Return(error(nil))
			},
			args: args{
				ctx:     context.Background(),
				reqData: GrantScopeRequestData{UserID: testUserID.String(), Scope: "user:admin"},
			},
			want:      ScopeResponseData{UserID: testUserID.String(), Scope: "user:admin"},
			assertion: assert.NoError,
		},
		{
			name:      "InvalidUserID",
			setupMock: func(um *mockUserMicro) {},
			args: args{
				ctx:     context.Background(),
				reqData: GrantScopeRequestData{UserID: "not-a-uuid", Scope: "user:admin"},
			},
			assertion: assert.Error,
		},
		{
			name: "StoreError",
			setupMock: func(um *mockUserMicro) {
				um.On("GrantScope", mock.Anything, testUserID, "user:admin").Return(assert.AnError)
			},
			args: args{
				ctx:     context.Background(),
				reqData: GrantScopeRequestData{UserID: testUserID.String(), Scope: "user:admin"},
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testUserMicro := &mockUserMicro{}
			tt.setupMock(testUserMicro)

			ae := authEndpoints{userMicro: testUserMicro}
			got, err := ae.GrantScope(tt.args.ctx, tt.args.reqData)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
			testUserMicro.AssertExpectations(t)
		})
	}
}

func Test_authEndpoints_RevokeScope(t *testing.T) {
	type args struct {
		ctx     context.Context
		reqData RevokeScopeRequestData
	}

	testUserID := uuid.New()

	tests := []struct {
		name      string
		setupMock func(*mockUserMicro)
		args      args
		want      ScopeResponseData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			setupMock: func(um *mockUserMicro) {
				um.On("RevokeScope", mock.Anything, testUserID, "user:admin").Return(error(nil))
			},
			args: args{
				ctx:     context.Background(),
				reqData: RevokeScopeRequestData{UserID: testUserID.String(), Scope: "user:admin"},
			},
			want:      ScopeResponseData{UserID: testUserID.String(), Scope: "user:admin"},
			assertion: assert.NoError,
		},
		{
			name:      "InvalidUserID",
			setupMock: func(um *mockUserMicro) {},
			args: args{
				ctx:     context.Background(),
				reqData: RevokeScopeRequestData{UserID: "not-a-uuid", Scope: "user:admin"},
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testUserMicro := &mockUserMicro{}
			tt.setupMock(testUserMicro)

			ae := authEndpoints{userMicro: testUserMicro}
			got, err := ae.RevokeScope(tt.args.ctx, tt.args.reqData)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
			testUserMicro.AssertExpectations(t)
		})
	}
}
