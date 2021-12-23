package artist

import (
	"time"

	"github.com/google/uuid"
)

// New returns an artist with an initialized ID.
func New() *Artist {
	return &Artist{ID: uuid.New().String()}
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
