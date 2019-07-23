// Package openweathermap implements openweathermap.org API.
package openweathermap

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/StalkR/goircbot/lib/transport"
)

// A Response represents the parameters in API response.
// https://openweathermap.org/current#parameter
type Response struct {
	Coord struct {
		Longitude float64 `json:"lon"`
		Latitude  float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string
	Main struct {
		Temperature         float64 `json:"temp"`
		Pressure            float64 `json:"pressure"`
		Humidity            int     `json:"humidity"`
		MinTemperature      float64 `json:"temp_min"`
		MaxTemperature      float64 `json:"temp_max"`
		SeaLevelPressure    float64 `json:"sea_level"`
		GroundLevelPressure float64 `json:"grnd_level"`
	} `json:"main"`
	Wind struct {
		Speed   float64 `json:"speed"` // meter/second
		Degrees float64 `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Rain struct {
		OneHour    float64 `json:"1h"`
		ThreeHours float64 `json:"3h"`
	} `json:"rain"`
	Snow struct {
		OneHour    float64 `json:"1h"`
		ThreeHours float64 `json:"3h"`
	} `json:"snow"`
	Timestamp uint `json:"dt"`
	Sys       struct {
		Type    string `json:"type"`
		ID      string `json:"id"`
		Message string `json:"message"`
		Country string `json:"country"`
		Sunrise uint   `json:"sunrise"`
		Sunset  uint   `json:"sunset"`
	} `json:"sys"`
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// String formats a weather result on one line.
func (s *Response) String() string {
	var conditions []string
	for _, v := range s.Weather {
		conditions = append(conditions, v.Description)
	}
	name := s.Name
	if s.Sys.Country != "" {
		name = fmt.Sprintf("%v, %v", s.Name, s.Sys.Country)
	}
	description := strings.Join(conditions, ", ")
	wind := degreesToCompass(s.Wind.Degrees)
	windKPH := s.Wind.Speed * 3600 / 1000
	more := fmt.Sprintf("https://openweathermap.org/city/%v", s.ID)
	return fmt.Sprintf("%v: %.0f°C %v, %v%% humidity, wind %v %.0f km/h - %v",
		name, s.Main.Temperature, description, s.Main.Humidity, wind, windKPH, more)
}

var compass = [...]string{"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"}

func degreesToCompass(degrees float64) string {
	n := int((float64(degrees) / 22.5) + 0.5)
	return compass[n%len(compass)]
}

// A findResponse represents the find API response.
type findResponse struct {
	Code    string     `json:"cod"`
	Message string     `json:"message"`
	Count   int        `json:"count"`
	List    []Response `json:"list"`
}

// Find finds weather conditions for a search query.
func Find(apiKey string, q string) (*Response, error) {
	v := url.Values{}
	v.Set("q", q)
	v.Set("units", "metric")
	v.Set("lang", "en")
	v.Set("appid", apiKey)
	dest := fmt.Sprintf("https://api.openweathermap.org/data/2.5/find?%s", v.Encode())
	client, err := transport.Client(dest)
	if err != nil {
		return nil, err
	}
	resp, err := client.Get(dest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var r findResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	if r.Code != "200" {
		return nil, fmt.Errorf("%v: %v", r.Code, r.Message)
	}
	if len(r.List) == 0 {
		return nil, fmt.Errorf("not found")
	}
	return &r.List[0], nil
}
