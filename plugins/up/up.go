// Package up implements a plugin to check if a web host is up or down.
package up

import (
	"fmt"
	bot "github.com/StalkR/goircbot"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func timeoutDialer(d time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		return net.DialTimeout(netw, addr, d)
	}
}

// Probe gets an URL and returns a boolean if it worked within 3 seconds.
func Probe(url string) bool {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: timeoutDialer(3 * time.Second),
		},
	}
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return true
}

func Up(b *bot.Bot, e *bot.Event) {
	arg := strings.TrimSpace(e.Args)
	if len(arg) == 0 {
		return
	}
	matched, err := regexp.Match("^[\\w._-]+$", []byte(arg))
	if err != nil || !matched {
		return
	}
	url := fmt.Sprintf("http://%s", arg)
	if !Probe(url) {
		b.Conn.Privmsg(e.Target, fmt.Sprintf("%s is down from here.", url))
		return
	}
	b.Conn.Privmsg(e.Target, fmt.Sprintf("%s is up from here.", url))
}

// Register registers the plugin with a bot.
func Register(b *bot.Bot) {
	b.AddCommand("up", bot.Command{
		Help:    "check if a web host is up or down",
		Handler: Up,
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
