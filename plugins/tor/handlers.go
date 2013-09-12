// Package tor implements a plugin to get info from a TOR node.
package tor

import (
	"fmt"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/tor"
)

func torInfo(b *bot.Bot, e *bot.Event, hostPort, key string) {
	i, err := tor.Info(hostPort, key)
	if err != nil {
		b.Conn.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	b.Conn.Privmsg(e.Target, i.String())
}

// Register registers the plugin with a bot.
func Register(b *bot.Bot, hostPort, key string) {
	b.AddCommand("tor", bot.Command{
		Help:    "get info on TOR node",
		Handler: func(b *bot.Bot, e *bot.Event) { torInfo(b, e, hostPort, key) },
		Pub:     true,
		Priv:    false,
		Hidden:  false})
}
