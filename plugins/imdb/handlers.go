// Package imdb implements a plugin to search IMDb titles.
package imdb

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/imdb"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36"

var client = &http.Client{
	Transport: &customTransport{http.DefaultTransport},
}

type customTransport struct {
	http.RoundTripper
}

func (e *customTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Accept-Language", "en") // avoid IP-based language detection
	r.Header.Set("User-Agent", userAgent)
	return e.RoundTripper.RoundTrip(r)
}

func search(e *bot.Event) {
	q := strings.TrimSpace(e.Args)
	if len(q) == 0 {
		return
	}
	titles, err := imdb.SearchTitle(client, q)
	if err != nil {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	if len(titles) == 0 {
		e.Bot.Privmsg(e.Target, "No results found.")
		return
	}
	title, err := imdb.NewTitle(client, titles[0].ID)
	if err != nil {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	e.Bot.Privmsg(e.Target, title.String())
}

// Register registers the plugin with a bot.
func Register(b bot.Bot) {
	b.Commands().Add("imdb", bot.Command{
		Help:    "imdb <title> - search a Title on IMDb",
		Handler: search,
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
