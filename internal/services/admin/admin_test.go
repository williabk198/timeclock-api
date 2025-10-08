package admin

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/williabk198/timeclock/internal/datastores"
	"github.com/williabk198/timeclock/internal/models"
)

func TestNewService(t *testing.T) {
	type args struct {
		adminStore   datastores.PersonStore
		contactStore datastores.ContactDatastore
	}

	testPersonStore := &mockPersonStore{}

	tests := []struct {
		name string
		args args
		want Service
	}{
		{
			name: "Success",
			args: args{
				adminStore: testPersonStore,
			},
			want: adminService{
				personStore: testPersonStore,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewService(tt.args.adminStore, tt.args.contactStore))
		})
	}
}

type mockPersonStore struct {
	mock.Mock
}

func (mps *mockPersonStore) Add(ctx context.Context, item models.Person) (id uuid.UUID, err error) {
	args := mps.Called(ctx, item)
	return args.Get(0).(uuid.UUID), args.Error(1)
}
func (mps *mockPersonStore) Delete(ctx context.Context, id uuid.UUID) (item models.Person, err error) {
	args := mps.Called(ctx, id)
	return args.Get(0).(models.Person), args.Error(1)
}
func (mps *mockPersonStore) GetAll(ctx context.Context) (items []models.Person, err error) {
	args := mps.Called(ctx)
	return args.Get(0).([]models.Person), args.Error(1)
}
func (mps *mockPersonStore) GetAllPaginated(ctx context.Context, offset uint, limit uint) (items []models.Person, err error) {
	args := mps.Called(ctx, offset, limit)
	return args.Get(0).([]models.Person), args.Error(1)
}
func (mps *mockPersonStore) GetSpecific(ctx context.Context, id uuid.UUID) (item models.Person, err error) {
	args := mps.Called(ctx, id)
	return args.Get(0).(models.Person), args.Error(1)
}
func (mps *mockPersonStore) Update(ctx context.Context, id uuid.UUID, item models.Person) (err error) {
	args := mps.Called(ctx, id, item)
	return args.Error(0)
}

type mockContactStore struct {
	mock.Mock
}

func (mcs *mockContactStore) AddPersonAddress(ctx context.Context, address models.ContactAddress) (uuid.UUID, error) {
	args := mcs.Called(ctx, address)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (mcs *mockContactStore) AddPersonEmail(ctx context.Context, email models.ContactEmail) (uuid.UUID, error) {
	args := mcs.Called(ctx, email)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (mcs *mockContactStore) AddPersonPhone(ctx context.Context, phone models.ContactPhone) (uuid.UUID, error) {
	args := mcs.Called(ctx, phone)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (mcs *mockContactStore) GetPersonAddresses(ctx context.Context, id uuid.UUID) ([]models.ContactAddress, error) {
	args := mcs.Called(ctx, id)
	return args.Get(0).([]models.ContactAddress), args.Error(1)
}

func (mcs *mockContactStore) GetPersonEmails(ctx context.Context, id uuid.UUID) ([]models.ContactEmail, error) {
	args := mcs.Called(ctx, id)
	return args.Get(0).([]models.ContactEmail), args.Error(1)
}

func (mcs *mockContactStore) GetPersonPhones(ctx context.Context, id uuid.UUID) ([]models.ContactPhone, error) {
	args := mcs.Called(ctx, id)
	return args.Get(0).([]models.ContactPhone), args.Error(1)
}
