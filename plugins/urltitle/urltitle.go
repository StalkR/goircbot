// Package urltitle implements a plugin to watch web URLs, fetch and display title.
package urltitle

import (
	"errors"
	"fmt"
	bot "github.com/StalkR/goircbot"
	"github.com/fluffle/goirc/client"
	"html"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// When matched, urltitle do not read line.
var silenceRegexp = "(^|\\s)tg(\\s|$)"

func timeoutDialer(d time.Duration) func(net, addr string) (c net.Conn, err error) {
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

func watchLine(b *bot.Bot, line *client.Line, ignoremap map[string]bool) {
	target := line.Args[0]
	if !strings.HasPrefix(target, "#") {
		return
	}
	if _, ignore := ignoremap[line.Nick]; ignore {
		return
	}
	text := line.Args[1]
	if m, err := regexp.Match(silenceRegexp, []byte(text)); err != nil || m {
		return
	}
	r, err := regexp.Compile("(?:^|\\s)(https?://[^\\s]+)")
	if err != nil {
		return
	}
	matches := r.FindSubmatch([]byte(text))
	if len(matches) < 2 {
		return
	}
	url := string(matches[1])
	title, err := Title(url)
	if err != nil {
		log.Println("urltitle:", err)
		return
	}
	if len(title) > 200 {
		title = title[:200]
	}
	b.Conn.Privmsg(target, fmt.Sprintf("%s :: %s", url, title))
}

// Register registers the plugin with a bot.
func Register(b *bot.Bot, ignore []string) {
	ignoremap := make(map[string]bool)
	for _, nick := range ignore {
		ignoremap[nick] = true
	}

	b.Conn.HandleFunc("privmsg",
		func(conn *client.Conn, line *client.Line) {
			watchLine(b, line, ignoremap)
		})
}
