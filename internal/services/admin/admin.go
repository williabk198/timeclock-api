package admin

import (
	"github.com/williabk198/timeclock/internal/datastores"
)

type Service interface {
	Contact() ContactMicro
	Person() PersonMicro
}

type adminService struct {
	personStore  datastores.PersonStore
	contactStore datastores.ContactDatastore
}

// Contact implements Service.
func (a adminService) Contact() ContactMicro {
	return contactMicroImpl{
		contactStore: a.contactStore,
	}
}

// Person implements Service.
func (a adminService) Person() PersonMicro {
	return personMicroImpl{
		personStore: a.personStore,
	}
}

func NewService(personStore datastores.PersonStore, contactStore datastores.ContactDatastore) Service {
	return adminService{
		personStore:  personStore,
		contactStore: contactStore,
	}
}
