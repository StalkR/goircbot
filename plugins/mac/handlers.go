// Package mac implements a plugin to know manufacturer of a MAC (IEEE public OUI).
package mac

import (
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/ieeeoui"
)

func mac(e *bot.Event, r *ieeeoui.Resolver) {
	arg := strings.TrimSpace(e.Args)
	if len(arg) == 0 {
		return
	}
	manufacturer, err := r.Find(arg)
	if err != nil {
		e.Bot.Privmsg(e.Target, err.Error())
		return
	}
	e.Bot.Privmsg(e.Target, manufacturer)
}

// Register registers the plugin with a bot.
func Register(b bot.Bot) {
	resolver := ieeeoui.New()

	b.Commands().Add("mac", bot.Command{
		Help:    "find manufacturer of a MAC address (IEEE public OUI)",
		Handler: func(e *bot.Event) { mac(e, resolver) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
