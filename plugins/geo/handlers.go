// Package geo implements a plugin to get geo location of an IP/host.
package geo

import (
	"fmt"
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/geo"
)

func locate(e *bot.Event) {
	addr := strings.TrimSpace(e.Args)
	if len(addr) == 0 {
		return
	}
	g, err := geo.Location(addr)
	if err != nil {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	e.Bot.Privmsg(e.Target, g.String())
}

// Register registers the plugin with a bot.
func Register(b bot.Bot) {
	b.Commands().Add("geo", bot.Command{
		Help:    "geo locate an IP/host",
		Handler: locate,
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
