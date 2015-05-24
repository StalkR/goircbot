// Package metalarchives is a library to search for bands on metal archives.
package metalarchives

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/StalkR/goircbot/lib/metal"
	"github.com/StalkR/goircbot/lib/transport"
)

const baseURL = "http://www.metal-archives.com/search/ajax-band-search/"

// Search finds bands by name.
func Search(name string) ([]metal.Band, error) {
	client, err := transport.Client(baseURL)
	if err != nil {
		return nil, err
	}
	u := url.Values{"query": {name}, "field": {"name"}}
	resp, err := client.Get(fmt.Sprintf("%s?%s", baseURL, u.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result struct {
		Bands [][]string `json:"aaData"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	var results []metal.Band
	for _, r := range result.Bands {
		if len(r) < 3 {
			continue
		}
		results = append(results, metal.Band{
			Name:    strings.TrimSpace(stripTags(r[0])),
			Genre:   r[1],
			Country: r[2],
		})
	}
	return results, nil
}

var htmlTagsRE = regexp.MustCompile(`<[^>]*>`)

func stripTags(s string) string {
	return htmlTagsRE.ReplaceAllString(s, "")
}
