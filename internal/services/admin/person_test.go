package admin

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/williabk198/timeclock/internal/models"
)

func Test_personMicroImpl_Add(t *testing.T) {
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
		a         personMicroImpl
		args      args
		wants     wants
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			a: personMicroImpl{
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
			a: personMicroImpl{
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
			gotID, err := tt.a.Add(tt.args.ctx, tt.args.person)
			tt.assertion(t, err)
			assert.Equal(t, tt.wants.id, gotID)
		})
	}
}

func Test_personMicroImpl_Delete(t *testing.T) {
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
		pm        PersonMicro
		args      args
		want      models.Person
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			pm: personMicroImpl{
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
			pm: personMicroImpl{
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
			got, err := tt.pm.Delete(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_personMicroImpl_GetAll(t *testing.T) {
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
		as        PersonMicro
		args      args
		want      []models.Person
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			as: personMicroImpl{
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
			as: personMicroImpl{
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
			got, err := tt.as.GetAll(tt.args.ctx, tt.args.offset, tt.args.limit)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_personMicroImpl_GetSpecific(t *testing.T) {
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
		as        PersonMicro
		args      args
		want      models.Person
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			as: personMicroImpl{
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
			as: personMicroImpl{
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
			got, err := tt.as.GetSpecific(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_personMicroImpl_Update(t *testing.T) {
	type args struct {
		ctx    context.Context
		id     uuid.UUID
		newVal models.Person
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
		pm        PersonMicro
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			pm: personMicroImpl{
				personStore: testPersonStore,
			},
			args: args{
				ctx:    context.Background(),
				id:     testPersonID,
				newVal: testPerson,
			},
			assertion: assert.NoError,
		},
		{
			name: "Error",
			pm: personMicroImpl{
				personStore: testPersonStore,
			},
			args: args{
				ctx:    context.Background(),
				id:     testNotFoundID,
				newVal: testPerson,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, tt.pm.Update(tt.args.ctx, tt.args.id, tt.args.newVal))
		})
	}
}
