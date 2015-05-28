// Package errors implements a plugin to get Linux & Windows error and status information.
package errors

import (
	"strings"

	"github.com/StalkR/goircbot/bot"
)

func handle(e *bot.Event, table []info) {
	arg := strings.TrimSpace(e.Args)
	if len(arg) == 0 {
		return
	}
	r, err := find(table, arg)
	if err != nil {
		e.Bot.Privmsg(e.Target, err.Error())
		return
	}
	e.Bot.Privmsg(e.Target, r.String())
}

// Register registers the plugin with a bot.
func Register(b bot.Bot) error {
	b.Commands().Add("error", bot.Command{
		Help:    "get Windows error code information",
		Handler: func(e *bot.Event) { handle(e, winerrors) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})
	b.Commands().Add("status", bot.Command{
		Help:    "get Windows status code information",
		Handler: func(e *bot.Event) { handle(e, ntstatus) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})
	b.Commands().Add("errno", bot.Command{
		Help:    "get Linux error code information",
		Handler: func(e *bot.Event) { handle(e, errnos) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})
	return nil
}
