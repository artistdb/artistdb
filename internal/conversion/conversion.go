package conversion

import (
	"fmt"
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

	fmt.Println(len(ret))

	for i, a := range a {

		var re model_gen.Artist

		fmt.Println(a)
		re.ID = a.ID
		re.FirstName = a.FirstName
		re.LastName = a.LastName
		re.ArtistName = &a.ArtistName

		if len(a.Pronouns) > 0 {
			re.Pronouns = make([]*string, len(a.Pronouns))

			for i, pronoun := range a.Pronouns {
				re.Pronouns[i] = &pronoun
			}
		}
		t := a.Origin.DateOfBirth.Format(time.UnixDate)
		re.DateOfBirth = &t

		re.PlaceOfBirth = &a.Origin.PlaceOfBirth
		re.Language = &a.Language
		re.Facebook = &a.Socials.Facebook
		re.Instagram = &a.Socials.Instagram
		re.Bandcamp = &a.Socials.Bandcamp
		re.BioGer = &a.BioGerman
		re.BioEn = &a.BioEnglish

		ret[i] = &re
	}

	return ret, nil
}
