package endpoints

import (
	"context"
	"fmt"
	"time"

	"github.com/williabk198/timeclock/internal/models"
	"github.com/williabk198/timeclock/internal/services/admin"
	"github.com/williabk198/timeclock/internal/utils"
)

type adminPersonEndpoints struct {
	adminService admin.Service
}

// Add implements PersonEndpoints.
func (ape adminPersonEndpoints) Add(ctx context.Context, person PersonData) (PersonData, error) {
	pronouns, err := utils.ParsePronouns(person.Pronouns)
	if err != nil {
		return PersonData{}, fmt.Errorf("failed to process pronoun data: %w", err)
	}

	dbPerson := models.Person{
		Name:        person.Name,
		DateOfBirth: time.Unix(person.DateOfBirth, 0),
		Gender:      models.Gender(person.Gender),
		Pronouns:    pronouns,
	}

	id, err := ape.adminService.AddPerson(ctx, dbPerson)
	if err != nil {
		return PersonData{}, fmt.Errorf("failed to add person to database: %w", err)
	}

	person.ID = id.String()
	return person, nil
}
