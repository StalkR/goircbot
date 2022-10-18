// Package mldonkey implements a library to talk to MLDonkey.
package mldonkey

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var (
	downRE  = regexp.MustCompile(`Down: ([\d.]+ .B/s) `)
	upRE    = regexp.MustCompile(`Up: ([\d.]+ .B/s) `)
	totalRE = regexp.MustCompile(`Total\((\d+)\): ([\d.]+.?)/([\d.]+.?) @`)
	linkRE  = regexp.MustCompile(`^ed2k://`)
)

// A Statistics holds generic stats of MLDonkey.
type Statistics struct {
	DL, UL            string
	Count             int
	Downloaded, Total string
}

func (s *Statistics) String() string {
	if s.Count == 0 {
		return fmt.Sprintf("%v DL, %v UL, 0 total", s.DL, s.UL)
	}
	return fmt.Sprintf("%v DL, %v UL, %v total (%v/%v downloaded)",
		s.DL, s.UL, s.Count, s.Downloaded, s.Total)
}

// A Conn represents a connection to MLDonkey.
type Conn struct {
	url string
}

// New prepares an MLDonkey connection by returning a *Conn.
func New(serverURL string) (*Conn, error) {
	return &Conn{url: serverURL}, nil
}

// Stats returns current statistics (speed, total downloads, etc.).
func (c *Conn) Stats() (*Statistics, error) {
	resp, err := http.DefaultClient.Get(c.url + "/submit?q=bw_stats")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bwStats, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	m := downRE.FindSubmatch(bwStats)
	if m == nil {
		return nil, errors.New("mldonkey: cannot parse download speed")
	}
	DL := string(m[1])
	m = upRE.FindSubmatch(bwStats)
	if m == nil {
		return nil, errors.New("mldonkey: cannot parse upload speed")
	}
	UL := string(m[1])

	resp, err = http.DefaultClient.Get(c.url + "/submit?q=vd")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	vd, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	m = totalRE.FindSubmatch(vd)
	if m == nil { // No current download.
		return &Statistics{DL: DL, UL: UL}, nil
	}
	count, err := strconv.Atoi(string(m[1]))
	if err != nil {
		return nil, errors.New("mldonkey: cannot parse total int")
	}
	dled := appendBytesSuffix(string(m[2]))
	total := appendBytesSuffix(string(m[3]))

	return &Statistics{DL: DL, UL: UL, Count: count, Downloaded: dled, Total: total}, nil
}

// appendBytesSuffix appends B suffix if it is a number.
// When large enough, mldk appends suffix like KB, MB, but nothing for bytes.
// Appending B so size cannot be confused with a number.
func appendBytesSuffix(n string) string {
	if _, err := strconv.Atoi(n); err == nil {
		return n + "B"
	}
	return n
}

// Add adds a link by URL.
func (c *Conn) Add(link string) error {
	if !linkRE.MatchString(link) {
		return errors.New("mldonkey: invalid link")
	}
	params := url.Values{}
	params.Set("q", link)
	resp, err := http.DefaultClient.Get(c.url + "/submit?" + params.Encode())
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
