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

func Test_adminEmployeeEndpoints_Add(t *testing.T) {
	type args struct {
		ctx      context.Context
		employee EmployeeData
	}

	testGoodPersonID := uuid.New()
	testGoodEmployeeID := uuid.New()

	testGoodEmployeeMetadata := &EmployeeMetadata{
		Pay:       models.EmployeePay{Currency: "USD", Rate: 50, Cadence: models.PayCadenceHourly},
		HireDate:  1493956800,
		StartDate: 1496894400,
		SickTime:  24.0,
		TimeOff:   24.0,
		Exempt:    true,
		Status:    1,
	}
	testBadEmployeeMetadata := &EmployeeMetadata{
		Pay:       models.EmployeePay{Currency: "USD", Rate: 250_000.0, Cadence: models.PayCadenceYearly},
		HireDate:  -1,
		StartDate: -1,
		SickTime:  -1.0,
		TimeOff:   -1.0,
		Exempt:    true,
		Status:    -1,
	}

	testGoodEmployeeDataGoodMeta := EmployeeData{
		PersonID:    testGoodPersonID.String(),
		ReportsToID: uuid.Nil.String(),
		Title:       "Owner",
		Metadata:    testGoodEmployeeMetadata,
	}
	testGoodEmployeeDataBadMeta := EmployeeData{
		PersonID:    testGoodPersonID.String(),
		ReportsToID: uuid.Nil.String(),
		Title:       "Owner",
		Metadata:    testBadEmployeeMetadata,
	}

	testGoodEmployeeDB := models.Employee{
		PersonID:    testGoodPersonID,
		ReportsToID: uuid.Nil,
		Title:       testGoodEmployeeDataGoodMeta.Title,
	}
	testGoodEmployeeMetadataDB := models.EmployeeMetadata{
		Pay:       models.EmployeePay{Currency: "USD", Rate: 50.0, Cadence: models.PayCadenceHourly},
		HireDate:  time.UnixMilli(1493956800),
		StartDate: time.UnixMilli(1496894400),
		SickTime:  24.0,
		TimeOff:   24.0,
		Exempt:    true,
		Status:    models.EmployeeStatusActive,
	}

	testBadEmployeeData := EmployeeData{
		PersonID:    uuid.NewString(),
		ReportsToID: uuid.NewString(),
		Title:       "error val",
		Metadata:    testGoodEmployeeMetadata,
	}
	testBadEmployeeDB := models.Employee{
		PersonID:    uuid.MustParse(testBadEmployeeData.PersonID),
		ReportsToID: uuid.MustParse(testBadEmployeeData.ReportsToID),
		Title:       testBadEmployeeData.Title,
	}
	testBadEmployeeMetadataDB := models.EmployeeMetadata{
		Pay:       models.EmployeePay{Currency: "USD", Rate: 250_000.0, Cadence: models.PayCadenceYearly},
		HireDate:  time.UnixMilli(-1),
		StartDate: time.UnixMilli(-1),
		SickTime:  -1.0,
		TimeOff:   -1.0,
		Exempt:    true,
		Status:    models.EmployeeStatus(-1),
	}

	testEmployeeMicro := &mockEmployeeMicro{}
	testEmployeeMicro.On("Add", mock.Anything, testGoodEmployeeDB, testGoodEmployeeMetadataDB).Return(testGoodEmployeeID, error(nil))
	testEmployeeMicro.On("Add", mock.Anything, testBadEmployeeDB, testGoodEmployeeMetadataDB).Return(uuid.Nil, assert.AnError)
	testEmployeeMicro.On("Add", mock.Anything, testGoodEmployeeDB, testBadEmployeeMetadataDB).Return(uuid.Nil, assert.AnError)

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
			args: args{ctx: context.Background(), employee: testGoodEmployeeDataGoodMeta},
			want: EmployeeData{
				ID:          testGoodEmployeeID.String(),
				PersonID:    testGoodEmployeeDataGoodMeta.PersonID,
				ReportsToID: testGoodEmployeeDataGoodMeta.ReportsToID,
				Title:       testGoodEmployeeDataGoodMeta.Title,
				Metadata:    testGoodEmployeeMetadata,
			},
			assertion: assert.NoError,
		},
		{
			name:      "Error; Invalid Input",
			a:         adminEmployeeEndpoints{employeeMicro: testEmployeeMicro},
			args:      args{ctx: context.Background(), employee: EmployeeData{PersonID: "invalid_value"}},
			assertion: assert.Error,
		},
		{
			name:      "Error; Service Error; Bad Employee Data",
			a:         adminEmployeeEndpoints{employeeMicro: testEmployeeMicro},
			args:      args{ctx: context.Background(), employee: testBadEmployeeData},
			assertion: assert.Error,
		},
		{
			name:      "Error; Service Error; Bad Metadata",
			a:         adminEmployeeEndpoints{employeeMicro: testEmployeeMicro},
			args:      args{ctx: context.Background(), employee: testGoodEmployeeDataBadMeta},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.Add(tt.args.ctx, tt.args.employee)
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
		PersonID:    testEmployeeID,
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

func Test_adminEmployeeEndpoints_UpdateEmployee(t *testing.T) {
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
	testEmployeeMicro.On("UpdateEmployee", mock.Anything, testEmployeeID, testEmployeeDB).Return(error(nil))
	testEmployeeMicro.On("UpdateEmployee", mock.Anything, testNotFoundID, testEmployeeDB).Return(assert.AnError)

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

func Test_adminEmployeeEndpoints_UpdateExemptStatus(t *testing.T) {
	type args struct {
		urd UpdateRequestData[bool]
	}

	testEmployeeID := uuid.New()
	testNotFoundID := uuid.New()

	testEmployeeMicro := &mockEmployeeMicro{}
	testEmployeeMicro.On("UpdateExemptStatus", mock.Anything, testEmployeeID, true).Return(error(nil))
	testEmployeeMicro.On("UpdateExemptStatus", mock.Anything, testNotFoundID, true).Return(assert.AnError)

	tests := []struct {
		name      string // description of this test case
		a         adminEmployeeEndpoints
		args      args
		want      bool
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			a: adminEmployeeEndpoints{
				employeeMicro: testEmployeeMicro,
			},
			args: args{
				urd: UpdateRequestData[bool]{
					ID:   testEmployeeID.String(),
					Data: true,
				},
			},
			want:      true,
			assertion: assert.NoError,
		},
		{
			name: "Error",
			a: adminEmployeeEndpoints{
				employeeMicro: testEmployeeMicro,
			},
			args: args{
				urd: UpdateRequestData[bool]{
					ID:   testNotFoundID.String(),
					Data: true,
				},
			},
			want:      false,
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := tt.a.UpdateExemptStatus(context.Background(), tt.args.urd)
			tt.assertion(t, gotErr)
			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_adminEmployeeEndpoints_UpdatePay(t *testing.T) {
	type args struct {
		urd UpdateRequestData[models.EmployeePay]
	}

	testNotFoundID := uuid.New()
	testEmployeeID := uuid.New()

	testPayDataDB := models.EmployeePay{
		Currency: "USD",
		Rate:     37.0,
		Cadence:  models.PayCadenceHourly,
	}

	testEmployeeMicro := &mockEmployeeMicro{}
	testEmployeeMicro.On("UpdatePay", mock.Anything, testEmployeeID, testPayDataDB).Return(error(nil))
	testEmployeeMicro.On("UpdatePay", mock.Anything, testNotFoundID, testPayDataDB).Return(assert.AnError)

	tests := []struct {
		name      string // description of this test case
		a         adminEmployeeEndpoints
		args      args
		want      models.EmployeePay
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			a: adminEmployeeEndpoints{
				employeeMicro: testEmployeeMicro,
			},
			args: args{
				urd: UpdateRequestData[models.EmployeePay]{
					ID:   testEmployeeID.String(),
					Data: testPayDataDB,
				},
			},
			want:      testPayDataDB,
			assertion: assert.NoError,
		},
		{
			name: "Error; Not Found",
			a: adminEmployeeEndpoints{
				employeeMicro: testEmployeeMicro,
			},
			args: args{
				urd: UpdateRequestData[models.EmployeePay]{
					ID:   testNotFoundID.String(),
					Data: testPayDataDB,
				},
			},
			want:      models.EmployeePay{},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := tt.a.UpdatePay(context.Background(), tt.args.urd)
			tt.assertion(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_adminEmployeeEndpoints_UpdateSickTimeHours(t *testing.T) {
	type args struct {
		urd UpdateRequestData[float64]
	}

	testNotFoundID := uuid.New()
	testEmployeeID := uuid.New()

	testEmployeeMicro := &mockEmployeeMicro{}
	testEmployeeMicro.On("UpdateSickTime", mock.Anything, testEmployeeID, 20.0).Return(error(nil))
	testEmployeeMicro.On("UpdateSickTime", mock.Anything, testNotFoundID, 20.0).Return(assert.AnError)

	tests := []struct {
		name      string // description of this test case
		a         adminEmployeeEndpoints
		args      args
		want      float64
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			a: adminEmployeeEndpoints{
				employeeMicro: testEmployeeMicro,
			},
			args: args{
				urd: UpdateRequestData[float64]{
					ID:   testEmployeeID.String(),
					Data: 20.0,
				},
			},
			want:      20.0,
			assertion: assert.NoError,
		},
		{
			name: "Error; Not Found",
			a: adminEmployeeEndpoints{
				employeeMicro: testEmployeeMicro,
			},
			args: args{
				urd: UpdateRequestData[float64]{
					ID:   testNotFoundID.String(),
					Data: 20.0,
				},
			},
			want:      -1.0,
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := tt.a.UpdateSickTimeHours(context.Background(), tt.args.urd)
			tt.assertion(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_adminEmployeeEndpoints_UpdateStatus(t *testing.T) {
	type args struct {
		urd UpdateRequestData[int]
	}

	testNotFoundID := uuid.New()
	testEmployeeID := uuid.New()

	testEmployeeMicro := &mockEmployeeMicro{}
	testEmployeeMicro.On("UpdateStatus", mock.Anything, testEmployeeID, models.EmployeeStatusGone).Return(error(nil))
	testEmployeeMicro.On("UpdateStatus", mock.Anything, testNotFoundID, models.EmployeeStatusGone).Return(assert.AnError)

	tests := []struct {
		name      string // description of this test case
		a         adminEmployeeEndpoints
		args      args
		want      int
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			a: adminEmployeeEndpoints{
				employeeMicro: testEmployeeMicro,
			},
			args: args{
				urd: UpdateRequestData[int]{
					ID:   testEmployeeID.String(),
					Data: int(models.EmployeeStatusGone),
				},
			},
			want:      int(models.EmployeeStatusGone),
			assertion: assert.NoError,
		},
		{
			name: "Error; Not Found",
			a: adminEmployeeEndpoints{
				employeeMicro: testEmployeeMicro,
			},
			args: args{
				urd: UpdateRequestData[int]{
					ID:   testNotFoundID.String(),
					Data: int(models.EmployeeStatusGone),
				},
			},
			want:      -1,
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := tt.a.UpdateStatus(context.Background(), tt.args.urd)
			tt.assertion(t, gotErr)
			assert.Equal(t, tt.want, got)

		})
	}
}

func Test_adminEmployeeEndpoints_UpdateTimeOffHours(t *testing.T) {
	type args struct {
		urd UpdateRequestData[float64]
	}

	testNotFoundID := uuid.New()
	testEmployeeID := uuid.New()

	testEmployeeMicro := &mockEmployeeMicro{}
	testEmployeeMicro.On("UpdateTimeOff", mock.Anything, testEmployeeID, 20.0).Return(error(nil))
	testEmployeeMicro.On("UpdateTimeOff", mock.Anything, testNotFoundID, 20.0).Return(assert.AnError)
	tests := []struct {
		name      string // description of this test case
		a         adminEmployeeEndpoints
		args      args
		urd       UpdateRequestData[float64]
		want      float64
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			a: adminEmployeeEndpoints{
				employeeMicro: testEmployeeMicro,
			},
			args: args{
				urd: UpdateRequestData[float64]{
					ID:   testEmployeeID.String(),
					Data: 20.0,
				},
			},
			want:      20.0,
			assertion: assert.NoError,
		},
		{
			name: "Error; Not Found",
			a: adminEmployeeEndpoints{
				employeeMicro: testEmployeeMicro,
			},
			args: args{
				urd: UpdateRequestData[float64]{
					ID:   testNotFoundID.String(),
					Data: 20.0,
				},
			},
			want:      -1.0,
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := tt.a.UpdateTimeOffHours(context.Background(), tt.args.urd)
			tt.assertion(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
