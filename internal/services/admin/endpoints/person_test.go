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
	testGoodPersonDB := &models.Person{
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
	testBadPersonDB := &models.Person{
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

type mockAdminService struct {
	mock.Mock
}

func (mas *mockAdminService) AddPerson(ctx context.Context, person models.Person) (uuid.UUID, error) {
	args := mas.Called(ctx, person)
	return args.Get(0).(uuid.UUID), args.Error(1)
}
