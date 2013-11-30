// Package wunderground implements a plugin to get weather conditions from wunderground.com.
package wunderground

import (
	"fmt"
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/wunderground"
)

func conditions(e *bot.Event, APIKey string) {
	location := strings.TrimSpace(e.Args)
	if len(location) == 0 {
		return
	}
	r, err := wunderground.Conditions(APIKey, location)
	if err != nil {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	e.Bot.Privmsg(e.Target, r.String())
}

// Register registers the plugin with a bot.
func Register(b bot.Bot, APIKey string) {
	b.Commands().Add("weather", bot.Command{
		Help:    "get weather conditions from wunderground.com",
		Handler: func(e *bot.Event) { conditions(e, APIKey) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
