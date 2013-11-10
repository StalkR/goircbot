// Package imdb implements a plugin to search IMDb titles.
package imdb

import (
	"fmt"
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/imdb"
)

func search(b *bot.Bot, e *bot.Event) {
	q := strings.TrimSpace(e.Args)
	if len(q) == 0 {
		return
	}
	titles, err := imdb.FindTitle(q)
	if err != nil {
		b.Conn.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	if len(titles) == 0 {
		b.Conn.Privmsg(e.Target, "No results found.")
		return
	}
	title, err := imdb.NewTitle(titles[0].ID)
	if err != nil {
		b.Conn.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	b.Conn.Privmsg(e.Target, title.String())
}

// Register registers the plugin with a bot.
func Register(b *bot.Bot) {
	b.AddCommand("imdb", bot.Command{
		Help:    "imdb <title> - search a Title on IMDb",
		Handler: search,
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
