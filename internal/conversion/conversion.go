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

// ArtistToGenArtist takes Artist objects (model type) and converts it to
// an Artist Object from our generated models
func ArtistToGenArtist(a []*model.Artist) ([]*model_gen.Artist, error) {
	ret := make([]*model_gen.Artist, len(a))

	for i, a := range a {
		ret[i] = &model_gen.Artist{}

		ret[i].ID = a.ID
		ret[i].FirstName = a.FirstName
		ret[i].LastName = a.LastName
		ret[i].ArtistName = &a.ArtistName

		if len(a.Pronouns) > 0 {
			ret[i].Pronouns = make([]*string, len(a.Pronouns))

			for j, pronoun := range a.Pronouns {
				ret[i].Pronouns[j] = &pronoun
			}
		}
		t := a.Origin.DateOfBirth.Format(time.UnixDate)
		ret[i].DateOfBirth = &t
		ret[i].PlaceOfBirth = &a.Origin.PlaceOfBirth
		ret[i].Language = &a.Language
		ret[i].Facebook = &a.Socials.Facebook
		ret[i].Instagram = &a.Socials.Instagram
		ret[i].Bandcamp = &a.Socials.Bandcamp
		ret[i].BioGer = &a.BioGerman
		ret[i].BioEn = &a.BioEnglish
	}

	return ret, nil
}
