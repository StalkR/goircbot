// Package up implements a plugin to check if a web host is up or down.
package up

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/StalkR/goircbot/bot"
)

// Probe gets an URL and returns a boolean if it worked within imparted time.
func Probe(url string) bool {
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return true
}

var hostRE = regexp.MustCompile(`^[\w._-]+$`)

func up(e *bot.Event) {
	arg := strings.TrimSpace(e.Args)
	if len(arg) == 0 {
		return
	}
	target := arg
	if hostRE.MatchString(arg) {
		target = fmt.Sprintf("http://%s", arg)
	}
	u, err := url.Parse(target)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("%s must be either a hostname or an http(s) URL.", arg))
		return
	}
	if !Probe(target) {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("%s is down from here.", target))
		return
	}
	e.Bot.Privmsg(e.Target, fmt.Sprintf("%s is up from here.", target))
}

// Register registers the plugin with a bot.
func Register(b bot.Bot) {
	b.Commands().Add("up", bot.Command{
		Help:    "check if a web host is up or down",
		Handler: up,
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
