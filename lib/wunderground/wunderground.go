// Package wunderground implements wunderground.com weather API.
package wunderground

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/StalkR/goircbot/lib/transport"
)

// An ConditionResult represents weather condition result from JSON.
type ConditionResult struct {
	Response APIResponse `json:"response"`
	Current  Weather     `json:"current_observation"`
}

// An APIResponse represents a response from JSON.
// Not all fields are represented.
type APIResponse struct {
	Version        string   `json:"version"`
	TermsOfService string   `json:"termsofService"`
	Error          APIError `json:"error"`
}

// An APIError represents an error from JSON.
// It implements error interface.
type APIError struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// IsZero reports whether there was an error or not.
func (e APIError) IsZero() bool {
	return e.Type == "" && e.Description == ""
}

// Error formats an API error with type and description.
func (e APIError) Error() string {
	return fmt.Sprintf("wunderground: %s, %s", e.Type, e.Description)
}

// A Weather represents a weather condition from JSON.
// Not all fields are represented.
type Weather struct {
	Display       Location `json:"display_location"`
	Observation   Location `json:"observation_location"`
	TZ            string   `json:"local_tz_long"`
	Weather       string   `json:"weather"`
	TempF         float64  `json:"temp_f"`
	TempC         float64  `json:"temp_c"`
	Humidity      string   `json:"relative_humidity"`
	Wind          string   `json:"wind_string"`
	WindDirection string   `json:"wind_dir"`
	WindMPH       float64  `json:"wind_mph"`
	WindKPH       float64  `json:"wind_kph"`
}

// A Location represents weather location from JSON.
type Location struct {
	Full           string `json:"full"`
	City           string `json:"city"`
	State          string `json:"state"`
	StateName      string `json:"state_name"`
	Country        string `json:"country"`
	CountryISO3166 string `json:"country_iso3166"`
	ZIP            string `json:"zip"`
	Latitude       string `json:"latitude"`
	Longitude      string `json:"longitude"`
	Elevation      string `json:"elevation"`
}

// Time returns local time at location.
func (w *Weather) Time() (time.Time, error) {
	loc, err := time.LoadLocation(w.TZ)
	if err != nil {
		return time.Time{}, fmt.Errorf("wunderground: invalid location %s: %v", w.TZ, err)
	}
	return time.Now().In(loc), nil
}

// String formats a weather result on one line.
func (w *Weather) String() string {
	var ts string
	if t, err := w.Time(); err != nil {
		log.Print(err)
	} else {
		ts = fmt.Sprintf(" - %s", t.Format(time.RFC1123))
	}
	return fmt.Sprintf("%v: %v (%.2f°C), humidity %v, wind %v (%v %.2f km/h)%s",
		w.Display.Full, w.Weather, w.TempC, w.Humidity, w.Wind, w.WindDirection, w.WindKPH, ts)
}

// Conditions requests weather conditions for a location.
// It performs auto complete to get the right weather location and uses first result.
// Documentation at http://www.wunderground.com/weather/api/d/docs?d=data/conditions.
func Conditions(APIKey, location string) (*Weather, error) {
	r, err := AutoComplete(location)
	if err != nil {
		return nil, err
	}
	if len(r) == 0 {
		return nil, errors.New("wunderground: location not found")
	}
	return r[0].Conditions(APIKey)
}

// Conditions requests weather conditions for an auto complete result.
func (a *ACElement) Conditions(APIKey string) (*Weather, error) {
	url := fmt.Sprintf("https://api.wunderground.com/api/%s/conditions%s.json",
		APIKey, a.URLPath)
	client, err := transport.Client(url)
	if err != nil {
		return nil, err
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	contents = bytes.Replace(contents, []byte("\\'"), []byte("'"), -1)
	r := ConditionResult{}
	err = json.Unmarshal(contents, &r)
	if err != nil {
		return nil, err
	}
	if !r.Response.Error.IsZero() {
		return nil, r.Response.Error
	}
	return &r.Current, nil
}
