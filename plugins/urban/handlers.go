// Package urban implements a plugin to get Urban Dictionary definition of words.
package urban

import (
	"fmt"
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/urbandictionary"
)

func define(b *bot.Bot, e *bot.Event) {
	term := strings.TrimSpace(e.Args)
	if len(term) == 0 {
		return
	}
	r, err := urbandictionary.Define(term)
	if err != nil {
		b.Conn.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	b.Conn.Privmsg(e.Target, r.String())
}

// Register registers the plugin with a bot.
func Register(b *bot.Bot) {
	b.AddCommand("urban", bot.Command{
		Help:    "get definition of word from Urban Dictionary",
		Handler: define,
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
