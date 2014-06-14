// Package time implements a plugin to get local time from wunderground.com.
package time

import (
	"fmt"
	"strings"
	"time"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/wunderground"
)

func showTime(e *bot.Event, APIKey string) {
	location := strings.TrimSpace(e.Args)
	if len(location) == 0 {
		return
	}
	r, err := wunderground.Conditions(APIKey, location)
	if err != nil {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	t, err := r.Time()
	if err != nil {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	e.Bot.Privmsg(e.Target, fmt.Sprintf("%s: %s", r.Display.Full, t.Format(time.RFC1123)))
}

// Register registers the plugin with a bot.
func Register(b bot.Bot, APIKey string) {
	b.Commands().Add("time", bot.Command{
		Help:    "get time from wunderground.com",
		Handler: func(e *bot.Event) { showTime(e, APIKey) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
