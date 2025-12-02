package endpoints

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/williabk198/timeclock/internal/models"
)

func Test_adminEmployeeEndpoints_Add(t *testing.T) {
	type args struct {
		ctx    context.Context
		person EmployeeData
	}

	testGoodPersonID := uuid.New()
	testGoodEmployeeID := uuid.New()
	testGoodEmployeeData := EmployeeData{
		PersonID:    testGoodPersonID.String(),
		ReportsToID: uuid.Nil.String(),
		Title:       "Owner",
	}
	testGoodEmployeeDB := models.Employee{
		PersonID:    testGoodPersonID,
		ReportsToID: uuid.Nil,
		Title:       testGoodEmployeeData.Title,
	}

	testBadEmployeeData := EmployeeData{
		PersonID:    uuid.NewString(),
		ReportsToID: uuid.NewString(),
		Title:       "error val",
	}
	testBadEmployeeDB := models.Employee{
		PersonID:    uuid.MustParse(testBadEmployeeData.PersonID),
		ReportsToID: uuid.MustParse(testBadEmployeeData.ReportsToID),
		Title:       testBadEmployeeData.Title,
	}

	testEmployeeMicro := &mockEmployeeMicro{}
	testEmployeeMicro.On("Add", mock.Anything, testGoodEmployeeDB).Return(testGoodEmployeeID, error(nil))
	testEmployeeMicro.On("Add", mock.Anything, testBadEmployeeDB).Return(uuid.Nil, assert.AnError)

	tests := []struct {
		name      string
		a         adminEmployeeEndpoints
		args      args
		want      EmployeeData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			a:    adminEmployeeEndpoints{employeeMicro: testEmployeeMicro},
			args: args{ctx: context.Background(), person: testGoodEmployeeData},
			want: EmployeeData{
				ID:          testGoodEmployeeID.String(),
				PersonID:    testGoodEmployeeData.PersonID,
				ReportsToID: testGoodEmployeeData.ReportsToID,
				Title:       testGoodEmployeeData.Title,
			},
			assertion: assert.NoError,
		},
		{
			name:      "Error; Invalid Input",
			a:         adminEmployeeEndpoints{employeeMicro: testEmployeeMicro},
			args:      args{ctx: context.Background(), person: EmployeeData{PersonID: "invalid_value"}},
			assertion: assert.Error,
		},
		{
			name:      "Error; Service Error",
			a:         adminEmployeeEndpoints{employeeMicro: testEmployeeMicro},
			args:      args{ctx: context.Background(), person: testBadEmployeeData},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.Add(tt.args.ctx, tt.args.person)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_adminEmployeeEndpoints_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}

	testNotFoundID := uuid.New()
	testEmployeeID := uuid.New()
	testEmployeeDB := models.Employee{
		PersonID:    uuid.New(),
		ReportsToID: uuid.Nil,
		Title:       "President",
	}
	testEmployeeData := EmployeeData{
		PersonID:    testEmployeeID.String(),
		ReportsToID: uuid.Nil.String(),
		Title:       "President",
	}

	testEmployeeMicro := &mockEmployeeMicro{}
	testEmployeeMicro.On("Delete", mock.Anything, testEmployeeID).Return(testEmployeeDB, error(nil))
	testEmployeeMicro.On("Delete", mock.Anything, testNotFoundID).Return(models.Employee{}, assert.AnError)

	tests := []struct {
		name      string
		a         adminEmployeeEndpoints
		args      args
		want      EmployeeData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "Success",
			a:         adminEmployeeEndpoints{employeeMicro: testEmployeeMicro},
			args:      args{ctx: context.Background(), id: testEmployeeID.String()},
			want:      testEmployeeData,
			assertion: assert.NoError,
		},
		{
			name:      "Error",
			a:         adminEmployeeEndpoints{employeeMicro: testEmployeeMicro},
			args:      args{ctx: context.Background(), id: testNotFoundID.String()},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.Delete(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_adminEmployeeEndpoints_GetAll(t *testing.T) {
	type args struct {
		ctx     context.Context
		reqData GetPaginatedRequestData
	}

	testEmployees := []models.Employee{
		{
			ID:          uuid.New(),
			PersonID:    uuid.New(),
			ReportsToID: uuid.New(),
			Title:       "Software Engineer",
		},
		{
			ID:          uuid.New(),
			PersonID:    uuid.New(),
			ReportsToID: uuid.New(),
			Title:       "Facilities Manager",
		},
		{
			ID:          uuid.New(),
			PersonID:    uuid.New(),
			ReportsToID: uuid.New(),
			Title:       "HR Manager",
		},
	}
	resultEmployees := []EmployeeData{
		{
			ID:          testEmployees[1].ID.String(),
			PersonID:    testEmployees[1].PersonID.String(),
			ReportsToID: testEmployees[1].ReportsToID.String(),
			Title:       testEmployees[1].Title,
		},
		{
			ID:          testEmployees[2].ID.String(),
			PersonID:    testEmployees[2].PersonID.String(),
			ReportsToID: testEmployees[2].ReportsToID.String(),
			Title:       testEmployees[2].Title,
		},
	}

	testEmployeeMicro := &mockEmployeeMicro{}
	testEmployeeMicro.On("GetAll", mock.Anything, uint(1), uint(2)).Return(testEmployees[1:3], error(nil))
	testEmployeeMicro.On("GetAll", mock.Anything, uint(0), uint(0)).Return([]models.Employee{}, assert.AnError)

	tests := []struct {
		name      string
		a         adminEmployeeEndpoints
		args      args
		want      []EmployeeData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "Success",
			a:         adminEmployeeEndpoints{employeeMicro: testEmployeeMicro},
			args:      args{ctx: context.Background(), reqData: GetPaginatedRequestData{Offset: 1, Limit: 2}},
			want:      resultEmployees,
			assertion: assert.NoError,
		},
		{
			name:      "Error",
			a:         adminEmployeeEndpoints{employeeMicro: testEmployeeMicro},
			args:      args{ctx: context.Background(), reqData: GetPaginatedRequestData{Offset: 0, Limit: 0}},
			want:      []EmployeeData{},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.GetAll(tt.args.ctx, tt.args.reqData)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_adminEmployeeEndpoints_GetSpecific(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}

	testNotFoundID := uuid.New()
	testEmployeeID := uuid.New()
	testPersonID := uuid.New()
	testEmployeeDB := models.Employee{
		ID:          testEmployeeID,
		PersonID:    testPersonID,
		ReportsToID: uuid.Nil,
		Title:       "President",
	}
	testEmployeeData := EmployeeData{
		ID:          testEmployeeID.String(),
		PersonID:    testPersonID.String(),
		ReportsToID: uuid.Nil.String(),
		Title:       "President",
	}

	testEmployeeMicro := &mockEmployeeMicro{}
	testEmployeeMicro.On("GetSpecific", mock.Anything, testEmployeeID).Return(testEmployeeDB, error(nil))
	testEmployeeMicro.On("GetSpecific", mock.Anything, testNotFoundID).Return(models.Employee{}, assert.AnError)

	tests := []struct {
		name      string
		a         adminEmployeeEndpoints
		args      args
		want      EmployeeData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "Success",
			a:         adminEmployeeEndpoints{employeeMicro: testEmployeeMicro},
			args:      args{ctx: context.Background(), id: testEmployeeID.String()},
			want:      testEmployeeData,
			assertion: assert.NoError,
		},
		{
			name:      "Error",
			a:         adminEmployeeEndpoints{employeeMicro: testEmployeeMicro},
			args:      args{ctx: context.Background(), id: testNotFoundID.String()},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.GetSpecific(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_adminEmployeeEndpoints_Update(t *testing.T) {
	type args struct {
		ctx context.Context
		urd UpdateRequestData[EmployeeData]
	}

	testNotFoundID := uuid.New()
	testEmployeeID := uuid.New()
	testEmployeeDB := models.Employee{
		ReportsToID: uuid.Nil,
		Title:       "Co-Owner",
	}
	testEmployeeData := EmployeeData{
		ReportsToID: testEmployeeDB.ReportsToID.String(),
		Title:       testEmployeeDB.Title,
	}

	testEmployeeMicro := &mockEmployeeMicro{}
	testEmployeeMicro.On("Update", mock.Anything, testEmployeeID, testEmployeeDB).Return(error(nil))
	testEmployeeMicro.On("Update", mock.Anything, testNotFoundID, testEmployeeDB).Return(assert.AnError)

	tests := []struct {
		name      string
		a         adminEmployeeEndpoints
		args      args
		want      EmployeeData
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			a:    adminEmployeeEndpoints{employeeMicro: testEmployeeMicro},
			args: args{
				ctx: context.Background(),
				urd: UpdateRequestData[EmployeeData]{
					ID:   testEmployeeID.String(),
					Data: testEmployeeData,
				},
			},
			want:      testEmployeeData,
			assertion: assert.NoError,
		},
		{
			name: "Error",
			a:    adminEmployeeEndpoints{employeeMicro: testEmployeeMicro},
			args: args{
				ctx: context.Background(),
				urd: UpdateRequestData[EmployeeData]{
					ID:   testNotFoundID.String(),
					Data: testEmployeeData,
				},
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.Update(tt.args.ctx, tt.args.urd)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
