package geonames

import (
  "encoding/json"
  "fmt"
  "net/http"
  "net/url"
  "time"

  "github.com/StalkR/goircbot/lib/transport"
)

// Timezone returns the timezone for latitude and longitude coordinates.
func Timezone(username string, latitude, longitude string) (*time.Location, error) {
  v := url.Values{}
  v.Set("lat", latitude)
  v.Set("lng", longitude)
  v.Set("username", username)
  dest := fmt.Sprintf("http://api.geonames.org/timezoneJSON?%s", v.Encode())
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
  var r struct {
    TimezoneID string `json:"timezoneId"`
  }
  if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
    return nil, err
  }
  return time.LoadLocation(r.TimezoneID)
}
