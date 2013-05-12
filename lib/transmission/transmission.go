// Package transmission implements a library to pull statistics from transmission.
package transmission

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

// Stats returns Statistics for a given transmission URL.
func Stats(url string) (*Statistics, error) {
	c, err := newConn(url)
	if err != nil {
		return nil, err
	}
	return c.stats()
}

// A Statistics holds generic stats of transmission.
type Statistics struct {
	DownloadSpeed, UploadSpeed                           int
	TorrentCount, ActiveTorrentCount, PausedTorrentcount int
	CurrentStats, CumulativeStats                        TotalStats
}

func (s *Statistics) String() string {
	return fmt.Sprintf("%v KB/s DL, %v KB/s UL, %v torrents (%v active, %v paused)",
		s.DownloadSpeed/1024, s.UploadSpeed/1024, s.TorrentCount,
		s.ActiveTorrentCount, s.PausedTorrentcount)
}

// A TotalStats holds total stats of transmission.
type TotalStats struct {
	DownloadedBytes, UploadedBytes          int
	FilesAdded, SecondsActive, SessionCount int
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

func (c *conn) sessionId() (string, error) {
	resp, err := c.client.Get(c.url + "/transmission/rpc")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	v, ok := resp.Header["X-Transmission-Session-Id"]
	if !ok || len(v) < 1 {
		return "", errors.New("transmission: sessionId not found")
	}
	return v[0], nil
}

func (c *conn) stats() (*Statistics, error) {
	sessId, err := c.sessionId()
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBufferString(`{"method":"session-stats"}`)
	req, err := http.NewRequest("POST", c.url+"/transmission/rpc", buf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Transmission-Session-Id", sessId)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var s sessionStats
	if err := json.Unmarshal(body, &s); err != nil {
		return nil, err
	}
	if s.Result != "success" {
		return nil, fmt.Errorf("transmission: no success: %s", s.Result)
	}
	return &s.Arguments, nil
}

type sessionStats struct {
	Arguments Statistics
	Result    string
}
