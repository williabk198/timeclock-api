package admin

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/williabk198/timeclock/internal/datastores"
	"github.com/williabk198/timeclock/internal/models"
)

type Service interface {
	AddPerson(context.Context, models.Person) (uuid.UUID, error)
	AddPersonContactAddress(context.Context, models.ContactAddress) (uuid.UUID, error)
	AddPersonContactEmail(context.Context, models.ContactEmail) (uuid.UUID, error)
	AddPersonContactPhone(context.Context, models.ContactPhone) (uuid.UUID, error)
	DeletePerson(context.Context, uuid.UUID) (models.Person, error)
	GetAllPersons(ctx context.Context, offset uint, limit uint) ([]models.Person, error)
	GetPersonContacts(context.Context, uuid.UUID) (models.Contacts, error)
	GetPersonContactAddresses(context.Context, uuid.UUID) ([]models.ContactAddress, error)
	GetPersonContactEmails(context.Context, uuid.UUID) ([]models.ContactEmail, error)
	GetPersonContactPhones(context.Context, uuid.UUID) ([]models.ContactPhone, error)
	GetPerson(context.Context, uuid.UUID) (models.Person, error)
	UpdatePerson(context.Context, uuid.UUID, models.Person) error
}

type adminService struct {
	personStore datastores.PersonStore
}

// AddPersonContactAddress implements Service.
func (as adminService) AddPersonContactAddress(ctx context.Context, address models.ContactAddress) (uuid.UUID, error) {
	return as.personStore.AddSpecificContactAddress(ctx, address)
}

// AddPersonContactEmail implements Service.
func (as adminService) AddPersonContactEmail(ctx context.Context, email models.ContactEmail) (uuid.UUID, error) {
	return as.personStore.AddSpecificContactEmail(ctx, email)
}

// AddPersonContactPhone implements Service.
func (as adminService) AddPersonContactPhone(ctx context.Context, phone models.ContactPhone) (uuid.UUID, error) {
	return as.personStore.AddSpecificContactPhone(ctx, phone)
}

// AddPunch implements Service.
func (as adminService) AddPerson(ctx context.Context, person models.Person) (uuid.UUID, error) {
	return as.personStore.Add(ctx, person)
}

// DeletePerson implements Service.
func (as adminService) DeletePerson(ctx context.Context, id uuid.UUID) (models.Person, error) {
	return as.personStore.Delete(ctx, id)
}

// GetAllPersons implements Service.
func (as adminService) GetAllPersons(ctx context.Context, offset uint, limit uint) ([]models.Person, error) {
	return as.personStore.GetAllPaginated(ctx, offset, limit)
}

// GetPerson implements Service.
func (as adminService) GetPerson(ctx context.Context, id uuid.UUID) (models.Person, error) {
	return as.personStore.GetSpecific(ctx, id)
}

// GetPersonContacts implements Service.
func (as adminService) GetPersonContacts(ctx context.Context, uuid uuid.UUID) (models.Contacts, error) {
	wg := &sync.WaitGroup{}
	wg.Add(3)

	var addresses []models.ContactAddress
	var addrErr error
	go func() {
		defer wg.Done()
		addresses, addrErr = as.personStore.GetSpecificContactAddresses(ctx, uuid)
	}()

	var emails []models.ContactEmail
	var emailErr error
	go func() {
		defer wg.Done()
		emails, emailErr = as.personStore.GetSpecificContactEmails(ctx, uuid)
	}()

	var phones []models.ContactPhone
	var phoneErr error
	go func() {
		defer wg.Done()
		phones, phoneErr = as.personStore.GetSpecificContactPhones(ctx, uuid)
	}()

	wg.Wait()
	switch {
	case addrErr != nil:
		return models.Contacts{}, addrErr
	case emailErr != nil:
		return models.Contacts{}, emailErr
	case phoneErr != nil:
		return models.Contacts{}, phoneErr
	}

	return models.Contacts{
		Addresses: addresses,
		Email:     emails,
		Phone:     phones,
	}, nil
}

// GetPersonContactAddresses implements Service.
func (as adminService) GetPersonContactAddresses(ctx context.Context, id uuid.UUID) ([]models.ContactAddress, error) {
	return as.personStore.GetSpecificContactAddresses(ctx, id)
}

// GetPersonContactEmails implements Service.
func (as adminService) GetPersonContactEmails(ctx context.Context, id uuid.UUID) ([]models.ContactEmail, error) {
	return as.personStore.GetSpecificContactEmails(ctx, id)
}

// GetPersonContactPhones implements Service.
func (as adminService) GetPersonContactPhones(ctx context.Context, id uuid.UUID) ([]models.ContactPhone, error) {
	return as.personStore.GetSpecificContactPhones(ctx, id)
}

// UpdatePerson implements Service.
func (as adminService) UpdatePerson(ctx context.Context, id uuid.UUID, data models.Person) error {
	return as.personStore.Update(ctx, id, data)
}

func NewService(personStore datastores.PersonStore) Service {
	return adminService{
		personStore: personStore,
	}
}
