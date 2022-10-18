// Package up implements a plugin to check if a web host is up or down.
package up

import (
	"fmt"
	"net/http"
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

func up(e *bot.Event) {
	arg := strings.TrimSpace(e.Args)
	if len(arg) == 0 {
		return
	}
	matched, err := regexp.Match(`^[\w._-]+$`, []byte(arg))
	if err != nil || !matched {
		return
	}
	url := fmt.Sprintf("http://%s", arg)
	if !Probe(url) {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("%s is down from here.", url))
		return
	}
	e.Bot.Privmsg(e.Target, fmt.Sprintf("%s is up from here.", url))
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
