// Package mldonkey implements a library to pull statistics from mldonkey.
package mldonkey

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

var (
	downRE  = regexp.MustCompile(`Down: ([\d.]+ .B/s) `)
	upRE    = regexp.MustCompile(`Up: ([\d.]+ .B/s) `)
	totalRE = regexp.MustCompile(`Total\((\d+)\): ([\d.]+.)/([\d.]+.) @`)
)

// Stats returns Statistics for a given mldonkey URL.
func Stats(url string) (*Statistics, error) {
	c, err := newConn(url)
	if err != nil {
		return nil, err
	}
	return c.stats()
}

// A Statistics holds generic stats of mldonkey.
type Statistics struct {
	DL, UL            string
	Count             int
	Downloaded, Total string
}

func (s *Statistics) String() string {
	return fmt.Sprintf("%v DL, %v UL, %v total (%v/%v downloaded)",
		s.DL, s.UL, s.Count, s.Downloaded, s.Total)
}

type conn struct {
	url    string
	client http.Client
}

func timeoutDialer(d time.Duration) func(net, addr string) (net.Conn, error) {
	return func(netw, addr string) (net.Conn, error) {
		return net.DialTimeout(netw, addr, d)
	}
}

func newConn(rawurl string) (*conn, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	// TODO(StalkR): add CAcert properly instead of InsecureSkipVerify.
	return &conn{
		url: rawurl,
		client: http.Client{
			Transport: &http.Transport{
				Dial:            timeoutDialer(3 * time.Second),
				TLSClientConfig: &tls.Config{ServerName: u.Host, InsecureSkipVerify: true},
			},
		},
	}, nil
}

func (c *conn) stats() (*Statistics, error) {
	resp, err := c.client.Get(c.url + "/submit?q=bw_stats")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bw_stats, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	m := downRE.FindSubmatch(bw_stats)
	if m == nil {
		return nil, errors.New("mldonkey: cannot parse download speed")
	}
	DL := string(m[1])
	m = upRE.FindSubmatch(bw_stats)
	if m == nil {
		return nil, errors.New("mldonkey: cannot parse upload speed")
	}
	UL := string(m[1])

	resp, err = c.client.Get(c.url + "/submit?q=vd")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	vd, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	m = totalRE.FindSubmatch(vd)
	if m == nil {
		return nil, errors.New("mldonkey: cannot parse total")
	}
	count, err := strconv.Atoi(string(m[1]))
	if err != nil {
		return nil, errors.New("mldonkey: cannot parse total int")
	}
	dled := string(m[2])
	total := string(m[3])

	return &Statistics{DL: DL, UL: UL, Count: count, Downloaded: dled, Total: total}, nil
}
