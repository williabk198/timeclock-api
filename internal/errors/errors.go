package errors

import "fmt"

type invalidFormatError struct {
	typeName       string
	recievedValue  string
	expectedFormat string
}

func (ife invalidFormatError) Error() string {
	return fmt.Sprintf(
		"%s is not a valid format for type %s; expected format %q",
		ife.recievedValue, ife.typeName, ife.expectedFormat,
	)
}

func NewInvalidFormatError(typeName, recievedValue, expectedFormat string) error {
	return invalidFormatError{
		typeName:       typeName,
		recievedValue:  recievedValue,
		expectedFormat: expectedFormat,
	}
}
