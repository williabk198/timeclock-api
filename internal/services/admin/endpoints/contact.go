package endpoints

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/williabk198/timeclock/internal/models"
	"github.com/williabk198/timeclock/internal/services/admin"
)

type adminContactEndpoints struct {
	adminService admin.Service
}

// AddContactAddressForPerson implements ContactEndpoints.
func (ace adminContactEndpoints) AddContactAddressForPerson(ctx context.Context, reqData AddSubItemRequestData[PersonAddressData]) (PersonEmailData, error) {
	panic("unimplemented")
}

// AddContactEmailForPerson implements ContactEndpoints.
func (ace adminContactEndpoints) AddContactEmailForPerson(ctx context.Context, reqData AddSubItemRequestData[PersonEmailData]) (PersonEmailData, error) {
	personID, err := uuid.Parse(reqData.ParentID)
	if err != nil {
		return PersonEmailData{}, err
	}

	// Oversimplifed email validation. This will eventually need to be more robust.
	splitEmail := strings.Split(reqData.Data.Email, "@")
	if len(splitEmail) != 2 || splitEmail[0] == "" || splitEmail[1] == "" || !strings.Contains(splitEmail[1], ".") ||
		splitEmail[1][0] == '.' || splitEmail[1][len(splitEmail)-1] == '.' {

		return PersonEmailData{}, fmt.Errorf("provided email address was ill formated")
	}

	emailID, err := ace.adminService.AddPersonContactEmail(ctx, models.ContactEmail{
		PersonID: personID,
		Username: splitEmail[0],
		Provider: splitEmail[1],
		Primary:  reqData.Data.Primary,
	})
	if err != nil {
		return PersonEmailData{}, err
	}

	reqData.Data.ID = emailID.String()
	return reqData.Data, nil
}

// AddContactPhoneForPerson implements ContactEndpoints.
func (ace adminContactEndpoints) AddContactPhoneForPerson(ctx context.Context, reqData AddSubItemRequestData[PersonPhoneData]) (PersonPhoneData, error) {
	panic("unimplemented")
}

func (ace adminContactEndpoints) GetPersonContacts(ctx context.Context, idStr string) (PersonContactData, error) {
	// TODO: Change this function to use adminService.GetPersonContactAddress, adminService.GetPersonContactEmails, and adminService.GetPersonContactPhones

	id, err := uuid.Parse(idStr)
	if err != nil {
		return PersonContactData{}, err
	}

	contactData, err := ace.adminService.GetPersonContacts(ctx, id)
	if err != nil {
		return PersonContactData{}, err
	}

	return ace.convertContactsModelToPersonContactData(contactData), nil
}

func (ace adminContactEndpoints) GetPersonContactAddresses(ctx context.Context, idStr string) ([]PersonAddressData, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}

	addresses, err := ace.adminService.GetPersonContactAddresses(ctx, id)
	if err != nil {
		return nil, err
	}

	return ace.convertContactAddressSliceToPersonAddressDataSlice(addresses), nil
}

func (ace adminContactEndpoints) GetPersonContactEmails(ctx context.Context, idStr string) ([]PersonEmailData, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}

	emails, err := ace.adminService.GetPersonContactEmails(ctx, id)
	if err != nil {
		return nil, err
	}

	return ace.convertContactEmailSliceToPersonEmailDataSlice(emails), nil
}

func (ace adminContactEndpoints) GetPersonContactPhones(ctx context.Context, idStr string) ([]PersonPhoneData, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}

	phones, err := ace.adminService.GetPersonContactPhones(ctx, id)
	if err != nil {
		return nil, err
	}

	return ace.convertContactPhoneSliceToPersonPhoneDataSlice(phones), nil
}

func (ace adminContactEndpoints) convertContactAddressSliceToPersonAddressDataSlice(addresses []models.ContactAddress) []PersonAddressData {
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
			Type:       string(a.Type),
			Primary:    a.Primary,
		}
	}
	return result
}

func (ace adminContactEndpoints) convertContactEmailSliceToPersonEmailDataSlice(emails []models.ContactEmail) []PersonEmailData {
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

func (ace adminContactEndpoints) convertContactPhoneSliceToPersonPhoneDataSlice(phones []models.ContactPhone) []PersonPhoneData {
	result := make([]PersonPhoneData, len(phones))
	for i, p := range phones {
		result[i] = PersonPhoneData{
			ID:          p.ID.String(),
			CountryCode: p.CountryCode,
			PhoneNumber: p.PhoneNumber,
			Type:        string(p.Type),
			Primary:     p.Primary,
		}
	}
	return result
}

func (ace adminContactEndpoints) convertContactsModelToPersonContactData(contacts models.Contacts) PersonContactData {
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
			Type:       string(a.Type),
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
			Type:        string(p.Type),
			Primary:     p.Primary,
		}
	}

	return PersonContactData{
		Addresses:    adressess,
		Emails:       emails,
		PhoneNumbers: phones,
	}
}
