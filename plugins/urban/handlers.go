// Package urban implements a plugin to get Urban Dictionary definition of words.
package urban

import (
	"fmt"
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/urbandictionary"
)

func define(e *bot.Event) {
	term := strings.TrimSpace(e.Args)
	if len(term) == 0 {
		return
	}
	r, err := urbandictionary.Define(term)
	if err != nil {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	def := r.String()
	if len(def) > 200 {
		def = def[:200] + "..."
	}
	e.Bot.Privmsg(e.Target, def)
}

// Register registers the plugin with a bot.
func Register(b bot.Bot) {
	b.Commands().Add("urban", bot.Command{
		Help:    "get definition of word from Urban Dictionary",
		Handler: define,
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
