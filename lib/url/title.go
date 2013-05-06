// Package url implements a module to get a meaningful title of web URLs.
package url

import (
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func timeoutDialer(d time.Duration) func(net, addr string) (net.Conn, error) {
	return func(netw, addr string) (net.Conn, error) {
		return net.DialTimeout(netw, addr, d)
	}
}

// Title gets an URL and returns its title.
func Title(url string) (string, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: timeoutDialer(3 * time.Second),
		},
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
	return ParseTitle(url, string(contents))
}
