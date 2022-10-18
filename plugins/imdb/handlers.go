// Package imdb implements a plugin to search IMDb titles.
package imdb

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/imdb"
)

func search(e *bot.Event) {
	q := strings.TrimSpace(e.Args)
	if len(q) == 0 {
		return
	}
	titles, err := imdb.SearchTitle(http.DefaultClient, q)
	if err != nil {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	if len(titles) == 0 {
		e.Bot.Privmsg(e.Target, "No results found.")
		return
	}
	title, err := imdb.NewTitle(http.DefaultClient, titles[0].ID)
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
