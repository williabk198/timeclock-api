package endpoints

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/williabk198/timeclock/internal/models"
	"github.com/williabk198/timeclock/internal/services/admin"
)

type mockAdminService struct {
	mock.Mock
}

func (mas *mockAdminService) Contact() admin.ContactMicro {
	args := mas.Called()
	return args.Get(0).(admin.ContactMicro)
}

func (mas *mockAdminService) Person() admin.PersonMicro {
	args := mas.Called()
	return args.Get(0).(admin.PersonMicro)
}

type mockContactMicro struct {
	mock.Mock
}

func (mcm *mockContactMicro) AddPersonAddress(ctx context.Context, address models.ContactAddress) (uuid.UUID, error) {
	args := mcm.Called(ctx, address)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (mcm *mockContactMicro) AddPersonEmail(ctx context.Context, email models.ContactEmail) (uuid.UUID, error) {
	args := mcm.Called(ctx, email)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (mcm *mockContactMicro) AddPersonPhone(ctx context.Context, phone models.ContactPhone) (uuid.UUID, error) {
	args := mcm.Called(ctx, phone)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (mcm *mockContactMicro) GetAllForPerson(ctx context.Context, id uuid.UUID) (models.Contacts, error) {
	args := mcm.Called(ctx, id)
	return args.Get(0).(models.Contacts), args.Error(1)
}

// GetPersonContactAddresses implements Service.
func (mcm *mockContactMicro) GetPersonAddresses(ctx context.Context, id uuid.UUID) ([]models.ContactAddress, error) {
	args := mcm.Called(ctx, id)
	return args.Get(0).([]models.ContactAddress), args.Error(1)
}

// GetPersonContactEmails implements Service.
func (mcm *mockContactMicro) GetPersonEmails(ctx context.Context, id uuid.UUID) ([]models.ContactEmail, error) {
	args := mcm.Called(ctx, id)
	return args.Get(0).([]models.ContactEmail), args.Error(1)
}

// GetPersonContactPhones implements Service.
func (mcm *mockContactMicro) GetPersonPhones(ctx context.Context, id uuid.UUID) ([]models.ContactPhone, error) {
	args := mcm.Called(ctx, id)
	return args.Get(0).([]models.ContactPhone), args.Error(1)
}

type mockPersonMicro struct {
	mock.Mock
}

func (mpm *mockPersonMicro) Add(ctx context.Context, person models.Person) (uuid.UUID, error) {
	args := mpm.Called(ctx, person)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (mpm *mockPersonMicro) Delete(ctx context.Context, id uuid.UUID) (models.Person, error) {
	args := mpm.Called(ctx, id)
	return args.Get(0).(models.Person), args.Error(1)
}

func (mpm *mockPersonMicro) GetAll(ctx context.Context, offset, limit uint) ([]models.Person, error) {
	args := mpm.Called(ctx, offset, limit)
	return args.Get(0).([]models.Person), args.Error(1)
}

func (mpm *mockPersonMicro) GetSpecific(ctx context.Context, id uuid.UUID) (models.Person, error) {
	args := mpm.Called(ctx, id)
	return args.Get(0).(models.Person), args.Error(1)
}

func (mpm *mockPersonMicro) Update(ctx context.Context, id uuid.UUID, data models.Person) error {
	args := mpm.Called(ctx, id, data)
	return args.Error(0)
}
