package admin

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/williabk198/timeclock/internal/datastores"
	"github.com/williabk198/timeclock/internal/models"
)

func TestNewService(t *testing.T) {
	type args struct {
		adminStore datastores.PersonStore
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
			assert.Equal(t, tt.want, NewService(tt.args.adminStore))
		})
	}
}

func Test_authService_AddPerson(t *testing.T) {
	type args struct {
		ctx    context.Context
		person models.Person
	}

	type wants struct {
		id uuid.UUID
	}

	testPersonID := uuid.New()
	testPerson := models.Person{
		Name: models.Name{
			GivenName:       "Testy",
			FamilyName:      "McTesterson",
			FamilyNameFirst: models.FirstNameGiven,
		},
		DateOfBirth: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		Gender:      models.GenderNonBinary,
		Pronouns:    models.Pronouns{Subject: "they", Object: "them"},
	}
	testErrorPerson := models.Person{
		Name: models.Name{
			GivenName: "err_val",
		},
	}

	testPersonStore := &mockPersonStore{}
	testPersonStore.On("Add", mock.Anything, testPerson).Return(testPersonID, error(nil))
	testPersonStore.On("Add", mock.Anything, testErrorPerson).Return(uuid.Nil, assert.AnError)

	tests := []struct {
		name      string
		a         adminService
		args      args
		wants     wants
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			a: adminService{
				personStore: testPersonStore,
			},
			args: args{
				ctx:    context.Background(),
				person: testPerson,
			},
			wants: wants{
				id: testPersonID,
			},
			assertion: assert.NoError,
		},
		{
			name: "Error",
			a: adminService{
				personStore: testPersonStore,
			},
			args: args{
				ctx:    context.Background(),
				person: testErrorPerson,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, err := tt.a.AddPerson(tt.args.ctx, tt.args.person)
			tt.assertion(t, err)
			assert.Equal(t, tt.wants.id, gotID)
		})
	}
}

func Test_adminService_DeletePerson(t *testing.T) {
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}

	testNotFoundID := uuid.New()
	testPersonID := uuid.New()
	testPerson := models.Person{
		Name: models.Name{
			GivenName:       "Testy",
			FamilyName:      "McTesterson",
			FamilyNameFirst: models.FirstNameGiven,
		},
		DateOfBirth: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		Gender:      models.GenderNonBinary,
		Pronouns:    models.Pronouns{Subject: "they", Object: "them"},
	}

	testPersonStore := &mockPersonStore{}
	testPersonStore.On("Delete", mock.Anything, testPersonID).Return(testPerson, error(nil))
	testPersonStore.On("Delete", mock.Anything, testNotFoundID).Return(models.Person{}, assert.AnError)

	tests := []struct {
		name      string
		as        adminService
		args      args
		want      models.Person
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			as: adminService{
				personStore: testPersonStore,
			},
			args: args{
				ctx: context.Background(),
				id:  testPersonID,
			},
			want:      testPerson,
			assertion: assert.NoError,
		},
		{
			name: "Error",
			as: adminService{
				personStore: testPersonStore,
			},
			args: args{
				ctx: context.Background(),
				id:  testNotFoundID,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.as.DeletePerson(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_adminService_GetAllPersons(t *testing.T) {
	type args struct {
		ctx    context.Context
		offset uint
		limit  uint
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

	testPersonStore := &mockPersonStore{}
	testPersonStore.On("GetAllPaginated", mock.Anything, uint(0), uint(2)).Return(testPersons[:2], error(nil))
	testPersonStore.On("GetAllPaginated", mock.Anything, uint(0), uint(0)).Return([]models.Person(nil), assert.AnError)

	tests := []struct {
		name      string
		as        adminService
		args      args
		want      []models.Person
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			as: adminService{
				personStore: testPersonStore,
			},
			args: args{
				ctx:    context.Background(),
				offset: 0,
				limit:  2,
			},
			want:      testPersons[:2],
			assertion: assert.NoError,
		},
		{
			name: "Error",
			as: adminService{
				personStore: testPersonStore,
			},
			args: args{
				ctx:    context.Background(),
				offset: 0,
				limit:  0,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.as.GetAllPersons(tt.args.ctx, tt.args.offset, tt.args.limit)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_adminService_GetPerson(t *testing.T) {
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}

	testNotFoundID := uuid.New()
	testPersonID := uuid.New()
	testPerson := models.Person{
		Name: models.Name{
			GivenName:       "Testy",
			FamilyName:      "McTesterson",
			FamilyNameFirst: models.FirstNameGiven,
		},
		DateOfBirth: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		Gender:      models.GenderNonBinary,
		Pronouns:    models.Pronouns{Subject: "they", Object: "them"},
	}

	testPersonStore := &mockPersonStore{}
	testPersonStore.On("GetSpecific", mock.Anything, testPersonID).Return(testPerson, error(nil))
	testPersonStore.On("GetSpecific", mock.Anything, testNotFoundID).Return(models.Person{}, assert.AnError)

	tests := []struct {
		name      string
		as        adminService
		args      args
		want      models.Person
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			as: adminService{
				personStore: testPersonStore,
			},
			args: args{
				ctx: context.Background(),
				id:  testPersonID,
			},
			want:      testPerson,
			assertion: assert.NoError,
		},
		{
			name: "Error",
			as: adminService{
				personStore: testPersonStore,
			},
			args: args{
				ctx: context.Background(),
				id:  testNotFoundID,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.as.GetPerson(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_adminService_GetPersonContacts(t *testing.T) {
	type args struct {
		ctx      context.Context
		personID uuid.UUID
	}

	type wants struct {
		callAssertions map[string]int
		returnVal      models.Contacts
	}

	testPersonID := uuid.New()
	testPersonAddressID := uuid.New()
	testPersonEmailID := uuid.New()
	testPersonPhoneID := uuid.New()
	testBadPersonAddressID := uuid.New()
	testBadPersonEmailID := uuid.New()
	testBadPersonPhoneID := uuid.New()

	testPersonAddresses := []models.ContactAddress{
		{
			ID:         testPersonAddressID,
			PersonID:   testPersonID,
			Street1:    "123 Test Dr",
			Locality:   "Testville",
			Region:     "Testeria",
			PostalCode: "12345-6789",
			Country:    "Testopia",
			Type:       "physical",
			Primary:    true,
		},
	}
	testPersonEmails := []models.ContactEmail{
		{
			ID:       testPersonEmailID,
			PersonID: testPersonID,
			Username: "test123",
			Provider: "example.com",
			Primary:  true,
		},
	}
	testPersonPhones := []models.ContactPhone{
		{
			ID:          testPersonPhoneID,
			PersonID:    testPersonID,
			CountryCode: 1,
			PhoneNumber: "(315)559-1190",
			Type:        "home",
			Primary:     true,
		},
	}

	mockPersonStore := &mockPersonStore{}
	mockPersonStore.On("GetSpecificContactAddresses", mock.Anything, testPersonID).Return(
		testPersonAddresses, error(nil),
	)
	mockPersonStore.On("GetSpecificContactAddresses", mock.Anything, testBadPersonAddressID).Return(
		[]models.ContactAddress(nil), assert.AnError,
	)
	mockPersonStore.On("GetSpecificContactAddresses", mock.Anything, testBadPersonEmailID).Return(
		testPersonAddresses, error(nil),
	)
	mockPersonStore.On("GetSpecificContactAddresses", mock.Anything, testBadPersonPhoneID).Return(
		testPersonAddresses, error(nil),
	)

	mockPersonStore.On("GetSpecificContactEmails", mock.Anything, testPersonID).Return(
		testPersonEmails, error(nil),
	)
	mockPersonStore.On("GetSpecificContactEmails", mock.Anything, testBadPersonEmailID).Return(
		[]models.ContactEmail(nil), assert.AnError,
	)
	mockPersonStore.On("GetSpecificContactEmails", mock.Anything, testBadPersonPhoneID).Return(
		testPersonEmails, error(nil),
	)

	mockPersonStore.On("GetSpecificContactPhones", mock.Anything, testPersonID).Return(
		testPersonPhones, error(nil),
	)
	mockPersonStore.On("GetSpecificContactPhones", mock.Anything, testBadPersonPhoneID).Return(
		[]models.ContactPhone(nil), assert.AnError,
	)

	tests := []struct {
		name      string
		as        adminService
		args      args
		wants     wants
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			as: adminService{
				personStore: mockPersonStore,
			},
			args: args{
				ctx:      context.Background(),
				personID: testPersonID,
			},
			wants: wants{
				callAssertions: map[string]int{
					"GetSpecificContactAddresses": 1,
					"GetSpecificContactEmails":    1,
					"GetSpecificContactPhones":    1,
				},
				returnVal: models.Contacts{
					Addresses: testPersonAddresses,
					Email:     testPersonEmails,
					Phone:     testPersonPhones,
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; Get Addresses",
			as: adminService{
				personStore: mockPersonStore,
			},
			args: args{
				ctx:      context.Background(),
				personID: testBadPersonAddressID,
			},
			wants: wants{
				callAssertions: map[string]int{
					"GetSpecificContactAddresses": 1,
					"GetSpecificContactEmails":    0,
					"GetSpecificContactPhones":    0,
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Get Emails",
			as: adminService{
				personStore: mockPersonStore,
			},
			args: args{
				ctx:      context.Background(),
				personID: testBadPersonEmailID,
			},
			wants: wants{
				callAssertions: map[string]int{
					"GetSpecificContactAddresses": 1,
					"GetSpecificContactEmails":    1,
					"GetSpecificContactPhones":    0,
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Get Phone",
			as: adminService{
				personStore: mockPersonStore,
			},
			args: args{
				ctx:      context.Background(),
				personID: testBadPersonPhoneID,
			},
			wants: wants{
				callAssertions: map[string]int{
					"GetSpecificContactAddresses": 1,
					"GetSpecificContactEmails":    1,
					"GetSpecificContactPhones":    1,
				},
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.as.GetPersonContacts(tt.args.ctx, tt.args.personID)
			tt.assertion(t, err)
			assert.Equal(t, tt.wants.returnVal, got)

			// Since we are calling multiple functions from a mocked instance, we want to ensure that each mocked function
			// gets called as we expect it.
			for k, v := range tt.wants.callAssertions {
				mockPersonStore.AssertNumberOfCalls(t, k, v)
			}
		})
	}
}

func Test_adminService_UpdatePerson(t *testing.T) {
	type args struct {
		ctx  context.Context
		id   uuid.UUID
		data models.Person
	}

	testNotFoundID := uuid.New()
	testPersonID := uuid.New()
	testPerson := models.Person{
		Name: models.Name{
			GivenName:       "Testy",
			FamilyName:      "McTesterson",
			FamilyNameFirst: models.FirstNameGiven,
		},
		DateOfBirth: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		Gender:      models.GenderNonBinary,
		Pronouns:    models.Pronouns{Subject: "they", Object: "them"},
	}

	testPersonStore := &mockPersonStore{}
	testPersonStore.On("Update", mock.Anything, testPersonID, testPerson).Return(error(nil))
	testPersonStore.On("Update", mock.Anything, testNotFoundID, testPerson).Return(assert.AnError)

	tests := []struct {
		name      string
		as        adminService
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			as: adminService{
				personStore: testPersonStore,
			},
			args: args{
				ctx:  context.Background(),
				id:   testPersonID,
				data: testPerson,
			},
			assertion: assert.NoError,
		},
		{
			name: "Error",
			as: adminService{
				personStore: testPersonStore,
			},
			args: args{
				ctx:  context.Background(),
				id:   testNotFoundID,
				data: testPerson,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, tt.as.UpdatePerson(tt.args.ctx, tt.args.id, tt.args.data))
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
func (mps *mockPersonStore) GetSpecificContactAddresses(ctx context.Context, id uuid.UUID) ([]models.ContactAddress, error) {
	args := mps.Called(ctx, id)
	return args.Get(0).([]models.ContactAddress), args.Error(1)
}
func (mps *mockPersonStore) GetSpecificContactEmails(ctx context.Context, id uuid.UUID) ([]models.ContactEmail, error) {
	args := mps.Called(ctx, id)
	return args.Get(0).([]models.ContactEmail), args.Error(1)
}
func (mps *mockPersonStore) GetSpecificContactPhones(ctx context.Context, id uuid.UUID) ([]models.ContactPhone, error) {
	args := mps.Called(ctx, id)
	return args.Get(0).([]models.ContactPhone), args.Error(1)
}
