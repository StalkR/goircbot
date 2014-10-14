// Package transport provides an net/http Transport in a single place.
// HTTP(S) timeout can be set with a flag.
// It is a separate package so libs and plugins use it instead of net/http.
// That way, in case HTTP(S) needs to go through a proxy or else, only this
// file needs to be edited.
// This also allows us to use our own TLS package, also extensible.
package transport

import (
	"flag"
	"net"
	"net/http"
	neturl "net/url"
	"time"

	"github.com/StalkR/goircbot/lib/tls"
)

var timeout = flag.Duration("http_timeout", 10*time.Second, "Timeout for HTTP(S) connections.")

// ByURL returns an transport given an URL for use with a client.
// It extracts the hostname from URL and give it to TLS for SNI.
func ByURL(url string) (*http.Transport, error) {
	u, err := neturl.Parse(url)
	if err != nil {
		return nil, err
	}
	return &http.Transport{
		Dial:            timeoutDialer(*timeout),
		TLSClientConfig: tls.Config(u.Host),
	}, nil
}

// Client returns a client given an URL.
func Client(url string) (*http.Client, error) {
	trans, err := ByURL(url)
	if err != nil {
		return nil, err
	}
	return &http.Client{Transport: trans}, nil
}

func timeoutDialer(d time.Duration) func(net, addr string) (net.Conn, error) {
	return func(netw, addr string) (net.Conn, error) {
		return net.DialTimeout(netw, addr, d)
	}
}
