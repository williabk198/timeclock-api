package endpoints

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/williabk198/timeclock/internal/models"
)

func Test_adminContactEndpoints_AddContactAddressForPerson(t *testing.T) {
	type args struct {
		ctx     context.Context
		reqData AddSubItemRequestData[PersonAddressData]
	}

	testValidPersonID := uuid.New()
	testNotFoundPersonID := uuid.New()
	testAddressID := uuid.New()

	testValidAddressData := PersonAddressData{
		Street1:    "123 Test Dr",
		Street2:    "",
		Locality:   "Testerville",
		Region:     "Testia",
		PostalCode: "12345-6789",
		Country:    "Testopia",
		Type:       "mailing",
		Primary:    true,
	}
	testValidAddressDB := models.ContactAddress{
		PersonID:   testValidPersonID,
		Street1:    "123 Test Dr",
		Street2:    "",
		Locality:   "Testerville",
		Region:     "Testia",
		PostalCode: "12345-6789",
		Country:    "Testopia",
		Type:       models.AddressTypeMailing,
		Primary:    true,
	}

	testInvalidAddressDB := models.ContactAddress{
		PersonID:   testNotFoundPersonID,
		Street1:    "123 Test Dr",
		Street2:    "",
		Locality:   "Testerville",
		Region:     "Testia",
		PostalCode: "12345-6789",
		Country:    "Testopia",
		Type:       models.AddressTypeMailing,
		Primary:    true,
	}

	testContactMicro := &mockContactMicro{}
	testContactMicro.On("AddPersonAddress", mock.Anything, testValidAddressDB).Return(testAddressID, error(nil))
	testContactMicro.On("AddPersonAddress", mock.Anything, testInvalidAddressDB).Return(uuid.Nil, assert.AnError)

	tests := []struct {
		name      string
		ace       adminContactEndpoints
		args      args
		want      PersonAddressData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: AddSubItemRequestData[PersonAddressData]{
					ParentID: testValidPersonID.String(),
					Data:     testValidAddressData,
				},
			},
			want: PersonAddressData{
				ID:         testAddressID.String(),
				Street1:    testValidAddressData.Street1,
				Street2:    testValidAddressData.Street2,
				Locality:   testValidAddressData.Locality,
				Region:     testValidAddressData.Region,
				PostalCode: testValidAddressData.PostalCode,
				Country:    testValidAddressData.Country,
				Type:       testValidAddressData.Type,
				Primary:    testValidAddressData.Primary,
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; Person DNE",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: AddSubItemRequestData[PersonAddressData]{
					ParentID: testNotFoundPersonID.String(),
					Data:     testValidAddressData,
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Invalid Person ID",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: AddSubItemRequestData[PersonAddressData]{
					ParentID: "ivalid_id",
					Data:     testValidAddressData,
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Invalid Address; Missing Street 1",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: AddSubItemRequestData[PersonAddressData]{
					ParentID: testValidPersonID.String(),
					Data: PersonAddressData{
						Locality:   "Testetville",
						Region:     "Testaria",
						PostalCode: "12345-6789",
						Country:    "Testopia",
						Type:       "mailing",
						Primary:    true,
					},
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Invalid Address; Missing Locality",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: AddSubItemRequestData[PersonAddressData]{
					ParentID: testValidPersonID.String(),
					Data: PersonAddressData{
						Locality:   "Testetville",
						Region:     "Testaria",
						PostalCode: "12345-6789",
						Country:    "Testopia",
						Type:       "mailing",
						Primary:    true,
					},
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Invalid Address; Missing Region",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: AddSubItemRequestData[PersonAddressData]{
					ParentID: testValidPersonID.String(),
					Data: PersonAddressData{
						Street1:    "123 Test Dr",
						Locality:   "Testetville",
						PostalCode: "12345-6789",
						Country:    "Testopia",
						Type:       "mailing",
						Primary:    true,
					},
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Invalid Address; Missing Postal Code",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: AddSubItemRequestData[PersonAddressData]{
					ParentID: testValidPersonID.String(),
					Data: PersonAddressData{
						Street1:  "123 Test Dr",
						Locality: "Testetville",
						Region:   "Testaria",
						Country:  "Testopia",
						Type:     "mailing",
						Primary:  true,
					},
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Invalid Address; Missing Country",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: AddSubItemRequestData[PersonAddressData]{
					ParentID: testValidPersonID.String(),
					Data: PersonAddressData{
						Street1:    "123 Test Dr",
						Locality:   "Testetville",
						Region:     "Testaria",
						PostalCode: "12345-6789",
						Type:       "mailing",
						Primary:    true,
					},
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Invalid Address; Missing Type",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: AddSubItemRequestData[PersonAddressData]{
					ParentID: testValidPersonID.String(),
					Data: PersonAddressData{
						Street1:    "123 Test Dr",
						Locality:   "Testetville",
						Region:     "Testaria",
						PostalCode: "12345-6789",
						Country:    "Testopia",
						Primary:    true,
					},
				},
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ace.AddContactAddressForPerson(tt.args.ctx, tt.args.reqData)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_adminContactEndpoints_AddContactEmailForPerson(t *testing.T) {
	type args struct {
		ctx     context.Context
		reqData AddSubItemRequestData[PersonEmailData]
	}

	testValidPersonID := uuid.New()
	testNotFoundPersonID := uuid.New()
	testEmailID := uuid.New()

	testValidEmailData := PersonEmailData{
		Email:   "test@example.com",
		Primary: true,
	}
	testValidEmailDB := models.ContactEmail{
		PersonID: testValidPersonID,
		Username: "test",
		Provider: "example.com",
		Primary:  true,
	}

	testInvalidEmailDB := models.ContactEmail{
		PersonID: testNotFoundPersonID,
		Username: "test",
		Provider: "example.com",
		Primary:  true,
	}

	testInvalidEmailData1 := PersonEmailData{
		Email: "@example.com",
	}
	testInvalidEmailData2 := PersonEmailData{
		Email: "user@",
	}
	testInvalidEmailData3 := PersonEmailData{
		Email: "user@example",
	}

	testContactMicro := &mockContactMicro{}
	testContactMicro.On("AddPersonEmail", mock.Anything, testValidEmailDB).Return(testEmailID, error(nil))
	testContactMicro.On("AddPersonEmail", mock.Anything, testInvalidEmailDB).Return(uuid.Nil, assert.AnError)

	tests := []struct {
		name      string
		ace       adminContactEndpoints
		args      args
		want      PersonEmailData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: AddSubItemRequestData[PersonEmailData]{
					ParentID: testValidPersonID.String(),
					Data:     testValidEmailData,
				},
			},
			want: PersonEmailData{
				ID:      testEmailID.String(),
				Email:   testValidEmailData.Email,
				Primary: testValidEmailData.Primary,
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; Person DNE",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: AddSubItemRequestData[PersonEmailData]{
					ParentID: testNotFoundPersonID.String(),
					Data:     testValidEmailData,
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Invalid Person ID",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: AddSubItemRequestData[PersonEmailData]{
					ParentID: "ivalid_id",
					Data:     testValidEmailData,
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; No Email Given",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: AddSubItemRequestData[PersonEmailData]{
					ParentID: testValidPersonID.String(),
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Invalid Email; No Username",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: AddSubItemRequestData[PersonEmailData]{
					ParentID: testValidPersonID.String(),
					Data:     testInvalidEmailData1,
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Invalid Email; No Domain",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: AddSubItemRequestData[PersonEmailData]{
					ParentID: testValidPersonID.String(),
					Data:     testInvalidEmailData2,
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Invalid Email; Missing Top-level Domain",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: AddSubItemRequestData[PersonEmailData]{
					ParentID: testValidPersonID.String(),
					Data:     testInvalidEmailData3,
				},
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ace.AddContactEmailForPerson(tt.args.ctx, tt.args.reqData)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_adminContactEndpoints_AddContactPhoneForPerson(t *testing.T) {
	type args struct {
		ctx     context.Context
		reqData AddSubItemRequestData[PersonPhoneData]
	}

	testValidPersonID := uuid.New()
	testNotFoundPersonID := uuid.New()
	testPhoneID := uuid.New()

	testValidPhoneData := PersonPhoneData{
		CountryCode: 1,
		PhoneNumber: "(555)555-5555",
		Type:        "home",
		Primary:     true,
	}
	testValidPhoneDB := models.ContactPhone{
		PersonID:    testValidPersonID,
		CountryCode: 1,
		PhoneNumber: "(555)555-5555",
		Type:        models.PhoneTypeHome,
		Primary:     true,
	}

	testInvalidPhoneDB := models.ContactPhone{
		PersonID:    testNotFoundPersonID,
		CountryCode: 1,
		PhoneNumber: "(555)555-5555",
		Type:        models.PhoneTypeHome,
		Primary:     true,
	}

	testContactMicro := &mockContactMicro{}
	testContactMicro.On("AddPersonPhone", mock.Anything, testValidPhoneDB).Return(testPhoneID, error(nil))
	testContactMicro.On("AddPersonPhone", mock.Anything, testInvalidPhoneDB).Return(uuid.Nil, assert.AnError)

	tests := []struct {
		name      string
		ace       adminContactEndpoints
		args      args
		want      PersonPhoneData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: AddSubItemRequestData[PersonPhoneData]{
					ParentID: testValidPersonID.String(),
					Data:     testValidPhoneData,
				},
			},
			want: PersonPhoneData{
				ID:          testPhoneID.String(),
				CountryCode: testValidPhoneData.CountryCode,
				PhoneNumber: testValidPhoneData.PhoneNumber,
				Type:        string(testValidPhoneData.Type),
				Primary:     testValidPhoneData.Primary,
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; Person DNE",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: AddSubItemRequestData[PersonPhoneData]{
					ParentID: testNotFoundPersonID.String(),
					Data:     testValidPhoneData,
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Invalid Person ID",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: AddSubItemRequestData[PersonPhoneData]{
					ParentID: "invalid_id",
					Data:     testValidPhoneData,
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Invalid Phone; Bad Country Code",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: AddSubItemRequestData[PersonPhoneData]{
					ParentID: testValidPersonID.String(),
					Data: PersonPhoneData{
						CountryCode: -1,
						PhoneNumber: "555-555-5555",
						Type:        "home",
						Primary:     false,
					},
				},
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ace.AddContactPhoneForPerson(tt.args.ctx, tt.args.reqData)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_adminContactEndpoints_DeleteContactAddressForPerson(t *testing.T) {
	type args struct {
		ctx     context.Context
		reqData DeleteContactRequestData
	}

	testPersonID := uuid.New()
	testPersonNotFoundID := uuid.New()

	testAddressID := uuid.New()
	testAddressNotFoundID := uuid.New()

	testAddress := models.ContactAddress{
		ID:         testAddressID,
		PersonID:   testPersonID,
		Street1:    "123 Test Dr",
		Street2:    "",
		Locality:   "Testerville",
		Region:     "Testia",
		PostalCode: "12345-6789",
		Country:    "Testopia",
		Type:       models.AddressTypePhysical,
		Primary:    true,
	}

	testContactMicro := &mockContactMicro{}
	testContactMicro.On("DeletePersonAddress", mock.Anything, testPersonID, testAddressID).Return(testAddress, error(nil))
	testContactMicro.On("DeletePersonAddress", mock.Anything, testPersonNotFoundID, testAddressID).Return(models.ContactAddress{}, assert.AnError)
	testContactMicro.On("DeletePersonAddress", mock.Anything, testPersonID, testAddressNotFoundID).Return(models.ContactAddress{}, assert.AnError)

	tests := []struct {
		name      string
		ace       adminContactEndpoints
		args      args
		want      PersonAddressData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: DeleteContactRequestData{
					PerosnID:  testPersonID.String(),
					ContactID: testAddressID.String(),
				},
			},
			want: PersonAddressData{
				ID:         testAddressID.String(),
				Street1:    testAddress.Street1,
				Street2:    testAddress.Street2,
				Locality:   testAddress.Locality,
				Region:     testAddress.Region,
				PostalCode: testAddress.PostalCode,
				Country:    testAddress.Country,
				Type:       string(testAddress.Type),
				Primary:    testAddress.Primary,
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; Person DNE",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: DeleteContactRequestData{
					PerosnID:  testPersonNotFoundID.String(),
					ContactID: testAddressID.String(),
				},
			},
		},
		{
			name: "Error; Address DNE",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: DeleteContactRequestData{
					PerosnID:  testPersonID.String(),
					ContactID: testAddressNotFoundID.String(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ace.DeleteContactAddressForPerson(tt.args.ctx, tt.args.reqData)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_adminContactEndpoints_DeleteContactEmailForPerson(t *testing.T) {
	type args struct {
		ctx     context.Context
		reqData DeleteContactRequestData
	}

	testPersonID := uuid.New()
	testPersonNotFoundID := uuid.New()

	testEmailID := uuid.New()
	testEmailNotFoundID := uuid.New()

	testEmail := models.ContactEmail{
		ID:       testEmailID,
		PersonID: testPersonID,
		Username: "",
		Provider: "",
		Primary:  true,
	}

	testContactMicro := &mockContactMicro{}
	testContactMicro.On("DeletePersonEmail", mock.Anything, testPersonID, testEmailID).Return(testEmail, error(nil))
	testContactMicro.On("DeletePersonEmail", mock.Anything, testPersonNotFoundID, testEmailID).Return(models.ContactEmail{}, assert.AnError)
	testContactMicro.On("DeletePersonEmail", mock.Anything, testPersonID, testEmailNotFoundID).Return(models.ContactEmail{}, assert.AnError)

	tests := []struct {
		name      string
		ace       adminContactEndpoints
		args      args
		want      PersonEmailData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: DeleteContactRequestData{
					PerosnID:  testPersonID.String(),
					ContactID: testEmailID.String(),
				},
			},
			want: PersonEmailData{
				ID:      testEmailID.String(),
				Email:   testEmail.String(),
				Primary: testEmail.Primary,
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; Person DNE",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: DeleteContactRequestData{
					PerosnID:  testPersonNotFoundID.String(),
					ContactID: testEmailID.String(),
				},
			},
		},
		{
			name: "Error; Address DNE",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: DeleteContactRequestData{
					PerosnID:  testPersonID.String(),
					ContactID: testEmailNotFoundID.String(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ace.DeleteContactEmailForPerson(tt.args.ctx, tt.args.reqData)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_adminContactEndpoints_DeleteContactPhoneForPerson(t *testing.T) {
	type args struct {
		ctx     context.Context
		reqData DeleteContactRequestData
	}

	testPersonID := uuid.New()
	testPersonNotFoundID := uuid.New()

	testPhoneID := uuid.New()
	testPhoneNotFoundID := uuid.New()

	testPhone := models.ContactPhone{
		ID:          testPhoneID,
		PersonID:    testPersonID,
		CountryCode: 1,
		PhoneNumber: "",
		Type:        models.PhoneTypeHome,
		Primary:     true,
	}

	testContactMicro := &mockContactMicro{}
	testContactMicro.On("DeletePersonPhone", mock.Anything, testPersonID, testPhoneID).Return(testPhone, error(nil))
	testContactMicro.On("DeletePersonPhone", mock.Anything, testPersonNotFoundID, testPhoneID).Return(models.ContactPhone{}, assert.AnError)
	testContactMicro.On("DeletePersonPhone", mock.Anything, testPersonID, testPhoneNotFoundID).Return(models.ContactPhone{}, assert.AnError)

	tests := []struct {
		name      string
		ace       adminContactEndpoints
		args      args
		want      PersonPhoneData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: DeleteContactRequestData{
					PerosnID:  testPersonID.String(),
					ContactID: testPhoneID.String(),
				},
			},
			want: PersonPhoneData{
				ID:          testPhoneID.String(),
				CountryCode: 1,
				PhoneNumber: testPhone.PhoneNumber,
				Type:        string(testPhone.Type),
				Primary:     testPhone.Primary,
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; Person DNE",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: DeleteContactRequestData{
					PerosnID:  testPersonNotFoundID.String(),
					ContactID: testPhoneID.String(),
				},
			},
		},
		{
			name: "Error; Address DNE",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: DeleteContactRequestData{
					PerosnID:  testPersonID.String(),
					ContactID: testPhoneNotFoundID.String(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ace.DeleteContactPhoneForPerson(tt.args.ctx, tt.args.reqData)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_adminContactEndpoint_GetPersonContacts(t *testing.T) {
	type args struct {
		ctx   context.Context
		idStr string
	}

	testDoesNotExistID := uuid.New()
	testPersonID := uuid.New()

	testAddress := models.ContactAddress{
		ID:         uuid.New(),
		PersonID:   testPersonID,
		Street1:    "123 Test Dr",
		Locality:   "Testville",
		Region:     "Testeria",
		PostalCode: "12345-6789",
		Country:    "Testopia",
		Type:       "physical",
		Primary:    true,
	}
	testEmail := models.ContactEmail{
		ID:       uuid.New(),
		PersonID: testPersonID,
		Username: "test123",
		Provider: "example.com",
		Primary:  true,
	}
	testPhone := models.ContactPhone{
		ID:          uuid.New(),
		PersonID:    testPersonID,
		CountryCode: 1,
		PhoneNumber: "555 555-5555",
		Type:        "home",
		Primary:     true,
	}

	testContacts := models.Contacts{
		Addresses: []models.ContactAddress{testAddress},
		Email:     []models.ContactEmail{testEmail},
		Phone:     []models.ContactPhone{testPhone},
	}

	testContactMicro := &mockContactMicro{}
	testContactMicro.On("GetAllForPerson", mock.Anything, testPersonID).Return(testContacts, error(nil))
	testContactMicro.On("GetAllForPerson", mock.Anything, testDoesNotExistID).Return(models.Contacts{}, assert.AnError)

	tests := []struct {
		name      string
		ape       adminContactEndpoints
		args      args
		want      PersonContactData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ape: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx:   context.Background(),
				idStr: testPersonID.String(),
			},
			want: PersonContactData{
				Addresses: []PersonAddressData{
					{
						ID:         testAddress.ID.String(),
						Street1:    testAddress.Street1,
						Street2:    testAddress.Street2,
						Locality:   testAddress.Locality,
						Region:     testAddress.Region,
						PostalCode: testAddress.PostalCode,
						Country:    testAddress.Country,
						Type:       string(testAddress.Type),
						Primary:    testAddress.Primary,
					},
				},
				Emails: []PersonEmailData{
					{
						ID:      testEmail.ID.String(),
						Email:   testEmail.String(),
						Primary: testEmail.Primary,
					},
				},
				PhoneNumbers: []PersonPhoneData{
					{
						ID:          testPhone.ID.String(),
						CountryCode: testPhone.CountryCode,
						PhoneNumber: testPhone.PhoneNumber,
						Type:        string(testPhone.Type),
						Primary:     testPhone.Primary,
					},
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; Bad ID",
			ape: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx:   context.Background(),
				idStr: "bad_val",
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Service",
			ape: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx:   context.Background(),
				idStr: testDoesNotExistID.String(),
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ape.GetPersonContacts(tt.args.ctx, tt.args.idStr)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_adminContactEndpoint_GetPersonContactAddresses(t *testing.T) {
	type args struct {
		ctx   context.Context
		idStr string
	}

	testDoesNotExistID := uuid.New()
	testPersonID := uuid.New()

	testAddresses := []models.ContactAddress{
		{
			ID:         uuid.New(),
			PersonID:   testPersonID,
			Street1:    "123 Test Dr",
			Locality:   "Testville",
			Region:     "Testeria",
			PostalCode: "12345-6789",
			Country:    "Testopia",
			Type:       "physical",
			Primary:    true,
		},
		{
			ID:         uuid.New(),
			PersonID:   testPersonID,
			Street1:    "P.O. Box 9876",
			Locality:   "Testville",
			Region:     "Testeria",
			PostalCode: "12345-7890",
			Country:    "Testopia",
			Type:       "mailing",
			Primary:    true,
		},
	}

	testContactMicro := &mockContactMicro{}
	testContactMicro.On("GetPersonAddresses", mock.Anything, testPersonID).Return(testAddresses, error(nil))
	testContactMicro.On("GetPersonAddresses", mock.Anything, testDoesNotExistID).Return([]models.ContactAddress(nil), assert.AnError)

	tests := []struct {
		name      string
		ape       adminContactEndpoints
		args      args
		want      []PersonAddressData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ape: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx:   context.Background(),
				idStr: testPersonID.String(),
			},
			want: []PersonAddressData{
				{
					ID:         testAddresses[0].ID.String(),
					Street1:    testAddresses[0].Street1,
					Street2:    testAddresses[0].Street2,
					Locality:   testAddresses[0].Locality,
					Region:     testAddresses[0].Region,
					PostalCode: testAddresses[0].PostalCode,
					Country:    testAddresses[0].Country,
					Type:       string(testAddresses[0].Type),
					Primary:    testAddresses[0].Primary,
				},
				{
					ID:         testAddresses[1].ID.String(),
					Street1:    testAddresses[1].Street1,
					Street2:    testAddresses[1].Street2,
					Locality:   testAddresses[1].Locality,
					Region:     testAddresses[1].Region,
					PostalCode: testAddresses[1].PostalCode,
					Country:    testAddresses[1].Country,
					Type:       string(testAddresses[1].Type),
					Primary:    testAddresses[1].Primary,
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; Bad ID",
			ape: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx:   context.Background(),
				idStr: "bad_val",
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Service",
			ape: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx:   context.Background(),
				idStr: testDoesNotExistID.String(),
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ape.GetPersonContactAddresses(tt.args.ctx, tt.args.idStr)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_adminContactEndpoint_GetPersonContactEmails(t *testing.T) {
	type args struct {
		ctx   context.Context
		idStr string
	}

	testDoesNotExistID := uuid.New()
	testPersonID := uuid.New()
	testEmails := []models.ContactEmail{
		{
			ID:       uuid.New(),
			PersonID: testPersonID,
			Username: "test123",
			Provider: "example.com",
			Primary:  true,
		},
		{
			ID:       uuid.New(),
			PersonID: testPersonID,
			Username: "test123",
			Provider: "example.org",
			Primary:  false,
		},
	}

	testContactMicro := &mockContactMicro{}
	testContactMicro.On("GetPersonEmails", mock.Anything, testPersonID).Return(testEmails, error(nil))
	testContactMicro.On("GetPersonEmails", mock.Anything, testDoesNotExistID).Return([]models.ContactEmail(nil), assert.AnError)

	tests := []struct {
		name      string
		ape       adminContactEndpoints
		args      args
		want      []PersonEmailData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ape:  adminContactEndpoints{contactMicro: testContactMicro},
			args: args{
				ctx:   context.Background(),
				idStr: testPersonID.String(),
			},
			want: []PersonEmailData{
				{
					ID:      testEmails[0].ID.String(),
					Email:   testEmails[0].String(),
					Primary: testEmails[0].Primary,
				},
				{
					ID:      testEmails[1].ID.String(),
					Email:   testEmails[1].String(),
					Primary: testEmails[1].Primary,
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; Bad ID",
			ape: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx:   context.Background(),
				idStr: "bad_val",
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Service",
			ape:  adminContactEndpoints{contactMicro: testContactMicro},
			args: args{
				ctx:   context.Background(),
				idStr: testDoesNotExistID.String(),
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ape.GetPersonContactEmails(tt.args.ctx, tt.args.idStr)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_adminContactEndpoint_GetPersonContactPhones(t *testing.T) {
	type args struct {
		ctx   context.Context
		idStr string
	}

	testDoesNotExistID := uuid.New()
	testPersonID := uuid.New()
	testPhones := []models.ContactPhone{
		{
			ID:          uuid.New(),
			PersonID:    testPersonID,
			CountryCode: 1,
			PhoneNumber: "555 555-5555",
			Type:        "home",
			Primary:     true,
		},
	}

	testContactMicro := &mockContactMicro{}
	testContactMicro.On("GetPersonPhones", mock.Anything, testPersonID).Return(testPhones, error(nil))
	testContactMicro.On("GetPersonPhones", mock.Anything, testDoesNotExistID).Return([]models.ContactPhone{}, assert.AnError)

	tests := []struct {
		name      string
		ape       adminContactEndpoints
		args      args
		want      []PersonPhoneData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ape:  adminContactEndpoints{contactMicro: testContactMicro},
			args: args{
				ctx:   context.Background(),
				idStr: testPersonID.String(),
			},
			want: []PersonPhoneData{
				{
					ID:          testPhones[0].ID.String(),
					CountryCode: testPhones[0].CountryCode,
					PhoneNumber: testPhones[0].PhoneNumber,
					Type:        string(testPhones[0].Type),
					Primary:     testPhones[0].Primary,
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; Bad ID",
			ape: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx:   context.Background(),
				idStr: "bad_val",
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Service",
			ape:  adminContactEndpoints{contactMicro: testContactMicro},
			args: args{
				ctx:   context.Background(),
				idStr: testDoesNotExistID.String(),
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ape.GetPersonContactPhones(tt.args.ctx, tt.args.idStr)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_adminContactEndpoints_UpdatePersonContactAddress(t *testing.T) {
	type args struct {
		ctx     context.Context
		reqData UpdateContactRequestData[PersonAddressData]
	}

	testValidPersonID := uuid.New()
	testNotFoundPersonID := uuid.New()

	testValidAddressID := uuid.New()
	testNotFoundAddressID := uuid.New()

	testValidAddressData := PersonAddressData{
		Street1:    "987 Testing Ln",
		Street2:    "APT 3",
		Locality:   "Testington",
		Region:     "Testoria",
		PostalCode: "98765-4321",
		Country:    "Testopia",
		Type:       "mailing",
		Primary:    true,
	}
	testValidAddressDB := models.ContactAddress{
		Street1:    "987 Testing Ln",
		Street2:    "APT 3",
		Locality:   "Testington",
		Region:     "Testoria",
		PostalCode: "98765-4321",
		Country:    "Testopia",
		Type:       models.AddressTypeMailing,
		Primary:    true,
	}

	testContactMicro := &mockContactMicro{}
	testContactMicro.On("UpdatePersonAddress", mock.Anything, testValidPersonID, testValidAddressID, testValidAddressDB).Return(error(nil))
	testContactMicro.On("UpdatePersonAddress", mock.Anything, testNotFoundPersonID, testValidAddressID, testValidAddressDB).Return(assert.AnError)
	testContactMicro.On("UpdatePersonAddress", mock.Anything, testValidPersonID, testNotFoundAddressID, testValidAddressDB).Return(assert.AnError)

	tests := []struct {
		name      string
		ace       adminContactEndpoints
		args      args
		want      PersonAddressData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: UpdateContactRequestData[PersonAddressData]{
					PersonID:  testValidPersonID.String(),
					ContactID: testValidAddressID.String(),
					Data:      testValidAddressData,
				},
			},
			want:      testValidAddressData,
			assertion: assert.NoError,
		},
		{
			name: "Error; Invalid Person ID",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: UpdateContactRequestData[PersonAddressData]{
					PersonID: "invalid_id",
					Data:     testValidAddressData,
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Person DNE",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: UpdateContactRequestData[PersonAddressData]{
					PersonID:  testNotFoundPersonID.String(),
					ContactID: testValidAddressID.String(),
					Data:      testValidAddressData,
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Address DNE",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: UpdateContactRequestData[PersonAddressData]{
					PersonID:  testValidPersonID.String(),
					ContactID: testNotFoundAddressID.String(),
					Data:      testValidAddressData,
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Invalid Address; Missing Street 1",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: UpdateContactRequestData[PersonAddressData]{
					PersonID:  testValidPersonID.String(),
					ContactID: testValidAddressID.String(),
					Data: PersonAddressData{
						Locality:   "Testetville",
						Region:     "Testaria",
						PostalCode: "12345-6789",
						Country:    "Testopia",
						Type:       "mailing",
						Primary:    true,
					},
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Invalid Address; Missing Locality",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: UpdateContactRequestData[PersonAddressData]{
					PersonID:  testValidPersonID.String(),
					ContactID: testValidAddressID.String(),
					Data: PersonAddressData{
						Locality:   "Testetville",
						Region:     "Testaria",
						PostalCode: "12345-6789",
						Country:    "Testopia",
						Type:       "mailing",
						Primary:    true,
					},
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Invalid Address; Missing Region",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: UpdateContactRequestData[PersonAddressData]{
					PersonID:  testValidPersonID.String(),
					ContactID: testValidAddressID.String(),
					Data: PersonAddressData{
						Street1:    "123 Test Dr",
						Locality:   "Testetville",
						PostalCode: "12345-6789",
						Country:    "Testopia",
						Type:       "mailing",
						Primary:    true,
					},
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Invalid Address; Missing Postal Code",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: UpdateContactRequestData[PersonAddressData]{
					PersonID:  testValidPersonID.String(),
					ContactID: testValidAddressID.String(),
					Data: PersonAddressData{
						Street1:  "123 Test Dr",
						Locality: "Testetville",
						Region:   "Testaria",
						Country:  "Testopia",
						Type:     "mailing",
						Primary:  true,
					},
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Invalid Address; Missing Country",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: UpdateContactRequestData[PersonAddressData]{
					PersonID:  testValidPersonID.String(),
					ContactID: testValidAddressID.String(),
					Data: PersonAddressData{
						Street1:    "123 Test Dr",
						Locality:   "Testetville",
						Region:     "Testaria",
						PostalCode: "12345-6789",
						Type:       "mailing",
						Primary:    true,
					},
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Invalid Address; Missing Type",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: UpdateContactRequestData[PersonAddressData]{
					PersonID:  testValidPersonID.String(),
					ContactID: testValidAddressID.String(),
					Data: PersonAddressData{
						Street1:    "123 Test Dr",
						Locality:   "Testetville",
						Region:     "Testaria",
						PostalCode: "12345-6789",
						Country:    "Testopia",
						Primary:    true,
					},
				},
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ace.UpdatePersonContactAddress(tt.args.ctx, tt.args.reqData)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_adminContactEndpoints_UpdatePersonContactEmail(t *testing.T) {
	type args struct {
		ctx     context.Context
		reqData UpdateContactRequestData[PersonEmailData]
	}

	testValidPersonID := uuid.New()
	testNotFoundPersonID := uuid.New()
	testValidEmailID := uuid.New()
	testNotFoundEmailID := uuid.New()

	testValidEmailData := PersonEmailData{
		Email:   "test@example.com",
		Primary: true,
	}
	testValidEmailDB := models.ContactEmail{
		Username: "test",
		Provider: "example.com",
		Primary:  true,
	}

	testInvalidEmailData1 := PersonEmailData{
		Email: "@example.com",
	}
	testInvalidEmailData2 := PersonEmailData{
		Email: "user@",
	}
	testInvalidEmailData3 := PersonEmailData{
		Email: "user@example",
	}

	testContactMicro := &mockContactMicro{}
	testContactMicro.On("UpdatePersonEmail", mock.Anything, testValidPersonID, testValidEmailID, testValidEmailDB).Return(error(nil))
	testContactMicro.On("UpdatePersonEmail", mock.Anything, testNotFoundPersonID, testValidEmailID, testValidEmailDB).Return(assert.AnError)
	testContactMicro.On("UpdatePersonEmail", mock.Anything, testValidPersonID, testNotFoundEmailID, testValidEmailDB).Return(assert.AnError)

	tests := []struct {
		name      string
		ace       adminContactEndpoints
		args      args
		want      PersonEmailData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: UpdateContactRequestData[PersonEmailData]{
					PersonID:  testValidPersonID.String(),
					ContactID: testValidEmailID.String(),
					Data:      testValidEmailData,
				},
			},
			want:      testValidEmailData,
			assertion: assert.NoError,
		},
		{
			name: "Error; Invalid Person ID",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: UpdateContactRequestData[PersonEmailData]{
					PersonID: "invalid_id",
					Data:     testValidEmailData,
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Person DNE",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: UpdateContactRequestData[PersonEmailData]{
					PersonID:  testNotFoundPersonID.String(),
					ContactID: testValidEmailID.String(),
					Data:      testValidEmailData,
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Email DNE",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: UpdateContactRequestData[PersonEmailData]{
					PersonID:  testValidPersonID.String(),
					ContactID: testNotFoundEmailID.String(),
					Data:      testValidEmailData,
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; No Email Given",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: UpdateContactRequestData[PersonEmailData]{
					PersonID:  testValidPersonID.String(),
					ContactID: testValidEmailID.String(),
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Invalid Email; No Username",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: UpdateContactRequestData[PersonEmailData]{
					PersonID:  testValidPersonID.String(),
					ContactID: testValidEmailID.String(),
					Data:      testInvalidEmailData1,
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Invalid Email; No Domain",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: UpdateContactRequestData[PersonEmailData]{
					PersonID:  testValidPersonID.String(),
					ContactID: testValidEmailID.String(),
					Data:      testInvalidEmailData2,
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Invalid Email; Missing Top-level Domain",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: UpdateContactRequestData[PersonEmailData]{
					PersonID:  testValidPersonID.String(),
					ContactID: testValidEmailID.String(),
					Data:      testInvalidEmailData3,
				},
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ace.UpdatePersonContactEmail(tt.args.ctx, tt.args.reqData)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_adminContactEndpoints_UpdatePersonContactPhone(t *testing.T) {
	type args struct {
		ctx     context.Context
		reqData UpdateContactRequestData[PersonPhoneData]
	}

	testValidPersonID := uuid.New()
	testNotFoundPersonID := uuid.New()
	testValidPhoneID := uuid.New()
	testNotFoundPhoneID := uuid.New()

	testValidPhoneData := PersonPhoneData{
		CountryCode: 1,
		PhoneNumber: "(555)555-5555",
		Type:        "home",
		Primary:     true,
	}
	testValidPhoneDB := models.ContactPhone{
		CountryCode: 1,
		PhoneNumber: "(555)555-5555",
		Type:        models.PhoneTypeHome,
		Primary:     true,
	}

	testContactMicro := &mockContactMicro{}
	testContactMicro.On("UpdatePersonPhone", mock.Anything, testValidPersonID, testValidPhoneID, testValidPhoneDB).Return(error(nil))
	testContactMicro.On("UpdatePersonPhone", mock.Anything, testValidPersonID, testNotFoundPhoneID, testValidPhoneDB).Return(assert.AnError)
	testContactMicro.On("UpdatePersonPhone", mock.Anything, testNotFoundPersonID, testValidPhoneID, testValidPhoneDB).Return(assert.AnError)

	tests := []struct {
		name      string
		ace       adminContactEndpoints
		args      args
		want      PersonPhoneData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: UpdateContactRequestData[PersonPhoneData]{
					PersonID:  testValidPersonID.String(),
					ContactID: testValidPhoneID.String(),
					Data:      testValidPhoneData,
				},
			},
			want:      testValidPhoneData,
			assertion: assert.NoError,
		},
		{
			name: "Error; Person DNE",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: UpdateContactRequestData[PersonPhoneData]{
					PersonID:  testNotFoundPersonID.String(),
					ContactID: testValidPhoneID.String(),
					Data:      testValidPhoneData,
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Invalid Person ID",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: UpdateContactRequestData[PersonPhoneData]{
					PersonID: "invalid_id",
					Data:     testValidPhoneData,
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Phone DNE",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: UpdateContactRequestData[PersonPhoneData]{
					PersonID:  testValidPersonID.String(),
					ContactID: testNotFoundPhoneID.String(),
					Data:      testValidPhoneData,
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Invalid Phone; Bad Country Code",
			ace: adminContactEndpoints{
				contactMicro: testContactMicro,
			},
			args: args{
				ctx: context.Background(),
				reqData: UpdateContactRequestData[PersonPhoneData]{
					PersonID:  testValidPersonID.String(),
					ContactID: testValidPhoneID.String(),
					Data: PersonPhoneData{
						CountryCode: -1,
						PhoneNumber: "555-555-5555",
						Type:        "home",
						Primary:     false,
					},
				},
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ace.UpdatePersonContactPhone(tt.args.ctx, tt.args.reqData)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
