// Package spiritofmetal is a library to search for bands on spirit of metal.
package spiritofmetal

import (
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"net/url"
	"regexp"

	"github.com/StalkR/goircbot/lib/metal/band"
	"github.com/StalkR/goircbot/lib/transport"
)

const baseURL = "http://www.spirit-of-metal.com/find.php"

var (
	sectionRE     = regexp.MustCompile(`(?s)<div[^>]*><h1[^>]*>Results in the bands section(.*?)</div></div>`)
	resultRE      = regexp.MustCompile(`(?s)<ul[^>]*>\s*<a[^>]*>([^<]+)</a>\s*\(([^-]+)-([^)]+)\)\s*</ul>`)
	nameCountryRE = regexp.MustCompile(`^(.*?) \([A-Z]+\)$`)
)

// Search finds bands by name.
func Search(name string) ([]band.Band, error) {
	client, err := transport.Client(baseURL)
	if err != nil {
		return nil, err
	}
	u := url.Values{"search": {"all"}, "nom": {name}}
	resp, err := client.Get(fmt.Sprintf("%s?%s", baseURL, u.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	s := string(b)
	section := sectionRE.FindString(s)
	if section == "" {
		return nil, errors.New("spiritofmetal: results section not found")
	}
	var results []band.Band
	for _, r := range resultRE.FindAllStringSubmatch(section, -1) {
		name := html.UnescapeString(r[1])
		// Some names are "Band (UK)", strip country to have "Band".
		nameCountry := nameCountryRE.FindStringSubmatch(name)
		if nameCountry != nil {
			name = nameCountry[1]
		}
		results = append(results, band.Band{
			Name:    name,
			Genre:   html.UnescapeString(r[2]),
			Country: html.UnescapeString(r[3]),
		})
	}
	return results, nil
}
