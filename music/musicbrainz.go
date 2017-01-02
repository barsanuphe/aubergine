package music

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	musicBrainzAPIURL = "http://musicbrainz.org/ws/2/release/%s?inc=labels+artist-credits&fmt=json"
)

// MusicBrainzReleaseResults is a struct describing the JSON response for a MusicBreinz query about a speficif release.
type MusicBrainzReleaseResults struct {
	ArtistCredit []struct {
		Artist struct {
			Disambiguation string `json:"disambiguation"`
			ID             string `json:"id"`
			Name           string `json:"name"`
			SortName       string `json:"sort-name"`
		} `json:"artist"`
		Joinphrase string `json:"joinphrase"`
		Name       string `json:"name"`
	} `json:"artist-credit"`
	Asin            string `json:"asin"`
	Barcode         string `json:"barcode"`
	Country         string `json:"country"`
	CoverArtArchive struct {
		Artwork  bool `json:"artwork"`
		Back     bool `json:"back"`
		Count    int  `json:"count"`
		Darkened bool `json:"darkened"`
		Front    bool `json:"front"`
	} `json:"cover-art-archive"`
	Date           string `json:"date"`
	Disambiguation string `json:"disambiguation"`
	ID             string `json:"id"`
	LabelInfo      []struct {
		CatalogNumber string `json:"catalog-number"`
		Label         struct {
			Disambiguation string      `json:"disambiguation"`
			ID             string      `json:"id"`
			LabelCode      interface{} `json:"label-code"`
			Name           string      `json:"name"`
			SortName       string      `json:"sort-name"`
		} `json:"label"`
	} `json:"label-info"`
	Packaging     string `json:"packaging"`
	PackagingID   string `json:"packaging-id"`
	Quality       string `json:"quality"`
	ReleaseEvents []struct {
		Area struct {
			Disambiguation string   `json:"disambiguation"`
			ID             string   `json:"id"`
			Iso31661Codes  []string `json:"iso-3166-1-codes"`
			Name           string   `json:"name"`
			SortName       string   `json:"sort-name"`
		} `json:"area"`
		Date string `json:"date"`
	} `json:"release-events"`
	Status             string `json:"status"`
	StatusID           string `json:"status-id"`
	TextRepresentation struct {
		Language string `json:"language"`
		Script   string `json:"script"`
	} `json:"text-representation"`
	Title string `json:"title"`
}

// MusicBrainzRelease allows retrieving information from MusicBrainz
type MusicBrainzRelease struct {
	ID   string
	Info MusicBrainzReleaseResults
}

// NewMusicBrainzRelease set up with release ID
func NewMusicBrainzRelease(id string) *MusicBrainzRelease {
	return &MusicBrainzRelease{ID: id}
}

// GetInfo from MusicBrainz about a release
func (mb *MusicBrainzRelease) GetInfo() error {
	// musicbrainz lookup
	musicbrainzSearch := fmt.Sprintf(musicBrainzAPIURL, mb.ID)
	mbJSON, err := retrieveGetRequestData(&http.Client{}, musicbrainzSearch)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(mbJSON), &mb.Info)
	return nil
}
