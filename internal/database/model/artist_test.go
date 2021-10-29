package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewArtist(t *testing.T) {
	require.NotEmpty(t, NewArtist().ID)
}

func TestWrapPronouns(t *testing.T) {
	tt := []struct {
		name     string
		pronouns []string
		want     string
	}{
		{
			name: "empty pronouns yields empty string",
		},
		{
			name:     "empty strings are appended",
			pronouns: []string{"", ""},
			want:     "|",
		},
		{
			name:     "single pronoun yields single string",
			pronouns: []string{"she"},
			want:     "she",
		},
		{
			name:     "two pronouns are appended correctly",
			pronouns: []string{"she", "her"},
			want:     "she|her",
		},
		{
			name:     "two pronouns are appended correctly",
			pronouns: []string{"she", "her", "him"},
			want:     "she|her|him",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, WrapPronouns(tc.pronouns))
		})
	}
}

func TestUnwrapPronouns(t *testing.T) {
	tt := []struct {
		name  string
		dbVal string
		want  []string
	}{
		{
			name: "empty string yields empty result",
		},
		{
			name:  "single pronoun yields single string",
			dbVal: "she",
			want:  []string{"she"},
		},
		{
			name:  "double pronouns are unwrapped",
			dbVal: "she|her",
			want:  []string{"she", "her"},
		},
		{
			name:  "tripple pronouns are unwrapped",
			dbVal: "she|her|him",
			want:  []string{"she", "her", "him"},
		},
		{
			name:  "empty strings are unwrapped too",
			dbVal: "she|her|him|",
			want:  []string{"she", "her", "him", ""},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, UnwrapPronouns(tc.dbVal))
		})
	}
}
