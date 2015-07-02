// Package invite implements a plugin to join channels on invite.
package invite

import (
	"github.com/StalkR/goircbot/bot"
	"github.com/fluffle/goirc/client"
)

// Register registers the plugin with a bot.
func Register(b bot.Bot) {
	b.Conn().HandleFunc("invite",
		func(conn *client.Conn, line *client.Line) {
			if len(line.Args) < 2 {
				return
			}
			b.Join(line.Args[1])
		})
}
