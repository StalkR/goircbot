// Package cdecl implements a plugin to explain C declarations using http://cdecl.org.
package cdecl

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/StalkR/goircbot/bot"
)

const queryURL = "http://cdecl.org/query.php"

func explain(s string) (string, error) {
	resp, err := http.DefaultClient.PostForm(queryURL, url.Values{"q": {s}})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func cdecl(e *bot.Event) {
	c := strings.TrimSpace(e.Args)
	if len(c) == 0 {
		e.Bot.Privmsg(e.Target, "Usage: cdecl <C declaration>")
		return
	}
	s, err := explain(c)
	if err != nil {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	e.Bot.Privmsg(e.Target, s)
}

// Register registers the plugin with a bot.
func Register(b bot.Bot) {
	b.Commands().Add("cdecl", bot.Command{
		Help:    "explain a C declaration using http://cdecl.org",
		Handler: cdecl,
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
