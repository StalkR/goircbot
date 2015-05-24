// Package metalorgie is a library to search for bands on metalorgie.
package metalorgie

import (
	"fmt"
	"html"
	"io/ioutil"
	"net/url"
	"regexp"

	"github.com/StalkR/goircbot/lib/metal/band"
	"github.com/StalkR/goircbot/lib/transport"
)

const baseURL = "http://www.metalorgie.com/recherche"

var (
	sectionsRE = regexp.MustCompile(`(?s)<div class="fleft">(.*?)<div class="clear">`)
	nameRE     = regexp.MustCompile(`<a href="[^"]+" class="title">([^<]+)`)
	genreRE    = regexp.MustCompile(`Style :</span> ([^<]+)`)
	countryRE  = regexp.MustCompile(`Pays :</span> <a[^>]*>([^<]+)`)
)

// Search finds bands by name.
func Search(name string) ([]band.Band, error) {
	client, err := transport.Client(baseURL)
	if err != nil {
		return nil, err
	}
	u := url.Values{"q": {name}}
	resp, err := client.Get(fmt.Sprintf("%s?%s", baseURL, u.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var results []band.Band
	for _, r := range sectionsRE.FindAllString(string(b), -1) {
		name := nameRE.FindStringSubmatch(r)
		genre := genreRE.FindStringSubmatch(r)
		country := countryRE.FindStringSubmatch(r)
		if name == nil || genre == nil || country == nil {
			continue
		}
		results = append(results, band.Band{
			Name:    html.UnescapeString(name[1]),
			Genre:   html.UnescapeString(genre[1]),
			Country: html.UnescapeString(country[1]),
		})
	}
	return results, nil
}
