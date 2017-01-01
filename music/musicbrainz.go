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
	Artist_credit []struct {
		Artist struct {
			Disambiguation string `json:"disambiguation"`
			ID             string `json:"id"`
			Name           string `json:"name"`
			Sort_name      string `json:"sort-name"`
		} `json:"artist"`
		Joinphrase string `json:"joinphrase"`
		Name       string `json:"name"`
	} `json:"artist-credit"`
	Asin              string `json:"asin"`
	Barcode           string `json:"barcode"`
	Country           string `json:"country"`
	Cover_art_archive struct {
		Artwork  bool `json:"artwork"`
		Back     bool `json:"back"`
		Count    int  `json:"count"`
		Darkened bool `json:"darkened"`
		Front    bool `json:"front"`
	} `json:"cover-art-archive"`
	Date           string `json:"date"`
	Disambiguation string `json:"disambiguation"`
	ID             string `json:"id"`
	Label_info     []struct {
		Catalog_number string `json:"catalog-number"`
		Label          struct {
			Disambiguation string      `json:"disambiguation"`
			ID             string      `json:"id"`
			Label_code     interface{} `json:"label-code"`
			Name           string      `json:"name"`
			Sort_name      string      `json:"sort-name"`
		} `json:"label"`
	} `json:"label-info"`
	Packaging      string `json:"packaging"`
	Packaging_id   string `json:"packaging-id"`
	Quality        string `json:"quality"`
	Release_events []struct {
		Area struct {
			Disambiguation   string   `json:"disambiguation"`
			ID               string   `json:"id"`
			Iso_3166_1_codes []string `json:"iso-3166-1-codes"`
			Name             string   `json:"name"`
			Sort_name        string   `json:"sort-name"`
		} `json:"area"`
		Date string `json:"date"`
	} `json:"release-events"`
	Status              string `json:"status"`
	Status_id           string `json:"status-id"`
	Text_representation struct {
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
