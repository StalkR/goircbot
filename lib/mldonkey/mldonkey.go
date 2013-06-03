// Package mldonkey implements a library to talk to MLDonkey.
package mldonkey

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/StalkR/goircbot/lib/tls"
)

var (
	downRE  = regexp.MustCompile(`Down: ([\d.]+ .B/s) `)
	upRE    = regexp.MustCompile(`Up: ([\d.]+ .B/s) `)
	totalRE = regexp.MustCompile(`Total\((\d+)\): ([\d.]+.)/([\d.]+.) @`)
	linkRE  = regexp.MustCompile(`^ed2k://`)
)

// A Statistics holds generic stats of MLDonkey.
type Statistics struct {
	DL, UL            string
	Count             int
	Downloaded, Total string
}

func (s *Statistics) String() string {
	return fmt.Sprintf("%v DL, %v UL, %v total (%v/%v downloaded)",
		s.DL, s.UL, s.Count, s.Downloaded, s.Total)
}

// A Conn represents a connection to MLDonkey.
type Conn struct {
	url    string
	client http.Client
}

func timeoutDialer(d time.Duration) func(net, addr string) (net.Conn, error) {
	return func(netw, addr string) (net.Conn, error) {
		return net.DialTimeout(netw, addr, d)
	}
}

// New prepares an MLDonkey connection by returning a *Conn.
func New(rawurl string) (*Conn, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	return &Conn{
		url: rawurl,
		client: http.Client{
			Transport: &http.Transport{
				Dial:            timeoutDialer(5 * time.Second),
				TLSClientConfig: tls.Config(u.Host),
			},
		},
	}, nil
}

// Stats returns current statistics (speed, total downloads, etc.).
func (c *Conn) Stats() (*Statistics, error) {
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

// Add adds a link by URL.
func (c *Conn) Add(link string) error {
	if !linkRE.MatchString(link) {
		return errors.New("mldonkey: invalid link")
	}
	params := url.Values{}
	params.Set("q", link)
	resp, err := c.client.Get(c.url + "/submit?" + params.Encode())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	page, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if !strings.Contains(string(page), "Added link") {
		fmt.Println(string(page))
		return fmt.Errorf("mldonkey: no result")
	}
	return nil
}
