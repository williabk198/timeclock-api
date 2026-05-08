package admin

import (
	"github.com/williabk198/timeclock/internal/datastores"
)

type Service interface {
	Contact() ContactMicro
	Employee() EmployeeMicro
	Person() PersonMicro
}

type adminService struct {
	personStore   datastores.PersonStore
	contactStore  datastores.ContactDatastore
	employeeStore datastores.EmployeeDatastore
}

// Contact implements Service.
func (a adminService) Contact() ContactMicro {
	return contactMicroImpl{
		contactStore: a.contactStore,
	}
}

// Employee implements Service.
func (a adminService) Employee() EmployeeMicro {
	return employeeMicroImpl{
		employeeStore: a.employeeStore,
	}
}

// Person implements Service.
func (a adminService) Person() PersonMicro {
	return personMicroImpl{
		personStore: a.personStore,
	}
}

func NewService(personStore datastores.PersonStore, contactStore datastores.ContactDatastore, employeeStore datastores.EmployeeDatastore) Service {
	return adminService{
		personStore:   personStore,
		contactStore:  contactStore,
		employeeStore: employeeStore,
	}
}
