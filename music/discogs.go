package music

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	u "github.com/barsanuphe/helpers/ui"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/skratchdot/open-golang/open"
)

const (
	discogsSearchURL = "https://api.discogs.com/database/search"
	credentialsFile  = "discogs_credentials.json"
)

type DiscogsResults struct {
	Pagination struct {
		Items   int      `json:"items"`
		Page    int      `json:"page"`
		Pages   int      `json:"pages"`
		PerPage int      `json:"per_page"`
		Urls    struct{} `json:"urls"`
	} `json:"pagination"`
	Results []struct {
		Barcode   []string `json:"barcode"`
		Catno     string   `json:"catno"`
		Community struct {
			Have int `json:"have"`
			Want int `json:"want"`
		} `json:"community"`
		Country     string   `json:"country"`
		Format      []string `json:"format"`
		Genre       []string `json:"genre"`
		ID          int      `json:"id"`
		Label       []string `json:"label"`
		ResourceURL string   `json:"resource_url"`
		Style       []string `json:"style"`
		Thumb       string   `json:"thumb"`
		Title       string   `json:"title"`
		Type        string   `json:"type"`
		URI         string   `json:"uri"`
		Year        string   `json:"year"`
	} `json:"results"`
}

// DiscogsRelease retrieves information about a release on Discogs.
type DiscogsRelease struct {
	CredentialsFile string
	Token           string
	Secret          string
	UserToken       string
	UserSecret      string
	Client          oauth.Client
	Info            DiscogsResults
}

// NewDiscogsRelease set up with Discogs API authorization info.
func NewDiscogsRelease(token, secret string) *DiscogsRelease {
	return &DiscogsRelease{Token: token, Secret: secret, CredentialsFile: credentialsFile}
}

func (d *DiscogsRelease) readCredentials() error {
	b, err := ioutil.ReadFile(d.CredentialsFile)
	if err != nil {
		return err
	}
	userCred := oauth.Credentials{}
	if err := json.Unmarshal(b, &userCred); err != nil {
		return err
	}
	d.UserToken = userCred.Token
	d.UserSecret = userCred.Secret
	return nil
}

func (d *DiscogsRelease) saveCredentials() error {
	userCred := oauth.Credentials{Token: d.UserToken, Secret: d.UserSecret}
	jsonToSave, err := json.MarshalIndent(userCred, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(d.CredentialsFile, jsonToSave, 0777)
}

// Authorize with Discogs by OAuth
func (d *DiscogsRelease) Authorize(ui u.UserInterface) error {
	// init client
	d.Client = oauth.Client{
		TemporaryCredentialRequestURI: "https://api.discogs.com/oauth/request_token",
		ResourceOwnerAuthorizationURI: "https://www.discogs.com/oauth/authorize",
		TokenRequestURI:               "https://api.discogs.com/oauth/access_token",
		Header:                        http.Header{"User-Agent": {"AUBERGINE/1.0"}},
	}
	d.Client.Credentials.Token = d.Token
	d.Client.Credentials.Secret = d.Secret

	// get from d.CredentialsFile
	if err := d.readCredentials(); err != nil {
		ui.Warning("Could not get credentials, authorizing with Discogs.")
		// if we cant't, get them from discogs
		tempCred, err := d.Client.RequestTemporaryCredentials(nil, "", nil)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		// open in browser to authorize once
		if err := open.Start(d.Client.AuthorizationURL(tempCred, nil)); err != nil {
			fmt.Println("err redirecting for authorization")
			return err
		}

		// wait for user input (code given by discogs web page)
		fmt.Print("Enter token: ")
		tempToken, err := ui.GetInput()
		if err != nil {
			ui.Error("Could not get token!")
			return err
		}
		fmt.Println("got token: " + tempToken)

		//tempTokenCred := &oauth.Credentials{Token:tempToken}

		tokenCred, _, err := d.Client.RequestToken(nil, tempCred, tempToken)
		if err != nil {
			ui.Error("Could not request token!")
			return err
		}
		d.UserToken = tokenCred.Token
		d.UserSecret = tokenCred.Secret
		if err := d.saveCredentials(); err != nil {
			return err
		}
	}
	return nil
}

// LookUp release on Discogs and retrieve its information
func (d *DiscogsRelease) LookUp(artist, release string) error {
	// TODO check authorized
	// TODO see what to return

	// search
	searchURL, err := url.Parse(discogsSearchURL)
	if err != nil {
		return err
	}
	q := searchURL.Query()
	q.Set("type", "release")
	q.Set("artist", artist)
	q.Set("release_title", release)
	searchURL.RawQuery = q.Encode()

	respp, err := d.Client.Get(nil, &oauth.Credentials{Token: d.UserToken, Secret: d.UserSecret}, discogsSearchURL, q)
	if err != nil {
		return err
	}

	defer respp.Body.Close()
	if respp.StatusCode != http.StatusOK {
		return errors.New("Returned status: " + respp.Status)
	}

	resultDCBytes, err := ioutil.ReadAll(respp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(resultDCBytes, &d.Info); err != nil {
		return errors.New("Could not read JSON data from Discogs.")
	}

	// TODO retrieve track list!! GET https://api.discogs.com/releases/{release_id}
	return nil
}
