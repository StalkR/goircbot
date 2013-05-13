// Package url implements a library to get a meaningful title of web URLs.
package url

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/StalkR/goircbot/lib/tls"
)

func timeoutDialer(d time.Duration) func(net, addr string) (net.Conn, error) {
	return func(netw, addr string) (net.Conn, error) {
		return net.DialTimeout(netw, addr, d)
	}
}

// Title gets an URL and returns its title.
func Title(rawurl string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	client := &http.Client{
		Transport: &http.Transport{
			Dial:            timeoutDialer(3 * time.Second),
			TLSClientConfig: tls.Config(u.Host),
		},
	}
	resp, err := client.Get(rawurl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return ParseTitle(rawurl, string(contents), Parsers)
}
