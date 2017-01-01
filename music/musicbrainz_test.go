package music

import (
	"fmt"
	"testing"

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
		"Solitude",
		"9b7a83cd-15c2-4d2f-9eea-da740c033517",
		"Verve",
		"314 519 810-2",
	},
	{
		"Radiohead",
		"Kid A",
		"0e8a1994-f0a7-481d-9be2-6c2f80e14de5",
		"XL Recordings",
		"XLDA782",
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
		check.Equal(0, len(a.Info.Label_info)) // no label info, CC material
	}

	for _, t := range testMBReleases {
		fmt.Println("Testing with " + t.mbReleaseID)
		a := NewMusicBrainzRelease(t.mbReleaseID)

		err := a.GetInfo()
		check.Nil(err, "Unexpected error getting MusicBrainz info")

		check.Equal(t.mbReleaseID, a.Info.ID)
		check.Equal(t.albumTitle, a.Info.Title)
		check.Equal(t.artist, a.Info.Artist_credit[0].Name)
		check.NotEqual(0, len(a.Info.Label_info), "Release should have label info")
		check.Equal(t.expectedCatalogNumber, a.Info.Label_info[0].Catalog_number)
		check.Equal(t.expectedLabel, a.Info.Label_info[0].Label.Name)
	}

}
