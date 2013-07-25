// Package imdb implements Title find and information using AppEngine JSON API.
package imdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/StalkR/goircbot/lib/transport"
)

const appURL = "https://movie-db-api.appspot.com"

type Result struct {
	Id, Name string
	Year     int
}

type Title struct {
	Id, Name, Type, Rating, Duration, Description, Poster string
	Year, Year_production, Year_release                   int
	Aka, Genres, Languages, Nationalities                 []string
	Directors, Writers, Actors                            []Name
}

type Name struct {
	Id, Name string
}

func (t *Title) String() string {
	var infos []string
	name := t.Name
	if t.Year != 0 {
		name = fmt.Sprintf("%s (%d)", name, t.Year)
	}
	infos = append(infos, name)
	if len(t.Genres) > 0 {
		max := len(t.Genres)
		if max > 3 {
			max = 3
		}
		infos = append(infos, strings.Join(t.Genres[:max], ", "))
	}
	if len(t.Directors) > 0 {
		max := len(t.Directors)
		if max > 3 {
			max = 3
		}
		var directors []string
		for _, director := range t.Directors {
			directors = append(directors, director.String())
		}
		infos = append(infos, strings.Join(directors, ", "))
	}
	if len(t.Actors) > 0 {
		max := len(t.Actors)
		if max > 3 {
			max = 3
		}
		var actors []string
		for _, actor := range t.Actors[:max] {
			actors = append(actors, actor.String())
		}
		infos = append(infos, strings.Join(actors, ", "))
	}
	if t.Duration != "" {
		infos = append(infos, t.Duration)
	}
	if t.Rating != "" {
		infos = append(infos, t.Rating)
	}
	infos = append(infos, fmt.Sprintf("http://www.imdb.com/title/%s", t.Id))
	return strings.Join(infos, " - ")
}

func (n *Name) String() string {
	return n.Name
}

// GetRetry performs an HTTP GET with retries.
func GetRetry(rawurl string, retries int) (*http.Response, error) {
	client, err := transport.Client(rawurl)
	if err != nil {
		return nil, err
	}
	var resp *http.Response
	for i := 0; i < retries; i++ {
		resp, err := client.Get(rawurl)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode == 200 {
			return resp, nil
		}
		log.Print("imdb: get error, status ", resp.StatusCode)
	}
	return nil, fmt.Errorf("imdb: get error, status: %v", resp.StatusCode)
}

// Decode decodes json data from app.
func Decode(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	// Go < 1.1 do not accept mismatched null so just skip this error.
	// See https://code.google.com/p/go/issues/detail?id=2540
	if err != nil && !strings.Contains(fmt.Sprintf("%s", err), "cannot unmarshal null") {
		log.Print("imdb: decode error: ", fmt.Sprintf("%v", string(data)))
		return err
	}
	return nil
}

// NewTitle obtains a Title ID with its information and returns a Title.
func NewTitle(id string) (*Title, error) {
	base := appURL + "/title"
	resp, err := GetRetry(fmt.Sprintf("%s/%s", base, id), 3)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	c, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	t := &Title{}
	if err = Decode(c, t); err != nil {
		return nil, err
	}
	return t, nil
}

// FindTitle searches a Title and returns a list of Result.
func FindTitle(q string) ([]Result, error) {
	base := appURL + "/find"
	params := url.Values{}
	params.Set("s", "tt")
	params.Set("q", q)
	resp, err := GetRetry(fmt.Sprintf("%s?%s", base, params.Encode()), 3)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	c, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	r := make([]Result, 0)
	if err = Decode(c, &r); err != nil {
		return nil, err
	}
	return r, nil
}
