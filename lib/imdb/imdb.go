// Package imdb implements Title find and information using AppEngine JSON API.
package imdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/StalkR/goircbot/lib/transport"
	"github.com/StalkR/imdb"
)

const appURL = "https://movie-db-api.appspot.com"

// GetRetry performs an HTTP GET with retries.
func GetRetry(rawurl string, retries int) (*http.Response, error) {
	client, err := transport.Client(rawurl)
	if err != nil {
		return nil, err
	}
	for i := 0; i < retries; i++ {
		resp, err := client.Get(rawurl)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode == 200 {
			return resp, nil
		}
		log.Print("imdb: GET status ", resp.StatusCode)
	}
	return nil, fmt.Errorf("imdb: GET failed")
}

// NewTitle obtains a Title ID with its information and returns a Title.
func NewTitle(id string) (*imdb.Title, error) {
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
	var title imdb.Title
	if err = json.Unmarshal(c, &title); err != nil {
		return nil, err
	}
	return &title, nil
}

// FindTitle searches a Title and returns a list of titles that matched.
func FindTitle(q string) ([]imdb.Title, error) {
	base := appURL + "/find"
	params := url.Values{}
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
	var titles []imdb.Title
	if err = json.Unmarshal(c, &titles); err != nil {
		return nil, err
	}
	return titles, nil
}
