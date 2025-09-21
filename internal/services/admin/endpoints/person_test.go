package endpoints

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/williabk198/timeclock/internal/models"
)

func Test_personEndpoints_Add(t *testing.T) {
	type args struct {
		ctx    context.Context
		person PersonData
	}

	testGoodPersonID := uuid.New()
	testGoodPerson := PersonData{
		Name: models.Name{
			GivenName:       "Testy",
			FamilyName:      "McTesterson",
			FamilyNameFirst: models.FirstNameGiven,
		},
		DateOfBirth: 0,
		Gender:      "non-binary",
		Pronouns:    "they/them",
	}
	testGoodPersonDB := models.Person{
		Name:        testGoodPerson.Name,
		DateOfBirth: time.Unix(testGoodPerson.DateOfBirth, 0),
		Gender:      models.Gender(testGoodPerson.Gender),
		Pronouns:    models.Pronouns{Subject: "they", Object: "them"},
	}

	testBadPerson := PersonData{
		Name: models.Name{
			GivenName:       "Testy",
			FamilyName:      "McTesterson",
			FamilyNameFirst: models.FirstNameGiven,
		},
		DateOfBirth: 0,
		Gender:      "error value",
		Pronouns:    "they/them",
	}
	testBadPersonDB := models.Person{
		Name:        testBadPerson.Name,
		DateOfBirth: time.Unix(testBadPerson.DateOfBirth, 0),
		Gender:      models.Gender(testBadPerson.Gender),
		Pronouns:    models.Pronouns{Subject: "they", Object: "them"},
	}

	mockAdminService := &mockAdminService{}

	mockAdminService.On("AddPerson", mock.Anything, testGoodPersonDB).Return(testGoodPersonID, error(nil))
	mockAdminService.On("AddPerson", mock.Anything, testBadPersonDB).Return(uuid.Nil, assert.AnError)

	tests := []struct {
		name      string
		ape       adminPersonEndpoints
		args      args
		want      PersonData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ape:  adminPersonEndpoints{adminService: mockAdminService},
			args: args{
				ctx:    context.Background(),
				person: testGoodPerson,
			},
			want: PersonData{
				ID: testGoodPersonID.String(),
				Name: models.Name{
					GivenName:       "Testy",
					FamilyName:      "McTesterson",
					FamilyNameFirst: models.FirstNameGiven,
				},
				DateOfBirth: 0,
				Gender:      "non-binary",
				Pronouns:    "they/them",
			},
			assertion: assert.NoError,
		},
		{
			name: "Error: Invalid Input",
			ape:  adminPersonEndpoints{adminService: mockAdminService},
			args: args{
				ctx:    context.Background(),
				person: PersonData{Pronouns: "invalid format"},
			},
			want:      PersonData{},
			assertion: assert.Error,
		},
		{
			name: "Error: Service Error",
			ape:  adminPersonEndpoints{adminService: mockAdminService},
			args: args{
				ctx:    context.Background(),
				person: testBadPerson,
			},
			want:      PersonData{},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ape.Add(tt.args.ctx, tt.args.person)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_adminPersonEndpoints_Delete(t *testing.T) {
	type args struct {
		ctx   context.Context
		idStr string
	}

	testDoesNotExistID := uuid.New()
	testPersonID := uuid.New()
	testPerson := models.Person{
		Name: models.Name{
			GivenName:       "Testy",
			FamilyName:      "McTesterson",
			FamilyNameFirst: models.FirstNameGiven,
		},
		DateOfBirth: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		Gender:      models.GenderNonBinary,
		Pronouns: models.Pronouns{
			Subject: "they",
			Object:  "them",
		},
	}

	testAdminService := &mockAdminService{}
	testAdminService.On("DeletePerson", mock.Anything, testPersonID).Return(testPerson, error(nil))
	testAdminService.On("DeletePerson", mock.Anything, testDoesNotExistID).Return(models.Person{}, assert.AnError)

	tests := []struct {
		name      string
		ape       adminPersonEndpoints
		args      args
		want      PersonData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ape: adminPersonEndpoints{
				adminService: testAdminService,
			},
			args: args{
				ctx:   context.Background(),
				idStr: testPersonID.String(),
			},
			want: PersonData{
				ID:          testPersonID.String(),
				Name:        testPerson.Name,
				DateOfBirth: testPerson.DateOfBirth.Unix(),
				Gender:      string(testPerson.Gender),
				Pronouns:    testPerson.Pronouns.String(),
			},
			assertion: assert.NoError,
		},
		{
			name: "Error",
			ape: adminPersonEndpoints{
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
			got, err := tt.ape.Delete(tt.args.ctx, tt.args.idStr)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_adminPersonEndpoints_GetAll(t *testing.T) {
	type args struct {
		ctx     context.Context
		reqData GetPaginatedRequestData
	}

	testPersons := []models.Person{
		{
			ID:          uuid.New(),
			Name:        models.Name{GivenName: "Testy", FamilyName: "McTesterson", FamilyNameFirst: models.FirstNameGiven},
			DateOfBirth: time.Unix(0, 0),
			Gender:      models.GenderNonBinary,
			Pronouns:    models.Pronouns{Subject: "they", Object: "them"},
		},
		{
			ID:          uuid.New(),
			Name:        models.Name{GivenName: "Brandon", FamilyName: "Williams", FamilyNameFirst: models.FirstNameGiven},
			DateOfBirth: time.Date(1992, 1, 27, 0, 0, 0, 0, time.UTC),
			Gender:      models.GenderMale,
			Pronouns:    models.Pronouns{Subject: "he", Object: "him"},
		},
		{
			ID:          uuid.New(),
			Name:        models.Name{GivenName: "Testita", FamilyName: "Tester", FamilyNameFirst: models.FirstNameGiven},
			DateOfBirth: time.Date(1950, 1, 1, 0, 0, 0, 0, time.UTC),
			Gender:      models.GenderFemale,
			Pronouns:    models.Pronouns{Subject: "she", Object: "her"},
		},
	}
	resultData := []PersonData{
		{
			ID:          testPersons[1].ID.String(),
			Name:        testPersons[1].Name,
			DateOfBirth: testPersons[1].DateOfBirth.Unix(),
			Gender:      string(testPersons[1].Gender),
			Pronouns:    testPersons[1].Pronouns.String(),
		},
		{
			ID:          testPersons[2].ID.String(),
			Name:        testPersons[2].Name,
			DateOfBirth: testPersons[2].DateOfBirth.Unix(),
			Gender:      string(testPersons[2].Gender),
			Pronouns:    testPersons[2].Pronouns.String(),
		},
	}

	testAdminService := &mockAdminService{}
	testAdminService.On("GetAllPersons", mock.Anything, uint(1), uint(2)).Return(testPersons[1:3], error(nil))
	testAdminService.On("GetAllPersons", mock.Anything, uint(0), uint(0)).Return([]models.Person{}, assert.AnError)

	tests := []struct {
		name      string
		ape       adminPersonEndpoints
		args      args
		want      []PersonData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ape: adminPersonEndpoints{
				adminService: testAdminService,
			},
			args: args{
				ctx: context.Background(),
				reqData: GetPaginatedRequestData{
					Offset: 1,
					Limit:  2,
				},
			},
			want:      resultData,
			assertion: assert.NoError,
		},
		{
			name: "Error",
			ape: adminPersonEndpoints{
				adminService: testAdminService,
			},
			args: args{
				ctx: context.Background(),
				reqData: GetPaginatedRequestData{
					Offset: 0,
					Limit:  0,
				},
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ape.GetAll(tt.args.ctx, tt.args.reqData)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_adminPersonEndpoints_GetSpecific(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}

	testDoesNotExistID := uuid.New()
	testPersonID := uuid.New()
	testPerson := models.Person{
		Name: models.Name{
			GivenName:       "Testy",
			FamilyName:      "McTesterson",
			FamilyNameFirst: models.FirstNameGiven,
		},
		DateOfBirth: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		Gender:      models.GenderNonBinary,
		Pronouns: models.Pronouns{
			Subject: "they",
			Object:  "them",
		},
	}
	testPersonData := PersonData{
		ID:          testPersonID.String(),
		Name:        testPerson.Name,
		DateOfBirth: testPerson.DateOfBirth.Unix(),
		Gender:      string(testPerson.Gender),
		Pronouns:    testPerson.Pronouns.String(),
	}

	testAdminService := &mockAdminService{}
	testAdminService.On("GetPerson", mock.Anything, testPersonID).Return(testPerson, error(nil))
	testAdminService.On("GetPerson", mock.Anything, testDoesNotExistID).Return(models.Person{}, assert.AnError)

	tests := []struct {
		name      string
		ape       adminPersonEndpoints
		args      args
		want      PersonData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ape: adminPersonEndpoints{
				adminService: testAdminService,
			},
			args: args{
				ctx: context.Background(),
				id:  testPersonID.String(),
			},
			want:      testPersonData,
			assertion: assert.NoError,
		},
		{
			name: "Error",
			ape: adminPersonEndpoints{
				adminService: testAdminService,
			},
			args: args{
				ctx: context.Background(),
				id:  testDoesNotExistID.String(),
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ape.GetSpecific(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_adminPersonEndpoints_GetSpecificContacts(t *testing.T) {
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
		ape       adminPersonEndpoints
		args      args
		want      PersonContactData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ape: adminPersonEndpoints{
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
			ape: adminPersonEndpoints{
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
			ape: adminPersonEndpoints{
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
			got, err := tt.ape.GetSpecificContacts(tt.args.ctx, tt.args.idStr)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_adminPersonEndpoints_GetSpecificContactAddresses(t *testing.T) {
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
	testAdminService.On("GetPersonContactAddress", mock.Anything, testPersonID).Return(testAddresses, error(nil))
	testAdminService.On("GetPersonContactAddress", mock.Anything, testDoesNotExistID).Return([]models.ContactAddress(nil), assert.AnError)

	tests := []struct {
		name      string
		ape       adminPersonEndpoints
		args      args
		want      []PersonAddressData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ape: adminPersonEndpoints{
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
			ape: adminPersonEndpoints{
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
			ape: adminPersonEndpoints{
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
			got, err := tt.ape.GetSpecificContactAddresses(tt.args.ctx, tt.args.idStr)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_adminPersonEndpoints_GetSpecificContactEmails(t *testing.T) {
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
		ape       adminPersonEndpoints
		args      args
		want      []PersonEmailData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ape:  adminPersonEndpoints{adminService: testAdminService},
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
		},
		{
			name: "Error; Bad ID",
			ape: adminPersonEndpoints{
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
			ape:  adminPersonEndpoints{adminService: testAdminService},
			args: args{
				ctx:   context.Background(),
				idStr: testDoesNotExistID.String(),
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ape.GetSpecificContactEmails(tt.args.ctx, tt.args.idStr)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_adminPersonEndpoints_GetSpecificContactPhones(t *testing.T) {
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
	testAdminService.On("GetPersonContacts", mock.Anything, testPersonID).Return(testPhones, error(nil))
	testAdminService.On("GetPersonContacts", mock.Anything, testDoesNotExistID).Return([]models.ContactPhone{}, assert.AnError)

	tests := []struct {
		name      string
		ape       adminPersonEndpoints
		args      args
		want      []PersonPhoneData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ape:  adminPersonEndpoints{adminService: testAdminService},
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
			ape: adminPersonEndpoints{
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
			ape:  adminPersonEndpoints{adminService: testAdminService},
			args: args{
				ctx:   context.Background(),
				idStr: testDoesNotExistID.String(),
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ape.GetSpecificContactPhones(tt.args.ctx, tt.args.idStr)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_adminPersonEndpoints_Update(t *testing.T) {
	type args struct {
		ctx context.Context
		urd UpdateRequestData[PersonData]
	}
	testPersonID := uuid.New()
	testPerson := models.Person{
		Name: models.Name{
			GivenName:       "Tetsuya",
			FamilyName:      "Takahashi",
			FamilyNameFirst: models.FirstNameFamily,
		},
		DateOfBirth: time.Date(1966, 11, 18, 0, 0, 0, 0, time.UTC).Local(),
		Gender:      models.GenderMale,
		Pronouns: models.Pronouns{
			Subject: "he",
			Object:  "him",
		},
	}
	testPersonData := PersonData{
		Name:        testPerson.Name,
		DateOfBirth: testPerson.DateOfBirth.Unix(),
		Gender:      string(testPerson.Gender),
		Pronouns:    testPerson.Pronouns.String(),
	}
	testBadPerson := models.Person{
		Name: models.Name{
			GivenName:       "Testy",
			FamilyName:      "McTesterson",
			FamilyNameFirst: models.FirstNameGiven,
		},
	}
	testBadPersonData := PersonData{
		Name: testBadPerson.Name,
	}

	testAdminService := &mockAdminService{}
	testAdminService.On("UpdatePerson", mock.Anything, testPersonID, testPerson).Return(error(nil))
	testAdminService.On("UpdatePerson", mock.Anything, testPersonID, testBadPerson).Return(assert.AnError)

	tests := []struct {
		name      string
		ape       adminPersonEndpoints
		args      args
		want      PersonData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ape: adminPersonEndpoints{
				adminService: testAdminService,
			},
			args: args{
				ctx: context.Background(),
				urd: UpdateRequestData[PersonData]{
					ID:   testPersonID.String(),
					Data: testPersonData,
				},
			},
			want:      testPersonData,
			assertion: assert.NoError,
		},
		{
			name: "Error; Bad ID Value",
			ape: adminPersonEndpoints{
				adminService: testAdminService,
			},
			args: args{
				ctx: context.Background(),
				urd: UpdateRequestData[PersonData]{
					ID: "bad_value",
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Service",
			ape: adminPersonEndpoints{
				adminService: testAdminService,
			},
			args: args{
				ctx: context.Background(),
				urd: UpdateRequestData[PersonData]{
					ID:   testPersonID.String(),
					Data: testBadPersonData,
				},
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ape.Update(tt.args.ctx, tt.args.urd)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

type mockAdminService struct {
	mock.Mock
}

func (mas *mockAdminService) AddPerson(ctx context.Context, person models.Person) (uuid.UUID, error) {
	args := mas.Called(ctx, person)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (mas *mockAdminService) DeletePerson(ctx context.Context, id uuid.UUID) (models.Person, error) {
	args := mas.Called(ctx, id)
	return args.Get(0).(models.Person), args.Error(1)
}

func (mas *mockAdminService) GetAllPersons(ctx context.Context, offset, limit uint) ([]models.Person, error) {
	args := mas.Called(ctx, offset, limit)
	return args.Get(0).([]models.Person), args.Error(1)
}

func (mas *mockAdminService) GetPerson(ctx context.Context, id uuid.UUID) (models.Person, error) {
	args := mas.Called(ctx, id)
	return args.Get(0).(models.Person), args.Error(1)
}

func (mas *mockAdminService) GetPersonContacts(ctx context.Context, id uuid.UUID) (models.Contacts, error) {
	args := mas.Called(ctx, id)
	return args.Get(0).(models.Contacts), args.Error(1)
}

// GetPersonContactAddresses implements Service.
func (mas *mockAdminService) GetPersonContactAddresses(ctx context.Context, id uuid.UUID) ([]models.ContactAddress, error) {
	args := mas.Called(ctx, id)
	return args.Get(0).([]models.ContactAddress), args.Error(1)
}

// GetPersonContactEmails implements Service.
func (mas *mockAdminService) GetPersonContactEmails(ctx context.Context, id uuid.UUID) ([]models.ContactEmail, error) {
	args := mas.Called(ctx, id)
	return args.Get(0).([]models.ContactEmail), args.Error(1)
}

// GetPersonContactPhones implements Service.
func (mas *mockAdminService) GetPersonContactPhones(ctx context.Context, id uuid.UUID) ([]models.ContactPhone, error) {
	args := mas.Called(ctx, id)
	return args.Get(0).([]models.ContactPhone), args.Error(1)
}

func (mas *mockAdminService) UpdatePerson(ctx context.Context, id uuid.UUID, data models.Person) error {
	args := mas.Called(ctx, id, data)
	return args.Error(0)
}
