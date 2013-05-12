// Package mac implements a plugin to know manufacturer of a MAC (IEEE public OUI).
package mac

import (
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/ieeeoui"
)

func mac(b *bot.Bot, e *bot.Event, r *ieeeoui.Resolver) {
	arg := strings.TrimSpace(e.Args)
	if len(arg) == 0 {
		return
	}
	manufacturer, err := r.Find(arg)
	if err != nil {
		b.Conn.Privmsg(e.Target, err.Error())
		return
	}
	b.Conn.Privmsg(e.Target, manufacturer)
}

// Register registers the plugin with a bot.
func Register(b *bot.Bot) {
	resolver := ieeeoui.New()

	b.AddCommand("mac", bot.Command{
		Help:    "find manufacturer of a MAC address (IEEE public OUI)",
		Handler: func(b *bot.Bot, e *bot.Event) { mac(b, e, resolver) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
