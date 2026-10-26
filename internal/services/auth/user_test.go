package auth

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/williabk198/timeclock/internal/models"
)

func Test_userMicroImpl_Authenticate(t *testing.T) {
	type args struct {
		ctx      context.Context
		email    string
		password string
	}

	testUserID := uuid.New()
	hashedPass, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash test password: %v", err)
	}
	testUser := models.User{
		ID:         testUserID,
		Email:      "testy@example.com",
		HashedPass: string(hashedPass),
	}
	enabledMetadata := models.UserMetadata{UserID: testUserID, Status: models.UserStatusEnabled}
	lockedMetadata := models.UserMetadata{UserID: testUserID, Status: models.UserStatusLocked}
	disabledMetadata := models.UserMetadata{UserID: testUserID, Status: models.UserStatusDisabled}
	almostLockedMetadata := models.UserMetadata{UserID: testUserID, Status: models.UserStatusEnabled, FailedLoginCount: maxFailedLoginAttempts - 1}

	tests := []struct {
		name      string
		setupMock func(*mockUserStore)
		args      args
		want      models.User
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			setupMock: func(mus *mockUserStore) {
				mus.On("GetByEmail", mock.Anything, "testy@example.com").Return(testUser, error(nil))
				mus.On("GetMetadata", mock.Anything, testUserID).Return(enabledMetadata, error(nil))
				mus.On("RecordSuccessfulLogin", mock.Anything, testUserID, mock.AnythingOfType("time.Time")).Return(error(nil))
			},
			args: args{
				ctx:      context.Background(),
				email:    "testy@example.com",
				password: "correct-password",
			},
			want:      testUser,
			assertion: assert.NoError,
		},
		{
			name: "UnknownEmail",
			setupMock: func(mus *mockUserStore) {
				mus.On("GetByEmail", mock.Anything, "unknown@example.com").Return(models.User{}, sql.ErrNoRows)
			},
			args: args{
				ctx:      context.Background(),
				email:    "unknown@example.com",
				password: "whatever",
			},
			assertion: assert.Error,
		},
		{
			name: "WrongPassword",
			setupMock: func(mus *mockUserStore) {
				mus.On("GetByEmail", mock.Anything, "wrongpass@example.com").Return(testUser, error(nil))
				mus.On("GetMetadata", mock.Anything, testUserID).Return(enabledMetadata, error(nil))
				mus.On("UpdateMetadata", mock.Anything, testUserID, mock.AnythingOfType("models.UserMetadata")).Return(error(nil))
			},
			args: args{
				ctx:      context.Background(),
				email:    "wrongpass@example.com",
				password: "incorrect-password",
			},
			assertion: assert.Error,
		},
		{
			name: "LockedAccount",
			setupMock: func(mus *mockUserStore) {
				mus.On("GetByEmail", mock.Anything, "locked@example.com").Return(testUser, error(nil))
				mus.On("GetMetadata", mock.Anything, testUserID).Return(lockedMetadata, error(nil))
			},
			args: args{
				ctx:      context.Background(),
				email:    "locked@example.com",
				password: "correct-password",
			},
			assertion: assert.Error,
		},
		{
			name: "DisabledAccount",
			setupMock: func(mus *mockUserStore) {
				mus.On("GetByEmail", mock.Anything, "disabled@example.com").Return(testUser, error(nil))
				mus.On("GetMetadata", mock.Anything, testUserID).Return(disabledMetadata, error(nil))
			},
			args: args{
				ctx:      context.Background(),
				email:    "disabled@example.com",
				password: "correct-password",
			},
			assertion: assert.Error,
		},
		{
			name: "LocksAccountAfterTooManyFailures",
			setupMock: func(mus *mockUserStore) {
				mus.On("GetByEmail", mock.Anything, "almostlocked@example.com").Return(testUser, error(nil))
				mus.On("GetMetadata", mock.Anything, testUserID).Return(almostLockedMetadata, error(nil))
				mus.On("UpdateMetadata", mock.Anything, testUserID, models.UserMetadata{
					UserID:           testUserID,
					Status:           models.UserStatusLocked,
					FailedLoginCount: maxFailedLoginAttempts,
				}).Return(error(nil))
			},
			args: args{
				ctx:      context.Background(),
				email:    "almostlocked@example.com",
				password: "incorrect-password",
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testUserStore := &mockUserStore{}
			tt.setupMock(testUserStore)

			umi := userMicroImpl{userStore: testUserStore}
			got, err := umi.Authenticate(tt.args.ctx, tt.args.email, tt.args.password)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
			testUserStore.AssertExpectations(t)
		})
	}
}

func Test_userMicroImpl_Scopes(t *testing.T) {
	testUserID := uuid.New()

	tests := []struct {
		name      string
		setupMock func(*mockUserStore)
		want      string
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "NoGrantedScopes",
			setupMock: func(mus *mockUserStore) {
				mus.On("GetGrantedScopes", mock.Anything, testUserID).Return([]string{}, error(nil))
			},
			want:      "default",
			assertion: assert.NoError,
		},
		{
			name: "WithGrantedScopes",
			setupMock: func(mus *mockUserStore) {
				mus.On("GetGrantedScopes", mock.Anything, testUserID).Return([]string{"user:admin"}, error(nil))
			},
			want:      "default user:admin",
			assertion: assert.NoError,
		},
		{
			name: "DatastoreError",
			setupMock: func(mus *mockUserStore) {
				mus.On("GetGrantedScopes", mock.Anything, testUserID).Return([]string(nil), assert.AnError)
			},
			want:      "",
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testUserStore := &mockUserStore{}
			tt.setupMock(testUserStore)

			umi := userMicroImpl{userStore: testUserStore}
			got, err := umi.Scopes(context.Background(), testUserID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
			testUserStore.AssertExpectations(t)
		})
	}
}

func Test_userMicroImpl_GrantScope(t *testing.T) {
	testUserID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		testUserStore := &mockUserStore{}
		testUserStore.On("GrantScope", mock.Anything, testUserID, "user:admin").Return(error(nil))

		umi := userMicroImpl{userStore: testUserStore}
		err := umi.GrantScope(context.Background(), testUserID, "user:admin")

		assert.NoError(t, err)
		testUserStore.AssertExpectations(t)
	})

	t.Run("EmptyScope", func(t *testing.T) {
		umi := userMicroImpl{userStore: &mockUserStore{}}
		err := umi.GrantScope(context.Background(), testUserID, "")

		assert.ErrorIs(t, err, ErrInvalidRequest)
	})
}

func Test_userMicroImpl_RevokeScope(t *testing.T) {
	testUserID := uuid.New()

	testUserStore := &mockUserStore{}
	testUserStore.On("RevokeScope", mock.Anything, testUserID, "user:admin").Return(error(nil))

	umi := userMicroImpl{userStore: testUserStore}
	err := umi.RevokeScope(context.Background(), testUserID, "user:admin")

	assert.NoError(t, err)
	testUserStore.AssertExpectations(t)
}

func Test_userMicroImpl_Register(t *testing.T) {
	type args struct {
		ctx      context.Context
		personID uuid.UUID
		email    string
		password string
	}

	testPersonID := uuid.New()
	testUserID := uuid.New()

	tests := []struct {
		name      string
		setupMock func(*mockUserStore)
		args      args
		want      uuid.UUID
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			setupMock: func(mus *mockUserStore) {
				mus.On("GetByEmail", mock.Anything, "new@example.com").Return(models.User{}, sql.ErrNoRows)
				mus.On("Add", mock.Anything, mock.AnythingOfType("models.User"), mock.AnythingOfType("models.UserMetadata")).Return(testUserID, error(nil))
			},
			args: args{
				ctx:      context.Background(),
				personID: testPersonID,
				email:    "new@example.com",
				password: "a-good-password",
			},
			want:      testUserID,
			assertion: assert.NoError,
		},
		{
			name: "EmailInUse",
			setupMock: func(mus *mockUserStore) {
				mus.On("GetByEmail", mock.Anything, "taken@example.com").Return(models.User{ID: uuid.New()}, error(nil))
			},
			args: args{
				ctx:      context.Background(),
				personID: testPersonID,
				email:    "taken@example.com",
				password: "a-good-password",
			},
			want:      uuid.Nil,
			assertion: assert.Error,
		},
		{
			name:      "MissingPassword",
			setupMock: func(mus *mockUserStore) {},
			args: args{
				ctx:      context.Background(),
				personID: testPersonID,
				email:    "missing@example.com",
				password: "",
			},
			want:      uuid.Nil,
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testUserStore := &mockUserStore{}
			tt.setupMock(testUserStore)

			umi := userMicroImpl{userStore: testUserStore}
			got, err := umi.Register(tt.args.ctx, tt.args.personID, tt.args.email, tt.args.password)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
			testUserStore.AssertExpectations(t)
		})
	}
}
