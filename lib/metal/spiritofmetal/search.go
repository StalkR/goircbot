// Package spiritofmetal is a library to search for bands on spirit of metal.
package spiritofmetal

import (
	"fmt"
	"html"
	"io/ioutil"
	"net/url"
	"regexp"

	"github.com/StalkR/goircbot/lib/metal"
	"github.com/StalkR/goircbot/lib/transport"
)

const baseURL = "https://www.spirit-of-metal.com/liste_groupe.php"

var (
	resultRE      = regexp.MustCompile(`(?s)<div class="col-xs-9"><h3>([^<]+)</h3>(.*?) - (.*?) <`)
	nameCountryRE = regexp.MustCompile(`^(.*?) \([A-Z]+\)$`)
)

// Search finds bands by name.
func Search(name string) ([]metal.Band, error) {
	client, err := transport.Client(baseURL)
	if err != nil {
		return nil, err
	}
	u := url.Values{"recherche_groupe": {name}}
	resp, err := client.Get(fmt.Sprintf("%s?%s", baseURL, u.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var results []metal.Band
	for _, r := range resultRE.FindAllStringSubmatch(string(b), -1) {
		name := html.UnescapeString(r[1])
		// Some names are "Band (UK)", strip country to have "Band".
		nameCountry := nameCountryRE.FindStringSubmatch(name)
		if nameCountry != nil {
			name = nameCountry[1]
		}
		results = append(results, metal.Band{
			Name:    name,
			Genre:   html.UnescapeString(r[2]),
			Country: html.UnescapeString(r[3]),
		})
	}
	return results, nil
}
