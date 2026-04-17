// Package metalorgie is a library to search for bands on metalorgie.
package metalorgie

import (
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/StalkR/goircbot/lib/metal"
)

const baseURL = "https://www.metalorgie.com/recherche"

var (
	resultsRE = regexp.MustCompile(`(?s)<h1 class="title">Groupes</h1>(.*?)<h1`)
	bandRE    = regexp.MustCompile(`(?s)<div class="column[^"]*">(.*?)</div>\s*</div>\s*</div>`)
	nameRE    = regexp.MustCompile(`<a class="post__content__band" href="[^"]+">([^<]+)`)
	genresRE  = regexp.MustCompile(`<span class="tag">([^<]+)`)
	countryRE = regexp.MustCompile(`<div class="post__content__flag flag-icon-background flag-icon-([^"]*)"></div>`)
)

// Search finds bands by name.
func Search(name string) ([]metal.Band, error) {
	u := url.Values{"q": {name}}
	resp, err := http.DefaultClient.Get(fmt.Sprintf("%s?%s", baseURL, u.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("metalorgie: status %v; body: %v", resp.Status, string(b))
	}
	results := resultsRE.FindStringSubmatch(string(b))
	if results == nil {
		return nil, fmt.Errorf("metalorgies: search results section not found")
	}
	var bands []metal.Band
	for _, section := range bandRE.FindAllString(results[1], -1) {
		name := nameRE.FindStringSubmatch(section)
		if name == nil {
			return nil, fmt.Errorf("metalorgies: empty name")
		}
		var genres []string
		for _, tag := range genresRE.FindAllStringSubmatch(section, -1) {
			genres = append(genres, html.UnescapeString(tag[1]))
		}
		if genres == nil {
			return nil, fmt.Errorf("metalorgies: no genres?")
		}
		country := countryRE.FindStringSubmatch(section)
		if country == nil {
			return nil, fmt.Errorf("metalorgies: empty country?")
		}
		bands = append(bands, metal.Band{
			Name:    html.UnescapeString(name[1]),
			Genre:   strings.Join(genres, " / "),
			Country: strings.ToUpper(country[1]),
		})
	}
	return bands, nil
}
