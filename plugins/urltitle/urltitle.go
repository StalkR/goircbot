package urltitle

import (
	"errors"
	"html"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// When matched, urltitle do not read line.
var silenceRegexp = "(^|\\s)tg(\\s|$)"

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
	r, err := regexp.Compile("<title[^>]*>([^<]+)<")
	if err != nil {
		return "", err
	}
	matches := r.FindSubmatch([]byte(contents))
	if len(matches) < 2 {
		return "", errors.New("no title found in page")
	}
	s := string(matches[1])
	s = html.UnescapeString(s)
	r, err = regexp.Compile("\\s+")
	if err != nil {
		return "", err
	}
	s = r.ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)
	return s, nil
}
