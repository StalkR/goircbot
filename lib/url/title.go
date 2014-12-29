// Package url implements a library to get a meaningful title of web URLs.
package url

import (
	"io"
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
	// Fetch maximum 1MB to avoid denial of service.
	lr := &io.LimitedReader{R: resp.Body, N: 1 << 20}
	contents, err := ioutil.ReadAll(lr)
	if err != nil {
		return "", err
	}
	return ParseTitle(url, string(contents), Parsers)
}
