package admin

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/williabk198/timeclock/internal/models"
)

func Test_contactMicroImpl_AddPersonAddress(t *testing.T) {
	type args struct {
		ctx     context.Context
		address models.ContactAddress
	}
	testAddressID := uuid.New()
	testPersonAddress := models.ContactAddress{
		PersonID:   uuid.New(),
		Street1:    "123 Test Dr",
		Street2:    "",
		Locality:   "Testerville",
		Region:     "Testaria",
		PostalCode: "12345-6789",
		Country:    "Testopia",
		Type:       models.AddressTypePhysical,
		Primary:    true,
	}
	testErrorPersonAddress := models.ContactAddress{
		PersonID: uuid.New(),
		Street1:  "123 Test Dr",
		Street2:  "erronious_val",
	}

	testContactStore := &mockContactStore{}
	testContactStore.On("AddPersonAddress", mock.Anything, testPersonAddress).Return(testAddressID, error(nil))
	testContactStore.On("AddPersonAddress", mock.Anything, testErrorPersonAddress).Return(uuid.Nil, assert.AnError)

	tests := []struct {
		name      string
		as        contactMicroImpl
		args      args
		want      uuid.UUID
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			as: contactMicroImpl{
				contactStore: testContactStore,
			},
			args: args{
				ctx:     context.Background(),
				address: testPersonAddress,
			},
			want:      testAddressID,
			assertion: assert.NoError,
		},
		{
			name: "Error",
			as: contactMicroImpl{
				contactStore: testContactStore,
			},
			args: args{
				ctx:     context.Background(),
				address: testErrorPersonAddress,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.as.AddPersonAddress(tt.args.ctx, tt.args.address)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_contactMicroImpl_AddPersonEmail(t *testing.T) {
	type args struct {
		ctx   context.Context
		email models.ContactEmail
	}
	type wants struct {
		id uuid.UUID
	}

	testEmailID := uuid.New()
	testPersonEmail := models.ContactEmail{
		PersonID: uuid.New(),
		Username: "test",
		Provider: "example.com",
		Primary:  true,
	}
	testErrorPersonEmail := models.ContactEmail{
		PersonID: uuid.New(),
		Username: "test",
		Provider: "invalid",
	}

	testContactStore := &mockContactStore{}
	testContactStore.On("AddPersonEmail", mock.Anything, testPersonEmail).Return(testEmailID, error(nil))
	testContactStore.On("AddPersonEmail", mock.Anything, testErrorPersonEmail).Return(uuid.Nil, assert.AnError)

	tests := []struct {
		name      string
		cm        ContactMicro
		args      args
		wants     wants
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			cm: contactMicroImpl{
				contactStore: testContactStore,
			},
			args: args{
				ctx:   context.Background(),
				email: testPersonEmail,
			},
			wants:     wants{id: testEmailID},
			assertion: assert.NoError,
		},
		{
			name: "Error",
			cm: contactMicroImpl{
				contactStore: testContactStore,
			},
			args: args{
				ctx:   context.Background(),
				email: testErrorPersonEmail,
			},
			wants:     wants{id: uuid.Nil},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cm.AddPersonEmail(tt.args.ctx, tt.args.email)
			tt.assertion(t, err)
			assert.Equal(t, tt.wants.id, got)
		})
	}
}

func Test_contactMicroImpl_AddPersonPhone(t *testing.T) {
	type args struct {
		ctx   context.Context
		phone models.ContactPhone
	}
	testPhoneID := uuid.New()
	testPersonPhone := models.ContactPhone{
		PersonID:    uuid.New(),
		CountryCode: 1,
		PhoneNumber: "555 555 5555",
		Type:        models.PhoneTypeHome,
		Primary:     true,
	}
	testErrorPersonPhone := models.ContactPhone{
		PersonID:    uuid.New(),
		CountryCode: -77,
	}

	testContactStore := &mockContactStore{}
	testContactStore.On("AddPersonPhone", mock.Anything, testPersonPhone).Return(testPhoneID, error(nil))
	testContactStore.On("AddPersonPhone", mock.Anything, testErrorPersonPhone).Return(uuid.Nil, assert.AnError)

	tests := []struct {
		name      string
		cm        ContactMicro
		args      args
		want      uuid.UUID
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			cm: contactMicroImpl{
				contactStore: testContactStore,
			},
			args: args{
				ctx:   context.Background(),
				phone: testPersonPhone,
			},
			want:      testPhoneID,
			assertion: assert.NoError,
		},
		{
			name: "Error",
			cm: contactMicroImpl{
				contactStore: testContactStore,
			},
			args: args{
				ctx:   context.Background(),
				phone: testErrorPersonPhone,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cm.AddPersonPhone(tt.args.ctx, tt.args.phone)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_contactMicroImpl_DeletePerosnAddress(t *testing.T) {
	type args struct {
		ctx       context.Context
		personID  uuid.UUID
		addressID uuid.UUID
	}

	testValidAddressID := uuid.New()
	testNotFoundAddressID := uuid.New()

	testValidPersonID := uuid.New()
	testNotFoundPersonID := uuid.New()

	testAddress := models.ContactAddress{
		ID:         testValidAddressID,
		PersonID:   testValidPersonID,
		Street1:    "123 Test Dr",
		Street2:    "",
		Locality:   "Testerville",
		Region:     "Testaria",
		PostalCode: "12345-6789",
		Country:    "Testopia",
		Type:       models.AddressTypePhysical,
		Primary:    true,
	}

	testContactStore := &mockContactStore{}
	testContactStore.On("DeletePersonAddress", mock.Anything, testValidPersonID, testValidAddressID).Return(testAddress, error(nil))
	testContactStore.On("DeletePersonAddress", mock.Anything, testNotFoundPersonID, testValidAddressID).Return(models.ContactAddress{}, assert.AnError)
	testContactStore.On("DeletePersonAddress", mock.Anything, testValidPersonID, testNotFoundAddressID).Return(models.ContactAddress{}, assert.AnError)

	tests := []struct {
		name      string
		cmi       contactMicroImpl
		args      args
		want      models.ContactAddress
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			cmi: contactMicroImpl{
				contactStore: testContactStore,
			},
			args: args{
				ctx:       context.Background(),
				personID:  testValidPersonID,
				addressID: testValidAddressID,
			},
			want:      testAddress,
			assertion: assert.NoError,
		},
		{
			name: "Error; Person DNE",
			cmi: contactMicroImpl{
				contactStore: testContactStore,
			},
			args: args{
				ctx:       context.Background(),
				personID:  testNotFoundPersonID,
				addressID: testValidAddressID,
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Address DNE",
			cmi: contactMicroImpl{
				contactStore: testContactStore,
			},
			args: args{
				ctx:       context.Background(),
				personID:  testValidPersonID,
				addressID: testNotFoundAddressID,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cmi.DeletePerosnAddress(tt.args.ctx, tt.args.personID, tt.args.addressID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_contactMicroImpl_DeletePersonEmail(t *testing.T) {
	type args struct {
		ctx      context.Context
		personID uuid.UUID
		emailID  uuid.UUID
	}

	testValidEmailID := uuid.New()
	testNotFoundEmailID := uuid.New()

	testValidPersonID := uuid.New()
	testNotFoundPersonID := uuid.New()

	testEmail := models.ContactEmail{
		ID:       testValidEmailID,
		PersonID: testValidPersonID,
		Username: "tester",
		Provider: "test.com",
		Primary:  true,
	}

	testContactStore := &mockContactStore{}
	testContactStore.On("DeletePersonEmail", mock.Anything, testValidPersonID, testValidEmailID).Return(testEmail, error(nil))
	testContactStore.On("DeletePersonEmail", mock.Anything, testNotFoundPersonID, testValidEmailID).Return(models.ContactEmail{}, assert.AnError)
	testContactStore.On("DeletePersonEmail", mock.Anything, testValidPersonID, testNotFoundEmailID).Return(models.ContactEmail{}, assert.AnError)

	tests := []struct {
		name      string
		cmi       contactMicroImpl
		args      args
		want      models.ContactEmail
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			cmi: contactMicroImpl{
				contactStore: testContactStore,
			},
			args: args{
				ctx:      context.Background(),
				personID: testValidPersonID,
				emailID:  testValidEmailID,
			},
			want:      testEmail,
			assertion: assert.NoError,
		},
		{
			name: "Error; Person DNE",
			cmi: contactMicroImpl{
				contactStore: testContactStore,
			},
			args: args{
				ctx:      context.Background(),
				personID: testNotFoundPersonID,
				emailID:  testValidEmailID,
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Address DNE",
			cmi: contactMicroImpl{
				contactStore: testContactStore,
			},
			args: args{
				ctx:      context.Background(),
				personID: testValidPersonID,
				emailID:  testNotFoundEmailID,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cmi.DeletePersonEmail(tt.args.ctx, tt.args.personID, tt.args.emailID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_contactMicroImpl_DeletePersonPhone(t *testing.T) {
	type args struct {
		ctx      context.Context
		personID uuid.UUID
		phoneID  uuid.UUID
	}

	testValidPhoneID := uuid.New()
	testNotFoundPhoneID := uuid.New()

	testValidPersonID := uuid.New()
	testNotFoundPersonID := uuid.New()

	testPhone := models.ContactPhone{
		ID:          testValidPhoneID,
		PersonID:    testValidPersonID,
		CountryCode: 1,
		PhoneNumber: "(555)555-5555",
		Type:        models.PhoneTypeHome,
		Primary:     true,
	}

	testContactStore := &mockContactStore{}
	testContactStore.On("DeletePersonPhone", mock.Anything, testValidPersonID, testValidPhoneID).Return(testPhone, error(nil))
	testContactStore.On("DeletePersonPhone", mock.Anything, testNotFoundPersonID, testValidPhoneID).Return(models.ContactPhone{}, assert.AnError)
	testContactStore.On("DeletePersonPhone", mock.Anything, testValidPersonID, testNotFoundPhoneID).Return(models.ContactPhone{}, assert.AnError)

	tests := []struct {
		name      string
		cmi       contactMicroImpl
		args      args
		want      models.ContactPhone
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			cmi: contactMicroImpl{
				contactStore: testContactStore,
			},
			args: args{
				ctx:      context.Background(),
				personID: testValidPersonID,
				phoneID:  testValidPhoneID,
			},
			want:      testPhone,
			assertion: assert.NoError,
		},
		{
			name: "Error; Person DNE",
			cmi: contactMicroImpl{
				contactStore: testContactStore,
			},
			args: args{
				ctx:      context.Background(),
				personID: testNotFoundPersonID,
				phoneID:  testValidPhoneID,
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Address DNE",
			cmi: contactMicroImpl{
				contactStore: testContactStore,
			},
			args: args{
				ctx:      context.Background(),
				personID: testValidPersonID,
				phoneID:  testNotFoundPhoneID,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cmi.DeletePersonPhone(tt.args.ctx, tt.args.personID, tt.args.phoneID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_contactMicroImpl_GetAllForPerson(t *testing.T) {
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
			PhoneNumber: "(555)555-5555",
			Type:        "home",
			Primary:     true,
		},
	}

	mockContactStore := &mockContactStore{}
	mockContactStore.On("GetPersonAddresses", mock.Anything, testPersonID).Return(
		testPersonAddresses, error(nil),
	)
	mockContactStore.On("GetPersonAddresses", mock.Anything, testBadPersonAddressID).Return(
		[]models.ContactAddress(nil), assert.AnError,
	)
	mockContactStore.On("GetPersonAddresses", mock.Anything, testBadPersonEmailID).Return(
		testPersonAddresses, error(nil),
	)
	mockContactStore.On("GetPersonAddresses", mock.Anything, testBadPersonPhoneID).Return(
		testPersonAddresses, error(nil),
	)

	mockContactStore.On("GetPersonEmails", mock.Anything, testPersonID).Return(
		testPersonEmails, error(nil),
	)
	mockContactStore.On("GetPersonEmails", mock.Anything, testBadPersonAddressID).Return(
		testPersonEmails, error(nil),
	)
	mockContactStore.On("GetPersonEmails", mock.Anything, testBadPersonEmailID).Return(
		[]models.ContactEmail(nil), assert.AnError,
	)
	mockContactStore.On("GetPersonEmails", mock.Anything, testBadPersonPhoneID).Return(
		testPersonEmails, error(nil),
	)

	mockContactStore.On("GetPersonPhones", mock.Anything, testPersonID).Return(
		testPersonPhones, error(nil),
	)
	mockContactStore.On("GetPersonPhones", mock.Anything, testBadPersonAddressID).Return(
		[]models.ContactPhone(nil), assert.AnError,
	)
	mockContactStore.On("GetPersonPhones", mock.Anything, testBadPersonEmailID).Return(
		testPersonPhones, error(nil),
	)
	mockContactStore.On("GetPersonPhones", mock.Anything, testBadPersonPhoneID).Return(
		[]models.ContactPhone(nil), assert.AnError,
	)

	tests := []struct {
		name      string
		cm        ContactMicro
		args      args
		wants     wants
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			cm: contactMicroImpl{
				contactStore: mockContactStore,
			},
			args: args{
				ctx:      context.Background(),
				personID: testPersonID,
			},
			wants: wants{
				callAssertions: map[string]int{
					"GetPersonAddresses": 1,
					"GetPersonEmails":    1,
					"GetPersonPhones":    1,
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
			cm: contactMicroImpl{
				contactStore: mockContactStore,
			},
			args: args{
				ctx:      context.Background(),
				personID: testBadPersonAddressID,
			},
			wants: wants{
				callAssertions: map[string]int{ // TODO: look into better ways for handling call counts that doesn't accumulate with each test case
					"GetPersonAddresses": 2,
					"GetPersonEmails":    2,
					"GetPersonPhones":    2,
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Get Emails",
			cm: contactMicroImpl{
				contactStore: mockContactStore,
			},
			args: args{
				ctx:      context.Background(),
				personID: testBadPersonEmailID,
			},
			wants: wants{
				callAssertions: map[string]int{
					"GetPersonAddresses": 3,
					"GetPersonEmails":    3,
					"GetPersonPhones":    3,
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Get Phone",
			cm: contactMicroImpl{
				contactStore: mockContactStore,
			},
			args: args{
				ctx:      context.Background(),
				personID: testBadPersonPhoneID,
			},
			wants: wants{
				callAssertions: map[string]int{
					"GetPersonAddresses": 4,
					"GetPersonEmails":    4,
					"GetPersonPhones":    4,
				},
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cm.GetAllForPerson(tt.args.ctx, tt.args.personID)
			tt.assertion(t, err)
			assert.Equal(t, tt.wants.returnVal, got)

			// Since we are calling multiple functions from a mocked instance, we want to ensure that each mocked function
			// gets called as we expect it.
			for k, v := range tt.wants.callAssertions {
				mockContactStore.AssertNumberOfCalls(t, k, v)
			}
		})
	}
}

func Test_contactMicroImpl_GetPersonAddresses(t *testing.T) {
	type args struct {
		ctx      context.Context
		personID uuid.UUID
	}
	testPersonID := uuid.New()
	testAddressID := uuid.New()
	testDoesNotExistID := uuid.New()

	testPersonAddresses := []models.ContactAddress{
		{
			ID:         testAddressID,
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

	mockContactStore := &mockContactStore{}
	mockContactStore.On("GetPersonAddresses", mock.Anything, testPersonID).Return(
		testPersonAddresses, error(nil),
	)
	mockContactStore.On("GetPersonAddresses", mock.Anything, testDoesNotExistID).Return(
		[]models.ContactAddress(nil), assert.AnError,
	)

	tests := []struct {
		name      string
		cm        ContactMicro
		args      args
		want      []models.ContactAddress
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			cm: contactMicroImpl{
				contactStore: mockContactStore,
			},
			args: args{
				ctx:      context.Background(),
				personID: testPersonID,
			},
			want:      testPersonAddresses,
			assertion: assert.NoError,
		},
		{
			name: "Error; Invalid ID",
			cm: contactMicroImpl{
				contactStore: mockContactStore,
			},
			args: args{
				ctx:      context.Background(),
				personID: testDoesNotExistID,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cm.GetPersonAddresses(tt.args.ctx, tt.args.personID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_contactMicroImpl_GetPersonEmails(t *testing.T) {
	type args struct {
		ctx      context.Context
		personID uuid.UUID
	}
	testPersonID := uuid.New()
	testEmailID := uuid.New()
	testDoesNotExistID := uuid.New()
	testPersonEmails := []models.ContactEmail{
		{
			ID:       testEmailID,
			PersonID: testPersonID,
			Username: "test123",
			Provider: "example.com",
			Primary:  true,
		},
	}

	mockContactStore := &mockContactStore{}
	mockContactStore.On("GetPersonEmails", mock.Anything, testPersonID).Return(
		testPersonEmails, error(nil),
	)
	mockContactStore.On("GetPersonEmails", mock.Anything, testDoesNotExistID).Return(
		[]models.ContactEmail(nil), assert.AnError,
	)

	tests := []struct {
		name      string
		cm        ContactMicro
		args      args
		want      []models.ContactEmail
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			cm: contactMicroImpl{
				contactStore: mockContactStore,
			},
			args: args{
				ctx:      context.Background(),
				personID: testPersonID,
			},
			want:      testPersonEmails,
			assertion: assert.NoError,
		},
		{
			name: "Error; Invalid ID",
			cm: contactMicroImpl{
				contactStore: mockContactStore,
			},
			args: args{
				ctx:      context.Background(),
				personID: testDoesNotExistID,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cm.GetPersonEmails(tt.args.ctx, tt.args.personID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_contactMicroImpl_GetPersonContactPhones(t *testing.T) {
	type args struct {
		ctx      context.Context
		personID uuid.UUID
	}
	testPersonID := uuid.New()
	testPhoneID := uuid.New()
	testDoesNotExistID := uuid.New()
	testPersonPhones := []models.ContactPhone{
		{
			ID:          testPhoneID,
			PersonID:    testPersonID,
			CountryCode: 1,
			PhoneNumber: "(555)555-5555",
			Type:        "home",
			Primary:     true,
		},
	}
	mockContactStore := &mockContactStore{}
	mockContactStore.On("GetPersonPhones", mock.Anything, testPersonID).Return(
		testPersonPhones, error(nil),
	)
	mockContactStore.On("GetPersonPhones", mock.Anything, testDoesNotExistID).Return(
		[]models.ContactPhone(nil), assert.AnError,
	)

	tests := []struct {
		name      string
		cm        ContactMicro
		args      args
		want      []models.ContactPhone
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			cm: contactMicroImpl{
				contactStore: mockContactStore,
			},
			args: args{
				ctx:      context.Background(),
				personID: testPersonID,
			},
			want:      testPersonPhones,
			assertion: assert.NoError,
		},
		{
			name: "Error; Invalid ID",
			cm: contactMicroImpl{
				contactStore: mockContactStore,
			},
			args: args{
				ctx:      context.Background(),
				personID: testDoesNotExistID,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cm.GetPersonPhones(tt.args.ctx, tt.args.personID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_contactMicroImpl_UpdatePersonAddress(t *testing.T) {
	type args struct {
		ctx       context.Context
		personID  uuid.UUID
		addressID uuid.UUID
		newVal    models.ContactAddress
	}
	testPersonID := uuid.New()
	testNotFoundPersonID := uuid.New()

	testAddressID := uuid.New()
	testNotFoundAddressID := uuid.New()

	testNewAddresValue := models.ContactAddress{
		Street1:    "123 Test Dr",
		Locality:   "Testville",
		Region:     "Testeria",
		PostalCode: "12345-6789",
		Country:    "Testopia",
		Type:       models.AddressTypeMailing,
		Primary:    true,
	}

	testContactStore := &mockContactStore{}
	testContactStore.On("UpdatePersonAddress", mock.Anything, testPersonID, testAddressID, testNewAddresValue).Return(error(nil))
	testContactStore.On("UpdatePersonAddress", mock.Anything, testNotFoundPersonID, testAddressID, testNewAddresValue).Return(assert.AnError)
	testContactStore.On("UpdatePersonAddress", mock.Anything, testPersonID, testNotFoundAddressID, testNewAddresValue).Return(assert.AnError)

	tests := []struct {
		name      string
		cmi       contactMicroImpl
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			cmi: contactMicroImpl{
				contactStore: testContactStore,
			},
			args: args{
				ctx:       context.Background(),
				personID:  testPersonID,
				addressID: testAddressID,
				newVal:    testNewAddresValue,
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; Person DNE",
			cmi: contactMicroImpl{
				contactStore: testContactStore,
			},
			args: args{
				ctx:       context.Background(),
				personID:  testNotFoundPersonID,
				addressID: testAddressID,
				newVal:    testNewAddresValue,
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Address DNE",
			cmi: contactMicroImpl{
				contactStore: testContactStore,
			},
			args: args{
				ctx:       context.Background(),
				personID:  testPersonID,
				addressID: testNotFoundAddressID,
				newVal:    testNewAddresValue,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, tt.cmi.UpdatePersonAddress(tt.args.ctx, tt.args.personID, tt.args.addressID, tt.args.newVal))
		})
	}
}

func Test_contactMicroImpl_UpdatePersonEmail(t *testing.T) {
	type args struct {
		ctx      context.Context
		personID uuid.UUID
		emailID  uuid.UUID
		newVal   models.ContactEmail
	}
	testPersonID := uuid.New()
	testNotFoundPersonID := uuid.New()

	testEmailID := uuid.New()
	testNotFoundEmailID := uuid.New()

	testNewEmailValue := models.ContactEmail{
		Username: "tester",
		Provider: "test.com",
		Primary:  true,
	}

	testContactStore := &mockContactStore{}
	testContactStore.On("UpdatePersonEmail", mock.Anything, testPersonID, testEmailID, testNewEmailValue).Return(error(nil))
	testContactStore.On("UpdatePersonEmail", mock.Anything, testNotFoundPersonID, testEmailID, testNewEmailValue).Return(assert.AnError)
	testContactStore.On("UpdatePersonEmail", mock.Anything, testPersonID, testNotFoundEmailID, testNewEmailValue).Return(assert.AnError)

	tests := []struct {
		name      string
		cmi       contactMicroImpl
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			cmi: contactMicroImpl{
				contactStore: testContactStore,
			},
			args: args{
				ctx:      context.Background(),
				personID: testPersonID,
				emailID:  testEmailID,
				newVal:   testNewEmailValue,
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; Person DNE",
			cmi: contactMicroImpl{
				contactStore: testContactStore,
			},
			args: args{
				ctx:      context.Background(),
				personID: testNotFoundPersonID,
				emailID:  testEmailID,
				newVal:   testNewEmailValue,
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Email DNE",
			cmi: contactMicroImpl{
				contactStore: testContactStore,
			},
			args: args{
				ctx:      context.Background(),
				personID: testPersonID,
				emailID:  testNotFoundEmailID,
				newVal:   testNewEmailValue,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, tt.cmi.UpdatePersonEmail(tt.args.ctx, tt.args.personID, tt.args.emailID, tt.args.newVal))
		})
	}
}

func Test_contactMicroImpl_UpdatePersonPhone(t *testing.T) {
	type args struct {
		ctx      context.Context
		personID uuid.UUID
		phoneID  uuid.UUID
		newVal   models.ContactPhone
	}
	testPersonID := uuid.New()
	testNotFoundPersonID := uuid.New()

	testPhoneID := uuid.New()
	testNotFoundPhoneID := uuid.New()

	testNewPhoneValue := models.ContactPhone{
		CountryCode: 1,
		PhoneNumber: "(555)555-5555",
		Type:        models.PhoneTypeHome,
		Primary:     true,
	}

	testContactStore := &mockContactStore{}
	testContactStore.On("UpdatePersonPhone", mock.Anything, testPersonID, testPhoneID, testNewPhoneValue).Return(error(nil))
	testContactStore.On("UpdatePersonPhone", mock.Anything, testNotFoundPersonID, testPhoneID, testNewPhoneValue).Return(assert.AnError)
	testContactStore.On("UpdatePersonPhone", mock.Anything, testPersonID, testNotFoundPhoneID, testNewPhoneValue).Return(assert.AnError)

	tests := []struct {
		name      string
		cmi       contactMicroImpl
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			cmi: contactMicroImpl{
				contactStore: testContactStore,
			},
			args: args{
				ctx:      context.Background(),
				personID: testPersonID,
				phoneID:  testPhoneID,
				newVal:   testNewPhoneValue,
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; Person DNE",
			cmi: contactMicroImpl{
				contactStore: testContactStore,
			},
			args: args{
				ctx:      context.Background(),
				personID: testNotFoundPersonID,
				phoneID:  testPhoneID,
				newVal:   testNewPhoneValue,
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Phone DNE",
			cmi: contactMicroImpl{
				contactStore: testContactStore,
			},
			args: args{
				ctx:      context.Background(),
				personID: testPersonID,
				phoneID:  testNotFoundPhoneID,
				newVal:   testNewPhoneValue,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, tt.cmi.UpdatePersonPhone(tt.args.ctx, tt.args.personID, tt.args.phoneID, tt.args.newVal))
		})
	}
}
