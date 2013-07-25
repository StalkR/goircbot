// Package url implements a library to get a meaningful title of web URLs.
package url

import (
	"io/ioutil"

	"github.com/StalkR/goircbot/lib/transport"
)

// Title gets an URL and returns its title.
func Title(url string) (string, error) {
	client, err := transport.Client(url)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return ParseTitle(url, string(contents), Parsers)
}
