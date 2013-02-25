// Package urbandictionary implements a plugin to get urban definition of words.
package urbandictionary

import (
	"fmt"
	bot "github.com/StalkR/goircbot"
	urbandictionary "github.com/StalkR/misc/urbandictionary"
	"strings"
)

func Urban(b *bot.Bot, e *bot.Event) {
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
		Help:    "get definition of word from urbandictionary",
		Handler: Urban,
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
