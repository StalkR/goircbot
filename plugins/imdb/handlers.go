// Package imdb implements a plugin to search IMDb titles.
package imdb

import (
	"fmt"
	bot "github.com/StalkR/goircbot"
	"github.com/StalkR/goircbot/lib/imdb"
	"strings"
)

func Imdb(b *bot.Bot, e *bot.Event) {
	q := strings.TrimSpace(e.Args)
	if len(q) == 0 {
		return
	}
	r, err := imdb.FindTitle(q)
	if err != nil {
		b.Conn.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	if len(r) == 0 {
		b.Conn.Privmsg(e.Target, "No results found.")
		return
	}
	t, err := imdb.NewTitle(r[0].Id)
	if err != nil {
		b.Conn.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	b.Conn.Privmsg(e.Target, t.String())
}

// Register registers the plugin with a bot.
func Register(b *bot.Bot) {
	b.AddCommand("imdb", bot.Command{
		Help:    "imdb <title> - search a Title on IMDb",
		Handler: Imdb,
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
