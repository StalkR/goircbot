// Package renick implements a plugin to get back a specific nick if it was taken.
package renick

import (
	"time"

	"github.com/StalkR/goircbot/bot"
	"github.com/fluffle/goirc/client"
)

// renick tries to renick until it succeeds.
func renick(b bot.Bot, nick string) {
	for _ = range time.Tick(time.Minute) {
		if !b.Connected() || b.Me().Nick == nick {
			break
		}
		b.Nick(nick)
	}
}

// Register runs the renick goroutine if we didn't get the nick we wanted.
func Register(b bot.Bot, nick string) {
	b.Conn().HandleFunc("connected",
		func(conn *client.Conn, line *client.Line) {
			if conn.Me().Nick != nick {
				go renick(b, nick)
			}
		})
}
