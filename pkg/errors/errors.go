package errors

import "fmt"

type nilValueError struct {
	varName string
}

func NewNilValueError(name string) error {
	return nilValueError{varName: name}
}

func (nve nilValueError) Error() string {
	return fmt.Sprintf("%q value was nil", nve.varName)
}
