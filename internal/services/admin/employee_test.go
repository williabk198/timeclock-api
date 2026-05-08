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

func Test_employeeMicroImpl_Add(t *testing.T) {
	type args struct {
		ctx      context.Context
		employee models.Employee
		metadata models.EmployeeMetadata
	}
	type wants struct {
		id uuid.UUID
	}

	testEmployeeID := uuid.New()
	testPersonID := uuid.New()

	testEmployee := models.Employee{
		PersonID:    testPersonID,
		ReportsToID: uuid.Nil,
		Title:       "Owner/President",
	}
	testEmployeeError := models.Employee{
		Title: "error_val",
	}

	testEmployeeMetadata := models.EmployeeMetadata{
		EmployeeID: testEmployeeID,
		Pay: models.EmployeePay{
			Currency: "USD",
			Rate:     19.0,
			Cadence:  models.PayCadenceHourly,
		},
		HireDate:  time.Date(2017, 5, 4, 0, 0, 0, 0, time.UTC),
		StartDate: time.Date(2017, 6, 5, 13, 0, 0, 0, time.UTC),
		SickTime:  20.0,
		TimeOff:   20.0,
		Exempt:    false,
		Status:    models.EmployeeStatusActive,
	}

	testEmployeeMetadataError := models.EmployeeMetadata{
		EmployeeID: testEmployeeID,
		Pay: models.EmployeePay{
			Currency: "invalid",
			Rate:     19.0,
			Cadence:  models.PayCadenceHourly,
		},
	}

	testEmployeeStore := &mockEmployeeStore{}
	testEmployeeStore.On("Add", mock.Anything, testEmployee, testEmployeeMetadata).Return(testEmployeeID, error(nil))
	testEmployeeStore.On("Add", mock.Anything, testEmployeeError, testEmployeeMetadata).Return(uuid.Nil, assert.AnError)
	testEmployeeStore.On("Add", mock.Anything, testEmployee, testEmployeeMetadataError).Return(uuid.Nil, assert.AnError)

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
			args:      args{ctx: context.Background(), employee: testEmployee, metadata: testEmployeeMetadata},
			wants:     wants{id: testEmployeeID},
			assertion: assert.NoError,
		},
		{
			name:      "Error; EmployeeData",
			e:         employeeMicroImpl{employeeStore: testEmployeeStore},
			args:      args{ctx: context.Background(), employee: testEmployeeError, metadata: testEmployeeMetadata},
			assertion: assert.Error,
		},
		{
			name:      "Error; Metadata",
			e:         employeeMicroImpl{employeeStore: testEmployeeStore},
			args:      args{ctx: context.Background(), employee: testEmployee, metadata: testEmployeeMetadataError},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, err := tt.e.Add(tt.args.ctx, tt.args.employee, tt.args.metadata)
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

func Test_employeeMicroImpl_UpdateEmployee(t *testing.T) {
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
	testEmployeeStore.On("UpdateEmployee", mock.Anything, testEmployeeID, testEmployee).Return(error(nil))
	testEmployeeStore.On("UpdateEmployee", mock.Anything, testEmployeNotFoundID, testEmployee).Return(assert.AnError)

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
			tt.assertion(t, tt.e.UpdateEmployee(tt.args.ctx, tt.args.id, tt.args.newVal))
		})
	}
}

func Test_employeeMicroImpl_UpdateExemptStatus(t *testing.T) {
	type args struct {
		ctx    context.Context
		id     uuid.UUID
		newVal bool
	}

	testEmployeeNotFoundID := uuid.New()
	testEmployeeID := uuid.New()

	testEmployeeStore := &mockEmployeeStore{}
	testEmployeeStore.On("UpdateExemptStatus", mock.Anything, testEmployeeID, true).Return(error(nil))
	testEmployeeStore.On("UpdateExemptStatus", mock.Anything, testEmployeeNotFoundID, false).Return(assert.AnError)

	tests := []struct {
		name      string
		emi       employeeMicroImpl
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			emi: employeeMicroImpl{
				employeeStore: testEmployeeStore,
			},
			args: args{
				ctx:    context.Background(),
				id:     testEmployeeID,
				newVal: true,
			},
			assertion: assert.NoError,
		},
		{
			name: "Error",
			emi: employeeMicroImpl{
				employeeStore: testEmployeeStore,
			},
			args: args{
				ctx:    context.Background(),
				id:     testEmployeeNotFoundID,
				newVal: false,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, tt.emi.UpdateExemptStatus(tt.args.ctx, tt.args.id, tt.args.newVal))
		})
	}
}

func Test_employeeMicroImpl_UpdatePay(t *testing.T) {
	type args struct {
		ctx    context.Context
		id     uuid.UUID
		newVal models.EmployeePay
	}

	testEmployeNotFoundID := uuid.New()
	testEmployeeID := uuid.New()

	testPayDataDB := models.EmployeePay{
		Currency: "USD",
		Rate:     37.0,
		Cadence:  models.PayCadenceHourly,
	}

	testEmployeeStore := &mockEmployeeStore{}
	testEmployeeStore.On("UpdatePay", mock.Anything, testEmployeeID, testPayDataDB).Return(error(nil))
	testEmployeeStore.On("UpdatePay", mock.Anything, testEmployeNotFoundID, testPayDataDB).Return(assert.AnError)

	tests := []struct {
		name      string
		emi       employeeMicroImpl
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			emi:  employeeMicroImpl{employeeStore: testEmployeeStore},
			args: args{
				ctx:    context.Background(),
				id:     testEmployeeID,
				newVal: testPayDataDB,
			},
			assertion: assert.NoError,
		},
		{
			name: "Error",
			emi:  employeeMicroImpl{employeeStore: testEmployeeStore},
			args: args{
				ctx:    context.Background(),
				id:     testEmployeNotFoundID,
				newVal: testPayDataDB,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, tt.emi.UpdatePay(tt.args.ctx, tt.args.id, tt.args.newVal))
		})
	}
}

func Test_employeeMicroImpl_UpdateSickTime(t *testing.T) {
	type args struct {
		ctx    context.Context
		id     uuid.UUID
		newVal float64
	}

	testEmployeNotFoundID := uuid.New()
	testEmployeeID := uuid.New()

	testEmployeeStore := &mockEmployeeStore{}
	testEmployeeStore.On("UpdateSickTime", mock.Anything, testEmployeeID, 16.0).Return(error(nil))
	testEmployeeStore.On("UpdateSickTime", mock.Anything, testEmployeNotFoundID, 16.0).Return(assert.AnError)

	tests := []struct {
		name      string
		emi       employeeMicroImpl
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			emi:  employeeMicroImpl{employeeStore: testEmployeeStore},
			args: args{
				ctx:    context.Background(),
				id:     testEmployeeID,
				newVal: 16.0,
			},
			assertion: assert.NoError,
		},
		{
			name: "Error",
			emi:  employeeMicroImpl{employeeStore: testEmployeeStore},
			args: args{
				ctx:    context.Background(),
				id:     testEmployeNotFoundID,
				newVal: 16.0,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			tt.assertion(t, tt.emi.UpdateSickTime(context.Background(), tt.args.id, tt.args.newVal))
		})
	}
}

func Test_employeeMicroImpl_UpdateStatus(t *testing.T) {
	type args struct {
		ctx    context.Context
		id     uuid.UUID
		newVal models.EmployeeStatus
	}

	testEmployeNotFoundID := uuid.New()
	testEmployeeID := uuid.New()

	testEmployeeStore := &mockEmployeeStore{}
	testEmployeeStore.On("UpdateStatus", mock.Anything, testEmployeNotFoundID, models.EmployeeStatusInactive).Return(assert.AnError)
	testEmployeeStore.On("UpdateStatus", mock.Anything, testEmployeeID, models.EmployeeStatusInactive).Return(error(nil))

	tests := []struct {
		name      string
		emi       employeeMicroImpl
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			emi:  employeeMicroImpl{employeeStore: testEmployeeStore},
			args: args{
				ctx:    context.Background(),
				id:     testEmployeeID,
				newVal: 2,
			},
			assertion: assert.NoError,
		},
		{
			name: "Error",
			emi:  employeeMicroImpl{employeeStore: testEmployeeStore},
			args: args{
				ctx:    context.Background(),
				id:     testEmployeNotFoundID,
				newVal: 2,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, tt.emi.UpdateStatus(tt.args.ctx, tt.args.id, tt.args.newVal))
		})
	}
}

func Test_employeeMicroImpl_UpdateTimeOff(t *testing.T) {
	type args struct {
		ctx    context.Context
		id     uuid.UUID
		newVal float64
	}

	testEmployeNotFoundID := uuid.New()
	testEmployeeID := uuid.New()

	testEmployeeStore := &mockEmployeeStore{}
	testEmployeeStore.On("UpdateTimeOff", mock.Anything, testEmployeeID, 32.0).Return(error(nil))
	testEmployeeStore.On("UpdateTimeOff", mock.Anything, testEmployeNotFoundID, 32.0).Return(assert.AnError)

	tests := []struct {
		name      string
		emi       employeeMicroImpl
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			emi:  employeeMicroImpl{employeeStore: testEmployeeStore},
			args: args{
				ctx:    context.Background(),
				id:     testEmployeeID,
				newVal: 32.0,
			},
			assertion: assert.NoError,
		},
		{
			name: "Error",
			emi:  employeeMicroImpl{employeeStore: testEmployeeStore},
			args: args{
				ctx:    context.Background(),
				id:     testEmployeNotFoundID,
				newVal: 32.0,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, tt.emi.UpdateTimeOff(tt.args.ctx, tt.args.id, tt.args.newVal))
		})
	}
}
