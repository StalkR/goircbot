// Package weather implements a plugin to get weather conditions from OpenWeatherMap.
package weather

import (
	"fmt"
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/openweathermap"
)

func weather(e *bot.Event, apiKey string) {
	location := strings.TrimSpace(e.Args)
	if len(location) == 0 {
		return
	}
	r, err := openweathermap.Find(apiKey, location)
	if err != nil {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("error: %v", err))
		return
	}
	e.Bot.Privmsg(e.Target, r.String())
}

// Register registers the plugin with a bot.
func Register(b bot.Bot, apiKey string) {
	b.Commands().Add("weather", bot.Command{
		Help:    "get weather conditions from openweathermap.org",
		Handler: func(e *bot.Event) { weather(e, apiKey) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
