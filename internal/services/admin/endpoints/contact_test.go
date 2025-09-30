package endpoints

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/williabk198/timeclock/internal/models"
)

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

	testAdminService := &mockAdminService{}
	testAdminService.On("AddPersonContactEmail", mock.Anything, testValidEmailDB).Return(testEmailID, error(nil))
	testAdminService.On("AddPersonContactEmail", mock.Anything, testInvalidEmailDB).Return(uuid.Nil, assert.AnError)

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
				adminService: testAdminService,
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
				adminService: testAdminService,
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
				adminService: testAdminService,
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
				adminService: testAdminService,
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
				adminService: testAdminService,
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
				adminService: testAdminService,
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
				adminService: testAdminService,
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

	testAdminService := &mockAdminService{}
	testAdminService.On("GetPersonContacts", mock.Anything, testPersonID).Return(testContacts, error(nil))
	testAdminService.On("GetPersonContacts", mock.Anything, testDoesNotExistID).Return(models.Contacts{}, assert.AnError)

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
				adminService: testAdminService,
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
						Type:       testAddress.Type,
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
						Type:        testPhone.Type,
						Primary:     testPhone.Primary,
					},
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; Bad ID",
			ape: adminContactEndpoints{
				adminService: testAdminService,
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
				adminService: testAdminService,
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

	testAdminService := &mockAdminService{}
	testAdminService.On("GetPersonContactAddresses", mock.Anything, testPersonID).Return(testAddresses, error(nil))
	testAdminService.On("GetPersonContactAddresses", mock.Anything, testDoesNotExistID).Return([]models.ContactAddress(nil), assert.AnError)

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
				adminService: testAdminService,
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
					Type:       testAddresses[0].Type,
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
					Type:       testAddresses[1].Type,
					Primary:    testAddresses[1].Primary,
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; Bad ID",
			ape: adminContactEndpoints{
				adminService: testAdminService,
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
				adminService: testAdminService,
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

	testAdminService := &mockAdminService{}
	testAdminService.On("GetPersonContactEmails", mock.Anything, testPersonID).Return(testEmails, error(nil))
	testAdminService.On("GetPersonContactEmails", mock.Anything, testDoesNotExistID).Return([]models.ContactEmail(nil), assert.AnError)

	tests := []struct {
		name      string
		ape       adminContactEndpoints
		args      args
		want      []PersonEmailData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ape:  adminContactEndpoints{adminService: testAdminService},
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
				adminService: testAdminService,
			},
			args: args{
				ctx:   context.Background(),
				idStr: "bad_val",
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Service",
			ape:  adminContactEndpoints{adminService: testAdminService},
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

	testAdminService := &mockAdminService{}
	testAdminService.On("GetPersonContactPhones", mock.Anything, testPersonID).Return(testPhones, error(nil))
	testAdminService.On("GetPersonContactPhones", mock.Anything, testDoesNotExistID).Return([]models.ContactPhone{}, assert.AnError)

	tests := []struct {
		name      string
		ape       adminContactEndpoints
		args      args
		want      []PersonPhoneData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ape:  adminContactEndpoints{adminService: testAdminService},
			args: args{
				ctx:   context.Background(),
				idStr: testPersonID.String(),
			},
			want: []PersonPhoneData{
				{
					ID:          testPhones[0].ID.String(),
					CountryCode: testPhones[0].CountryCode,
					PhoneNumber: testPhones[0].PhoneNumber,
					Type:        testPhones[0].Type,
					Primary:     testPhones[0].Primary,
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; Bad ID",
			ape: adminContactEndpoints{
				adminService: testAdminService,
			},
			args: args{
				ctx:   context.Background(),
				idStr: "bad_val",
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Service",
			ape:  adminContactEndpoints{adminService: testAdminService},
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
