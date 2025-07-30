package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/williabk198/timeclock/internal/models"
)

func TestParsePronouns(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name      string
		args      args
		want      models.Pronouns
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				input: "they/them",
			},
			want: models.Pronouns{
				Subject: "they",
				Object:  "them",
			},
			assertion: assert.NoError,
		},
		{
			name: "Error",
			args: args{
				input: "invalid",
			},
			want:      models.Pronouns{},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePronouns(tt.args.input)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
