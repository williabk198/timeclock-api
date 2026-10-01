package core

import "github.com/williabk198/timeclock/internal/datastores"

type Service interface {
	Punch() PunchMicro
}

func NewCoreServices(punchesStore datastores.PunchStore) Service {
	panic("unimplemented")
}
