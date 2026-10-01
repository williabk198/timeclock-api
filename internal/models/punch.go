package models

import (
	"time"

	"github.com/google/uuid"
)

type PunchKind int

const (
	PunchKindIn PunchKind = iota + 1
	PunchKindMealOut
	PunchKindMealIn
	PunchKindOut
)

type Punch struct {
	ID           uuid.UUID `jagsqlb:"id;omit"`
	EmployeeID   uuid.UUID `jagsqlb:"employee_id"`
	Kind         PunchKind `jagsqlb:"kind"`
	ReportedTime time.Time `jagsqlb:"reported_time"`
	Timestamp    time.Time `jagsqlb:"timestamp"`
}
