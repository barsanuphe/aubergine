package music

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/barsanuphe/helpers"
	u "github.com/barsanuphe/helpers/ui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscogs(t *testing.T) {
	fmt.Println("+ Testing Discogs...")
	check := assert.New(t)
	ui := &u.UI{}

	// get api key from env
	key := os.Getenv("DISCOGS_TOKEN")
	require.NotEqual(t, 0, len(key), "Cannot get Discogs application token")
	// get api secret from env
	secret := os.Getenv("DISCOGS_SECRET")
	require.NotEqual(t, 0, len(key), "Cannot get Discogs application secret")

	for _, t := range testMBReleases {
		fmt.Println("Testing with " + t.artist + " - " + t.albumTitle)
		a := NewDiscogsRelease(key, secret)

		if err := a.readCredentials(); err != nil {
			fmt.Println("COULD NOT TEST, CREDENTIALS MISSING")
			fmt.Println("For now, this test requires a valid discogs_credentials.json file, obviously not included in the repository.")
			return
		}

		err := a.Authorize(ui)
		check.Nil(err, "Error authorizing")

		err = a.LookUp(t.artist, t.albumTitle)
		check.Nil(err, "Error searching Discogs")
		check.NotEqual(0, len(a.Info.Results), "Expected hits")

		// Only checking first hit, not trying to find the right release now
		check.Equal(strings.ToLower(fmt.Sprintf("%s - %s", t.artist, t.albumTitle)), strings.ToLower(a.Info.Results[0].Title))

		rp := strings.NewReplacer(" ", "",
			"-", "")
		catno := rp.Replace(t.expectedCatalogNumber)

		found := false
		for _, r := range a.Info.Results {
			_, knownLabel := helpers.StringInSlice(t.expectedLabel, r.Label)
			if knownLabel && catno == rp.Replace(r.Catno) {
				found = true
				break
			}
		}
		check.Equal(true, found, "Release was not found on Discogs!")
		/*
			for _, r := range a.Info.Results {
				fmt.Println(r.ID, r.Title, r.Year, r.Country, r.Format, r.Genre, r.Label, r.Catno)
			}
		*/
	}

}
