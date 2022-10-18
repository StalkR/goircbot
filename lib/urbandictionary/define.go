// Package urbandictionary implements Urban Dictionary definition of words
// using their JSON API.
package urbandictionary

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// A Result represents an Urban Dictionary search result.
type Result struct {
	List []Definition `json:"list"`
}

// A Definition represents an Urban Dictionary definition.
type Definition struct {
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

func (r *Result) String() string {
	if len(r.List) == 0 {
		return "no result"
	}
	return r.List[0].String()
}

func (d *Definition) String() string {
	s := d.Definition
	s = strings.Replace(s, "\r", "", -1)
	s = strings.Replace(s, "\n", " ", -1)
	s = strings.Replace(s, "[", "", -1)
	s = strings.Replace(s, "]", "", -1)
	return fmt.Sprintf("%v: %v", d.Word, s)
}

// Define gets definition of term on Urban Dictionary and populates a Result.
func Define(term string) (*Result, error) {
	base := "http://api.urbandictionary.com/v0/define"
	params := url.Values{}
	params.Set("term", term)
	resp, err := http.DefaultClient.Get(fmt.Sprintf("%s?%s", base, params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var r Result
	if err = json.Unmarshal(contents, &r); err != nil {
		return nil, err
	}
	return &r, nil
}
