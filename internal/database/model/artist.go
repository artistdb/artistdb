package model

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// NewArtist returns an artist with an initialized
func NewArtist() *Artist {
	return &Artist{ID: uuid.New().String()}
}

func WrapPronouns(pronouns []string) string {
	return strings.Join(pronouns, "|")
}

func UnwrapPronouns(dbValue string) []string {
	if len(dbValue) == 0 {
		return nil
	}

	return strings.Split(dbValue, "|")
}

// Artist represents an artist who is able to participate on an Event with some
// kind of Artwork.
type Artist struct {
	ID         string
	FirstName  string
	LastName   string
	ArtistName string
	Pronouns   []string
	Origin     Origin
	Language   string
	Socials    Socials
	BioGerman  string
	BioEnglish string
}

// Socials holds information about social media presences.
type Socials struct {
	Instagram string
	Facebook  string
	Bandcamp  string
}

// Origin holds information about the origin of an Artist.
type Origin struct {
	DateOfBirth  time.Time
	PlaceOfBirth string
	Nationality  string
}
