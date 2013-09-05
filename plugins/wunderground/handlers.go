// Package wunderground implements a plugin to get weather conditions from wunderground.com.
package wunderground

import (
	"fmt"
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/wunderground"
)

func conditions(b *bot.Bot, e *bot.Event, APIKey string) {
	location := strings.TrimSpace(e.Args)
	if len(location) == 0 {
		return
	}
	r, err := wunderground.Conditions(APIKey, location)
	if err != nil {
		b.Conn.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	b.Conn.Privmsg(e.Target, r.String())
}

// Register registers the plugin with a bot.
func Register(b *bot.Bot, APIKey string) {
	b.AddCommand("weather", bot.Command{
		Help:    "get weather conditions from wunderground.com",
		Handler: func(b *bot.Bot, e *bot.Event) { conditions(b, e, APIKey) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
