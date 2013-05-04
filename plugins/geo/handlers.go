// Package geo implements a plugin to get geo location of an IP/host.
package geo

import (
	"fmt"
	bot "github.com/StalkR/goircbot"
	"github.com/StalkR/goircbot/lib/geo"
	"strings"
)

func Geo(b *bot.Bot, e *bot.Event) {
	addr := strings.TrimSpace(e.Args)
	if len(addr) == 0 {
		return
	}
	g, err := geo.Location(addr)
	if err != nil {
		b.Conn.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	b.Conn.Privmsg(e.Target, g.String())
}

// Register registers the plugin with a bot.
func Register(b *bot.Bot) {
	b.AddCommand("geo", bot.Command{
		Help:    "geo locate an IP/host",
		Handler: Geo,
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
