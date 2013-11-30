// Package tor implements a plugin to get info from a TOR node.
package tor

import (
	"fmt"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/tor"
)

func torInfo(e *bot.Event, hostPort, key string) {
	i, err := tor.Info(hostPort, key)
	if err != nil {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	e.Bot.Privmsg(e.Target, i.String())
}

// Register registers the plugin with a bot.
func Register(b bot.Bot, hostPort, key string) {
	b.Commands().Add("tor", bot.Command{
		Help:    "get info on TOR node",
		Handler: func(e *bot.Event) { torInfo(e, hostPort, key) },
		Pub:     true,
		Priv:    false,
		Hidden:  false})
}
