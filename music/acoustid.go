package music

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os/exec"
	"regexp"

	h "github.com/barsanuphe/helpers"
)

const (
	acoustidURL = "http://api.acoustid.org/v2/lookup"
)

type AcoustidResults struct {
	Results []struct {
		ID         string `json:"id"`
		Recordings []struct {
			Artists []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"artists"`
			Duration int    `json:"duration"`
			ID       string `json:"id"`
			Releases []struct {
				Artists []struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"artists"`
				Country string `json:"country"`
				Date    struct {
					Day   int `json:"day"`
					Month int `json:"month"`
					Year  int `json:"year"`
				} `json:"date"`
				ID          string `json:"id"`
				MediumCount int    `json:"medium_count"`
				Mediums     []struct {
					Format     string `json:"format"`
					Position   int    `json:"position"`
					TrackCount int    `json:"track_count"`
					Tracks     []struct {
						Artists []struct {
							ID   string `json:"id"`
							Name string `json:"name"`
						} `json:"artists"`
						ID       string `json:"id"`
						Position int    `json:"position"`
						Title    string `json:"title"`
					} `json:"tracks"`
				} `json:"mediums"`
				Releaseevents []struct {
					Country string `json:"country"`
					Date    struct {
						Day   int `json:"day"`
						Month int `json:"month"`
						Year  int `json:"year"`
					} `json:"date"`
				} `json:"releaseevents"`
				Title      string `json:"title"`
				TrackCount int    `json:"track_count"`
			} `json:"releases"`
			Title string `json:"title"`
		} `json:"recordings"`
		Score float64 `json:"score"`
	} `json:"results"`
	Status string `json:"status"`
}

//------------------------

// AcousticID allows getting information about a track from its contents
type AcousticID struct {
	ApiKey      string
	Fingerprint string
	Duration    string
}

// NewAcoustid set up with api key
func NewAcoustid(key string) *AcousticID {
	return &AcousticID{ApiKey: key}
}

// CalculateFingerprint for a given track
func (a *AcousticID) CalculateFingerprint(path string) error {
	var err error
	fpcalc, err := exec.LookPath("fpcalc")
	if err != nil {
		fmt.Println("Needs fpcalc, installed with chromaprint!")
		return err
	}
	// run fpcalc on path
	out, err := exec.Command(fpcalc, path).Output()
	if err != nil {
		return err
	}
	data := string(out)
	// get duration & fingerprint
	r := regexp.MustCompile(`DURATION=(\d+)\nFINGERPRINT=([-\w]+)`)
	if r.MatchString(data) {
		a.Duration = r.FindStringSubmatch(data)[1]
		a.Fingerprint = r.FindStringSubmatch(data)[2]
		return nil
	}
	return errors.New("Could not find duration and fingerprint.")
}

// LookUp Acoustid database once we have the fingerprint
func (a *AcousticID) LookUp() (*AcoustidResults, error) {
	// get fingerprint
	if a.Fingerprint == "" {
		return nil, errors.New("Must fingerprint first")
	}
	// setting up the form
	buffer := new(bytes.Buffer)
	w := multipart.NewWriter(buffer)
	errs := []error{}
	errs = append(errs, w.WriteField("format", "json"))
	errs = append(errs, w.WriteField("client", a.ApiKey))
	errs = append(errs, w.WriteField("duration", a.Duration))
	errs = append(errs, w.WriteField("fingerprint", a.Fingerprint))
	errs = append(errs, w.WriteField("meta", "recordings releases tracks"))
	if err := h.CheckErrors(errs...); err != nil {
		return nil, err
	}
	w.Close()

	req, err := http.NewRequest("POST", acoustidURL, buffer)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	// TODO compress form? req.Header.Set("Content-Encoding", "gzip")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// parse response
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Returned status: " + resp.Status)
	}
	resultBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	results := AcoustidResults{}
	err = json.Unmarshal(resultBytes, &results)
	if results.Status != "ok" {
		return nil, errors.New("Acoustid Error")
	}
	return &results, nil
}
