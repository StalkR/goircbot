package geo

import (
	"fmt"
	"strings"
)

type City struct {
	Confidence int               `json:"confidence,omitempty"`
	GeoNameId  int               `json:"geoname_id,omitempty"`
	Names      map[string]string `json:"names,omitempty"`
}

type Continent struct {
	Code      string            `json:"code,omitempty"`
	GeoNameId int               `json:"geoname_id,omitempty"`
	Names     map[string]string `json:"names,omitempty"`
}

type Country struct {
	Confidence int               `json:"confidence,omitempty"`
	GeoNameId  int               `json:"geoname_id,omitempty"`
	IsoCode    string            `json:"iso_code,omitempty"`
	Names      map[string]string `json:"names,omitempty"`
}

type Location struct {
	AccuracyRadius    string  `json:"accuracy_radius,omitempty"`
	AverageIncome     int     `json:"average_income,omitempty"`
	Latitude          float64 `json:"latitude,omitempty"`
	Longitude         float64 `json:"longitude,omitempty"`
	MetroCode         int     `json:"metro_code,omitempty"`
	PopulationDensity int     `json:"population_density,omitempty"`
	TimeZone          string  `json:"time_zone,omitempty"`
}

type Postal struct {
	Code       string `json:"code,omitempty"`
	Confidence int    `json:"confidence,omitempty"`
}

type RegisteredCountry struct {
	GeoNameId int               `json:"geoname_id,omitempty"`
	IsoCode   string            `json:"iso_code,omitempty"`
	Names     map[string]string `json:"names,omitempty"`
}

type RepresentedCountry struct {
	GeoNameId int               `json:"geoname_id,omitempty"`
	IsoCode   string            `json:"iso_code,omitempty"`
	Names     map[string]string `json:"names,omitempty"`
	Type      string            `json:"type,omitempty"`
}

type Subdivision struct {
	Confidence int               `json:"confidence,omitempty"`
	GeoNameId  int               `json:"geoname_id,omitempty"`
	IsoCode    string            `json:"iso_code,omitempty"`
	Names      map[string]string `json:"names,omitempty"`
}

type Traits struct {
	AS                  string `json:"autonomous_system_number,omitempty"`
	ASOrg               string `json:"autonomous_system_organization,omitempty"`
	Domain              string `json:"domain,omitempty"`
	IsAnonymousProxy    bool   `json:"is_anonymous_proxy,omitempty"`
	IsSatelliteProvider bool   `json:"is_satellite_provider,omitempty"`
	ISP                 string `json:"isp,omitempty"`
	IP                  string `json:"ip_address,omitempty"`
	Org                 string `json:"organization,omitempty"`
	UserType            string `json:"user_type,omitempty"`
}

type MaxMind struct {
	QueriesRemaining int `json:"queries_remaining,omitempty"`
}

type Response struct {
	City               City               `json:"city,omitempty"`
	Continent          Continent          `json:"continent,omitempty"`
	Country            Country            `json:"country,omitempty"`
	Location           Location           `json:"location,omitempty"`
	Postal             Postal             `json:"postal,omitempty"`
	RegisteredCountry  RegisteredCountry  `json:"registered_country,omitempty"`
	RepresentedCountry RepresentedCountry `json:"represented_country,omitempty"`
	Subdivisions       []Subdivision      `json:"subdivisions,omitempty"`
	Traits             Traits             `json:"traits,omitempty"`
	MaxMind            MaxMind            `json:"maxmind,omitempty"`
}

func (r *Response) String() string {
	var fields []string
	add := func(name string, value interface{}) {
		if value == "" {
			return
		}
		fields = append(fields, fmt.Sprintf("%v: %v", name, value))
	}
	add("City", r.City.Names["en"])
	add("Country", r.Country.Names["en"])
	add("Org", r.Traits.Org)
	add("ISP", r.Traits.ISP)
	add("IP", r.Traits.IP)
	add("Domain", r.Traits.Domain)
	add("Latitude", r.Location.Latitude)
	add("Longitude", r.Location.Longitude)
	return strings.Join(fields, ", ")
}
