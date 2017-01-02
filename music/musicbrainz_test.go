package music

import (
	"fmt"
	"testing"

	"strings"

	"github.com/stretchr/testify/assert"
)

// Creative Commons tracks, see test/source.md
var testMBReleases = []struct {
	artist                string
	albumTitle            string
	mbReleaseID           string
	expectedLabel         string
	expectedCatalogNumber string
}{
	{
		"Billie Holiday",
		"Lady Sings the Blues",
		"9e4bfa2d-af1d-4f64-a495-88fe2144cabb",
		"PolyGram",
		"833 770-2",
	},
	{
		"Radiohead",
		"Kid A",
		"a3b0e5eb-fa3b-3e4d-b5e6-d0881984a183",
		"Parlophone",
		"527 7532",
	},
}

func TestMusicBrainz(t *testing.T) {
	fmt.Println("+ Testing MusicBrainz...")
	check := assert.New(t)

	for _, t := range testTracks {
		fmt.Println("Testing with " + t.path)
		a := NewMusicBrainzRelease(t.mbReleaseID)

		err := a.GetInfo()
		check.Nil(err, "Unexpected error getting MusicBrainz info")

		check.Equal(t.mbReleaseID, a.Info.ID)
		check.Equal(t.albumTitle, a.Info.Title)
		check.Equal(0, len(a.Info.LabelInfo)) // no label info, CC material
	}

	for _, t := range testMBReleases {
		fmt.Println("Testing with " + t.mbReleaseID)
		a := NewMusicBrainzRelease(t.mbReleaseID)

		err := a.GetInfo()
		check.Nil(err, "Unexpected error getting MusicBrainz info")

		check.Equal(t.mbReleaseID, a.Info.ID)
		check.Equal(strings.ToLower(t.albumTitle), strings.ToLower(a.Info.Title))
		check.Equal(t.artist, a.Info.ArtistCredit[0].Name)
		check.NotEqual(0, len(a.Info.LabelInfo), "Release should have label info")
		check.Equal(t.expectedCatalogNumber, a.Info.LabelInfo[0].CatalogNumber)
		check.Equal(t.expectedLabel, a.Info.LabelInfo[0].Label.Name)
		fmt.Println(a.Info.LabelInfo[0].Label.Disambiguation)
	}

}
