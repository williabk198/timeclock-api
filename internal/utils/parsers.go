package utils

import (
	"strings"

	"github.com/williabk198/timeclock/internal/errors"
	"github.com/williabk198/timeclock/internal/models"
)

func ParsePronouns(input string) (models.Pronouns, error) {
	splitInput := strings.Split(input, "/")

	if len(splitInput) != 2 {
		return models.Pronouns{}, errors.NewInvalidFormatError("Pronouns", input, "they/them, he/him, etc...")
	}

	return models.Pronouns{
		Subject: splitInput[0],
		Object:  splitInput[1],
	}, nil
}
