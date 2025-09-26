package endpoints

import (
	"context"

	"github.com/google/uuid"
	"github.com/williabk198/timeclock/internal/models"
	"github.com/williabk198/timeclock/internal/services/admin"
)

type adminContactEndpoints struct {
	adminService admin.Service
}

func (ape adminContactEndpoints) GetPersonContacts(ctx context.Context, idStr string) (PersonContactData, error) {
	// TODO: Change this function to use adminService.GetPersonContactAddress, adminService.GetPersonContactEmails, and adminService.GetPersonContactPhones

	id, err := uuid.Parse(idStr)
	if err != nil {
		return PersonContactData{}, err
	}

	contactData, err := ape.adminService.GetPersonContacts(ctx, id)
	if err != nil {
		return PersonContactData{}, err
	}

	return ape.convertContactsModelToPersonContactData(contactData), nil
}

func (ape adminContactEndpoints) GetPersonContactAddresses(ctx context.Context, idStr string) ([]PersonAddressData, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}

	addresses, err := ape.adminService.GetPersonContactAddresses(ctx, id)
	if err != nil {
		return nil, err
	}

	return ape.convertContactAddressSliceToPersonAddressDataSlice(addresses), nil
}

func (ape adminContactEndpoints) GetPersonContactEmails(ctx context.Context, idStr string) ([]PersonEmailData, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}

	emails, err := ape.adminService.GetPersonContactEmails(ctx, id)
	if err != nil {
		return nil, err
	}

	return ape.convertContactEmailSliceToPersonEmailDataSlice(emails), nil
}

func (ape adminContactEndpoints) GetPersonContactPhones(ctx context.Context, idStr string) ([]PersonPhoneData, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}

	phones, err := ape.adminService.GetPersonContactPhones(ctx, id)
	if err != nil {
		return nil, err
	}

	return ape.convertContactPhoneSliceToPersonPhoneDataSlice(phones), nil
}

func (ape adminContactEndpoints) convertContactAddressSliceToPersonAddressDataSlice(addresses []models.ContactAddress) []PersonAddressData {
	result := make([]PersonAddressData, len(addresses))
	for i, a := range addresses {
		result[i] = PersonAddressData{
			ID:         a.ID.String(),
			Street1:    a.Street1,
			Street2:    a.Street2,
			Locality:   a.Locality,
			Region:     a.Region,
			PostalCode: a.PostalCode,
			Country:    a.Country,
			Type:       a.Type,
			Primary:    a.Primary,
		}
	}
	return result
}

func (ape adminContactEndpoints) convertContactEmailSliceToPersonEmailDataSlice(emails []models.ContactEmail) []PersonEmailData {
	result := make([]PersonEmailData, len(emails))
	for i, e := range emails {
		result[i] = PersonEmailData{
			ID:      e.ID.String(),
			Email:   e.String(),
			Primary: e.Primary,
		}
	}
	return result
}

func (ape adminContactEndpoints) convertContactPhoneSliceToPersonPhoneDataSlice(phones []models.ContactPhone) []PersonPhoneData {
	result := make([]PersonPhoneData, len(phones))
	for i, p := range phones {
		result[i] = PersonPhoneData{
			ID:          p.ID.String(),
			CountryCode: p.CountryCode,
			PhoneNumber: p.PhoneNumber,
			Type:        p.Type,
			Primary:     p.Primary,
		}
	}
	return result
}

func (ape adminContactEndpoints) convertContactsModelToPersonContactData(contacts models.Contacts) PersonContactData {
	// TODO: Remove this function

	adressess := make([]PersonAddressData, len(contacts.Addresses))
	for i, a := range contacts.Addresses {
		adressess[i] = PersonAddressData{
			ID:         a.ID.String(),
			Street1:    a.Street1,
			Street2:    a.Street2,
			Locality:   a.Locality,
			Region:     a.Region,
			PostalCode: a.PostalCode,
			Country:    a.Country,
			Type:       a.Type,
			Primary:    a.Primary,
		}
	}

	emails := make([]PersonEmailData, len(contacts.Email))
	for i, e := range contacts.Email {
		emails[i] = PersonEmailData{
			ID:      e.ID.String(),
			Email:   e.String(),
			Primary: e.Primary,
		}
	}

	phones := make([]PersonPhoneData, len(contacts.Phone))
	for i, p := range contacts.Phone {
		phones[i] = PersonPhoneData{
			ID:          p.ID.String(),
			CountryCode: p.CountryCode,
			PhoneNumber: p.PhoneNumber,
			Type:        p.Type,
			Primary:     p.Primary,
		}
	}

	return PersonContactData{
		Addresses:    adressess,
		Emails:       emails,
		PhoneNumbers: phones,
	}
}
