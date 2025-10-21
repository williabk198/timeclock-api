package admin

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/williabk198/timeclock/internal/datastores"
	"github.com/williabk198/timeclock/internal/models"
)

type ContactMicro interface {
	AddPersonAddress(ctx context.Context, address models.ContactAddress) (uuid.UUID, error)
	AddPersonEmail(ctx context.Context, email models.ContactEmail) (uuid.UUID, error)
	AddPersonPhone(ctx context.Context, phone models.ContactPhone) (uuid.UUID, error)
	DeletePerosnAddress(ctx context.Context, personID, addressID uuid.UUID) (models.ContactAddress, error)
	DeletePersonEmail(ctx context.Context, personID, emailID uuid.UUID) (models.ContactEmail, error)
	DeletePersonPhone(ctx context.Context, personID, phoneID uuid.UUID) (models.ContactPhone, error)
	GetAllForPerson(ctx context.Context, personID uuid.UUID) (models.Contacts, error)
	GetPersonAddresses(ctx context.Context, personID uuid.UUID) ([]models.ContactAddress, error)
	GetPersonEmails(ctx context.Context, personID uuid.UUID) ([]models.ContactEmail, error)
	GetPersonPhones(ctx context.Context, personID uuid.UUID) ([]models.ContactPhone, error)
	UpdatePersonAddress(ctx context.Context, personID, addressID uuid.UUID, newVal models.ContactAddress) error
	UpdatePersonEmail(ctx context.Context, personID, emailID uuid.UUID, newVal models.ContactEmail) error
	UpdatePersonPhone(ctx context.Context, personID, phoneID uuid.UUID, newVal models.ContactPhone) error
}

type contactMicroImpl struct {
	contactStore datastores.ContactDatastore
}

func (cmi contactMicroImpl) AddPersonAddress(ctx context.Context, address models.ContactAddress) (uuid.UUID, error) {
	return cmi.contactStore.AddPersonAddress(ctx, address)
}

func (cmi contactMicroImpl) AddPersonEmail(ctx context.Context, email models.ContactEmail) (uuid.UUID, error) {
	return cmi.contactStore.AddPersonEmail(ctx, email)
}

func (cmi contactMicroImpl) AddPersonPhone(ctx context.Context, phone models.ContactPhone) (uuid.UUID, error) {
	return cmi.contactStore.AddPersonPhone(ctx, phone)
}

// DeletePerosnAddress implements ContactMicro.
func (cmi contactMicroImpl) DeletePerosnAddress(ctx context.Context, personID uuid.UUID, addressID uuid.UUID) (models.ContactAddress, error) {
	panic("unimplemented")
}

// DeletePersonEmail implements ContactMicro.
func (cmi contactMicroImpl) DeletePersonEmail(ctx context.Context, personID uuid.UUID, emailID uuid.UUID) (models.ContactEmail, error) {
	panic("unimplemented")
}

// DeletePersonPhone implements ContactMicro.
func (cmi contactMicroImpl) DeletePersonPhone(ctx context.Context, personID uuid.UUID, phoneID uuid.UUID) (models.ContactPhone, error) {
	panic("unimplemented")
}

func (cmi contactMicroImpl) GetAllForPerson(ctx context.Context, personID uuid.UUID) (models.Contacts, error) {
	wg := &sync.WaitGroup{}
	wg.Add(3)

	var addresses []models.ContactAddress
	var addrErr error
	go func() {
		defer wg.Done()
		addresses, addrErr = cmi.contactStore.GetPersonAddresses(ctx, personID)
	}()

	var emails []models.ContactEmail
	var emailErr error
	go func() {
		defer wg.Done()
		emails, emailErr = cmi.contactStore.GetPersonEmails(ctx, personID)
	}()

	var phones []models.ContactPhone
	var phoneErr error
	go func() {
		defer wg.Done()
		phones, phoneErr = cmi.contactStore.GetPersonPhones(ctx, personID)
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

func (cmi contactMicroImpl) GetPersonAddresses(ctx context.Context, personID uuid.UUID) ([]models.ContactAddress, error) {
	return cmi.contactStore.GetPersonAddresses(ctx, personID)
}

func (cmi contactMicroImpl) GetPersonEmails(ctx context.Context, personID uuid.UUID) ([]models.ContactEmail, error) {
	return cmi.contactStore.GetPersonEmails(ctx, personID)
}

func (cmi contactMicroImpl) GetPersonPhones(ctx context.Context, personID uuid.UUID) ([]models.ContactPhone, error) {
	return cmi.contactStore.GetPersonPhones(ctx, personID)
}

func (cmi contactMicroImpl) UpdatePersonAddress(ctx context.Context, personID uuid.UUID, addressID uuid.UUID, newVal models.ContactAddress) error {
	return cmi.contactStore.UpdatePersonAddress(ctx, personID, addressID, newVal)
}

func (cmi contactMicroImpl) UpdatePersonEmail(ctx context.Context, personID uuid.UUID, emailID uuid.UUID, newVal models.ContactEmail) error {
	return cmi.contactStore.UpdatePersonEmail(ctx, personID, emailID, newVal)
}

func (cmi contactMicroImpl) UpdatePersonPhone(ctx context.Context, personID uuid.UUID, phoneID uuid.UUID, newVal models.ContactPhone) error {
	return cmi.contactStore.UpdatePersonPhone(ctx, personID, phoneID, newVal)
}
