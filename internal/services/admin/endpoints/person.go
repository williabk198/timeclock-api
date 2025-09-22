package endpoints

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/williabk198/timeclock/internal/models"
	"github.com/williabk198/timeclock/internal/services/admin"
	"github.com/williabk198/timeclock/internal/utils"
)

type adminPersonEndpoints struct {
	adminService admin.Service
}

// Add implements PersonEndpoints.
func (ape adminPersonEndpoints) Add(ctx context.Context, person PersonData) (PersonData, error) {
	pronouns, err := utils.ParsePronouns(person.Pronouns)
	if err != nil {
		return PersonData{}, fmt.Errorf("failed to process pronoun data: %w", err)
	}

	dbPerson := models.Person{
		Name:        person.Name,
		DateOfBirth: time.Unix(person.DateOfBirth, 0),
		Gender:      models.Gender(person.Gender),
		Pronouns:    pronouns,
	}

	id, err := ape.adminService.AddPerson(ctx, dbPerson)
	if err != nil {
		return PersonData{}, fmt.Errorf("failed to add person to database: %w", err)
	}

	person.ID = id.String()
	return person, nil
}

// Delete implements PersonEndpoints.
func (ape adminPersonEndpoints) Delete(ctx context.Context, idStr string) (PersonData, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return PersonData{}, err
	}

	person, err := ape.adminService.DeletePerson(ctx, id)
	if err != nil {
		return PersonData{}, err
	}

	return PersonData{
		ID:          id.String(),
		Name:        person.Name,
		DateOfBirth: person.DateOfBirth.Unix(),
		Gender:      string(person.Gender),
		Pronouns:    person.Pronouns.String(),
	}, nil
}

// GetAll implements PersonEndpoints.
func (ape adminPersonEndpoints) GetAll(ctx context.Context, reqData GetPaginatedRequestData) ([]PersonData, error) {
	persons, err := ape.adminService.GetAllPersons(ctx, reqData.Offset, reqData.Limit)
	if err != nil {
		return nil, err
	}

	result := make([]PersonData, len(persons))
	for i, p := range persons {
		result[i] = PersonData{
			ID:          p.ID.String(),
			Name:        p.Name,
			DateOfBirth: p.DateOfBirth.Unix(),
			Gender:      string(p.Gender),
			Pronouns:    p.Pronouns.String(),
		}
	}

	return result, nil
}

// GetSpecific implements PersonEndpoints.
func (ape adminPersonEndpoints) GetSpecific(ctx context.Context, idStr string) (PersonData, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return PersonData{}, err
	}

	person, err := ape.adminService.GetPerson(ctx, id)
	if err != nil {
		return PersonData{}, err
	}

	return PersonData{
		ID:          id.String(),
		Name:        person.Name,
		DateOfBirth: person.DateOfBirth.Unix(),
		Gender:      string(person.Gender),
		Pronouns:    person.Pronouns.String(),
	}, nil
}

// GetSpecificContact implements PersonEndpoints.
func (ape adminPersonEndpoints) GetSpecificContacts(ctx context.Context, idStr string) (PersonContactData, error) {
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

// GetSpecificContactAddresses implements PersonEndpoints.
func (ape adminPersonEndpoints) GetSpecificContactAddresses(ctx context.Context, idStr string) ([]PersonAddressData, error) {
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

// GetSpecificContactEmails implements PersonEndpoints.
func (ape adminPersonEndpoints) GetSpecificContactEmails(ctx context.Context, idStr string) ([]PersonEmailData, error) {
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

// GetSpecificContactPhones implements PersonEndpoints.
func (ape adminPersonEndpoints) GetSpecificContactPhones(ctx context.Context, idStr string) ([]PersonPhoneData, error) {
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

// Update implements PersonEndpoints.
func (ape adminPersonEndpoints) Update(ctx context.Context, updateReqData UpdateRequestData[PersonData]) (PersonData, error) {
	id, err := uuid.Parse(updateReqData.ID)
	if err != nil {
		return PersonData{}, err
	}

	pronouns, err := utils.ParsePronouns(updateReqData.Data.Pronouns)
	if err != nil {
		return PersonData{}, err
	}

	updatedVals := models.Person{
		Name:        updateReqData.Data.Name,
		DateOfBirth: time.Unix(updateReqData.Data.DateOfBirth, 0),
		Gender:      models.Gender(updateReqData.Data.Gender),
		Pronouns:    pronouns,
	}

	err = ape.adminService.UpdatePerson(ctx, id, updatedVals)
	if err != nil {
		return PersonData{}, err
	}

	return updateReqData.Data, nil
}

func (ape adminPersonEndpoints) convertContactAddressSliceToPersonAddressDataSlice(addresses []models.ContactAddress) []PersonAddressData {
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

func (ape adminPersonEndpoints) convertContactEmailSliceToPersonEmailDataSlice(emails []models.ContactEmail) []PersonEmailData {
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

func (ape adminPersonEndpoints) convertContactPhoneSliceToPersonPhoneDataSlice(phones []models.ContactPhone) []PersonPhoneData {
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

func (ape adminPersonEndpoints) convertContactsModelToPersonContactData(contacts models.Contacts) PersonContactData {
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
