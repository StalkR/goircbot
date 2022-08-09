// Package geonames implements geonames.org APIs.
package geonames

import (
  "encoding/json"
  "errors"
  "fmt"
  "net/http"
  "net/url"

  "github.com/StalkR/goircbot/lib/transport"
)

// ErrNotFound is returned when no result is found.
var ErrNotFound = errors.New("not found")

// A Location represents a geo location with a few properties.
type Location struct {
  Name      string `json:"name"`
  Country   string `json:"countryName"`
  Latitude  string `json:"lat"`
  Longitude string `json:"lng"`
}

type searchResponse struct {
  Geonames []Location `json:"geonames"`
}

// Search searches a geo name query.
func Search(username string, q string) (*Location, error) {
  v := url.Values{}
  v.Set("q", q)
  v.Set("maxRows", "1")
  v.Set("type", "json")
  v.Set("username", username)
  dest := fmt.Sprintf("http://api.geonames.org/search?%s", v.Encode())
  client, err := transport.Client(dest)
  if err != nil {
    return nil, err
  }
  resp, err := client.Get(dest)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()
  if resp.StatusCode != http.StatusOK {
    return nil, fmt.Errorf("status %v", resp.Status)
  }
  var r searchResponse
  if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
    return nil, err
  }
  if len(r.Geonames) == 0 {
    return nil, ErrNotFound
  }
  return &r.Geonames[0], nil
}
