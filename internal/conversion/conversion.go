package conversion

import (
	"time"

	model_gen "github.com/obitech/artist-db/graph/model"
	"github.com/obitech/artist-db/internal/database/model"
)

func PointerToString(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func PointerToTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}

	return t.UTC()
}

func TimeToPointer(t time.Time) *time.Time {
	return &t
}

func TimeToPString(t time.Time) *string {
	ret := t.Format(time.UnixDate)
	return &ret
}

// ArtistToGenArtist takes Artist objects (model type) and converts it to
// an Artist Object from our generated models
func ArtistToGenArtist(a []*model.Artist) ([]*model_gen.Artist, error) {
	ret := make([]*model_gen.Artist, len(a))


	for i, a := range a {

		ret[i] = &model_gen.Artist{
			ID:           a.ID,
			FirstName:    a.FirstName,
			LastName:     a.LastName,
			ArtistName:   &a.ArtistName,
			Pronouns:     []*string{},
			DateOfBirth:  TimeToPString(a.Origin.DateOfBirth),
			PlaceOfBirth: &a.Origin.PlaceOfBirth,
			Nationality:  &a.Origin.Nationality,
			Language:     &a.Language,
			Facebook:     &a.Socials.Facebook,
			Instagram:    &a.Socials.Instagram,
			Bandcamp:     &a.Socials.Bandcamp,
			BioGer:       &a.BioGerman,
			BioEn:        &a.BioEnglish,
		}

		if len(a.Pronouns) > 0 {
			ret[i].Pronouns = make([]*string, len(a.Pronouns))

			for j, pronoun := range a.Pronouns {
				ret[i].Pronouns[j] = &pronoun
			}
		}
	}

	return ret, nil
}

func LocationToGenLocation(l []*model.Location) ([]*model_gen.Location, error) {
	ret := make([]*model_gen.Location, len(l))

	for i, l := range l {
		ret[i] = &model_gen.Location{
			ID:   l.ID,
			Name: l.Name,
		}
	}

	return ret, nil
}
