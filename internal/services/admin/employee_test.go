package admin

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/williabk198/timeclock/internal/models"
)

func Test_employeeMicroImpl_Add(t *testing.T) {
	type args struct {
		ctx    context.Context
		person models.Employee
	}
	type wants struct {
		id uuid.UUID
	}

	testEmployeeID := uuid.New()
	testPersonID := uuid.New()

	testEmployee := models.Employee{
		ID:          testEmployeeID,
		PersonID:    testPersonID,
		ReportsToID: uuid.Nil,
		Title:       "Owner/President",
	}
	testEmployeeError := models.Employee{
		Title: "error_val",
	}

	testEmployeeStore := &mockEmployeeStore{}
	testEmployeeStore.On("Add", mock.Anything, testEmployee).Return(testEmployeeID, error(nil))
	testEmployeeStore.On("Add", mock.Anything, testEmployeeError).Return(uuid.Nil, assert.AnError)

	tests := []struct {
		name      string
		e         employeeMicroImpl
		args      args
		wants     wants
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "Success",
			e:         employeeMicroImpl{employeeStore: testEmployeeStore},
			args:      args{ctx: context.Background(), person: testEmployee},
			wants:     wants{id: testEmployeeID},
			assertion: assert.NoError,
		},
		{
			name:      "Error",
			e:         employeeMicroImpl{employeeStore: testEmployeeStore},
			args:      args{ctx: context.Background(), person: testEmployeeError},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, err := tt.e.Add(tt.args.ctx, tt.args.person)
			tt.assertion(t, err)
			assert.Equal(t, tt.wants.id, gotID)
		})
	}
}

func Test_employeeMicroImpl_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}

	testEmployeNotFoundID := uuid.New()
	testEmployeeID := uuid.New()
	testEmployee := models.Employee{
		ID:          testEmployeeID,
		PersonID:    uuid.New(),
		ReportsToID: uuid.Nil,
		Title:       "Owner/President",
	}

	testEmployeeStore := &mockEmployeeStore{}
	testEmployeeStore.On("Delete", mock.Anything, testEmployeeID).Return(testEmployee, error(nil))
	testEmployeeStore.On("Delete", mock.Anything, testEmployeNotFoundID).Return(models.Employee{}, assert.AnError)

	tests := []struct {
		name      string
		e         employeeMicroImpl
		args      args
		want      models.Employee
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "Success",
			e:         employeeMicroImpl{employeeStore: testEmployeeStore},
			args:      args{ctx: context.Background(), id: testEmployeeID},
			want:      testEmployee,
			assertion: assert.NoError,
		},
		{
			name:      "Error",
			e:         employeeMicroImpl{employeeStore: testEmployeeStore},
			args:      args{ctx: context.Background(), id: testEmployeNotFoundID},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.Delete(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_employeeMicroImpl_GetAll(t *testing.T) {
	type args struct {
		ctx    context.Context
		offset uint
		limit  uint
	}

	ownerEmployeeID := uuid.New()
	testEmployees := []models.Employee{
		{
			ID:          ownerEmployeeID,
			PersonID:    uuid.New(),
			ReportsToID: uuid.Nil,
			Title:       "Owner",
		},
		{
			ID:          uuid.New(),
			PersonID:    uuid.New(),
			ReportsToID: ownerEmployeeID,
			Title:       "HR Manager",
		},
		{
			ID:          uuid.New(),
			PersonID:    uuid.New(),
			ReportsToID: ownerEmployeeID,
			Title:       "Director of Technology",
		},
	}

	testEmployeeStore := &mockEmployeeStore{}
	testEmployeeStore.On("GetAllPaginated", mock.Anything, uint(0), uint(2)).Return(testEmployees[:2], error(nil))
	testEmployeeStore.On("GetAllPaginated", mock.Anything, uint(0), uint(0)).Return([]models.Employee(nil), assert.AnError)

	tests := []struct {
		name      string
		e         employeeMicroImpl
		args      args
		want      []models.Employee
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			e: employeeMicroImpl{
				employeeStore: testEmployeeStore,
			},
			args:      args{ctx: context.Background(), offset: 0, limit: 2},
			want:      testEmployees[:2],
			assertion: assert.NoError,
		},
		{
			name: "Error",
			e: employeeMicroImpl{
				employeeStore: testEmployeeStore,
			},
			args:      args{ctx: context.Background(), offset: 0, limit: 0},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.GetAll(tt.args.ctx, tt.args.offset, tt.args.limit)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_employeeMicroImpl_GetSpecific(t *testing.T) {
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}

	testEmployeNotFoundID := uuid.New()
	testEmployeeID := uuid.New()
	testEmployee := models.Employee{
		ID:          testEmployeeID,
		PersonID:    uuid.New(),
		ReportsToID: uuid.Nil,
		Title:       "Owner/President",
	}

	testEmployeeStore := &mockEmployeeStore{}
	testEmployeeStore.On("GetSpecific", mock.Anything, testEmployeeID).Return(testEmployee, error(nil))
	testEmployeeStore.On("GetSpecific", mock.Anything, testEmployeNotFoundID).Return(models.Employee{}, assert.AnError)

	tests := []struct {
		name      string
		e         employeeMicroImpl
		args      args
		want      models.Employee
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "Success",
			e:         employeeMicroImpl{employeeStore: testEmployeeStore},
			args:      args{ctx: context.Background(), id: testEmployeeID},
			want:      testEmployee,
			assertion: assert.NoError,
		},
		{
			name:      "Error",
			e:         employeeMicroImpl{employeeStore: testEmployeeStore},
			args:      args{ctx: context.Background(), id: testEmployeNotFoundID},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.GetSpecific(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_employeeMicroImpl_Update(t *testing.T) {
	type args struct {
		ctx    context.Context
		id     uuid.UUID
		newVal models.Employee
	}

	testEmployeNotFoundID := uuid.New()
	testEmployeeID := uuid.New()
	testEmployee := models.Employee{
		ID:          testEmployeeID,
		PersonID:    uuid.New(),
		ReportsToID: uuid.Nil,
		Title:       "CEO",
	}

	testEmployeeStore := &mockEmployeeStore{}
	testEmployeeStore.On("Update", mock.Anything, testEmployeeID, testEmployee).Return(error(nil))
	testEmployeeStore.On("Update", mock.Anything, testEmployeNotFoundID, testEmployee).Return(assert.AnError)

	tests := []struct {
		name      string
		e         employeeMicroImpl
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "Success",
			e:         employeeMicroImpl{employeeStore: testEmployeeStore},
			args:      args{ctx: context.Background(), id: testEmployeeID, newVal: testEmployee},
			assertion: assert.NoError,
		},
		{
			name:      "Error",
			e:         employeeMicroImpl{employeeStore: testEmployeeStore},
			args:      args{ctx: context.Background(), id: testEmployeNotFoundID, newVal: testEmployee},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, tt.e.Update(tt.args.ctx, tt.args.id, tt.args.newVal))
		})
	}
}
