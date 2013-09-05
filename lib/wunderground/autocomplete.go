package wunderground

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/StalkR/goircbot/lib/transport"
)

// An ACResult represents auto complete result from JSON.
type ACResult struct {
	Results []ACElement `json:"RESULTS"`
}

// An ACElement represents an auto complete element from JSON.
type ACElement struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	Country       string `json:"c"`
	ZMW           string `json:"zmw"`
	TimeZone      string `json:"tz"`
	TimeZoneShort string `json:"tzs"`
	URLPath       string `json:"l"`
}

// AutoComplete asks wunderground to auto complete a query.
// It does not require an API key.
// Documentation at http://www.wunderground.com/weather/api/d/docs?d=autocomplete-api.
func AutoComplete(query string) ([]ACElement, error) {
	base := "https://autocomplete.wunderground.com/aq"
	client, err := transport.Client(base)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("query", query)
	params.Set("format", "JSON")
	resp, err := client.Get(fmt.Sprintf("%s?%s", base, params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	r := ACResult{}
	err = json.Unmarshal(contents, &r)
	if err != nil {
		return nil, err
	}
	return r.Results, nil
}
