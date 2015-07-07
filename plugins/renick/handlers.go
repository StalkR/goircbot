// Package renick implements a plugin to get back a specific nick if it was taken.
package renick

import (
	"strings"
	"time"

	"github.com/StalkR/goircbot/bot"
	"github.com/fluffle/goirc/client"
)

const (
	replyISupport = "005"
	replyLogOff   = "601"
	replyNowOff   = "605"
)

// handleNickFree grabs desiredNick.
func handleNickFree(b bot.Bot, desiredNick string, line *client.Line) {
	if line.Args[1] == desiredNick && b.Me().Nick != desiredNick {
		b.Nick(desiredNick)
	}
}

// handleISupport sets a WATCH on desiredNick if it's supported by the server.
func handleISupport(b bot.Bot, desiredNick string, line *client.Line) {
	if len(line.Args) < 2 {
		return
	}

	// First arg is our name, last arg is "are supported by this server".
	for _, support := range line.Args[1 : len(line.Args)-1] {
		parts := strings.SplitN(support, "=", 2)
		if parts[0] == "WATCH" {
			b.Conn().Raw("WATCH +" + desiredNick)
			return
		}
	}
}

// pollNick tries to grab desiredNick every minute.
func pollNick(b bot.Bot, desiredNick string) {
	for range time.Tick(time.Minute) {
		if b.Connected() && b.Me().Nick != desiredNick {
			b.Nick(desiredNick)
		}
	}
}

// Register registers the plugin with a bot.
func Register(b bot.Bot, desiredNick string) {
	b.HandleFunc(replyISupport, func(conn *client.Conn, line *client.Line) {
		handleISupport(b, desiredNick, line)
	})
	b.HandleFunc(replyLogOff, func(conn *client.Conn, line *client.Line) {
		handleNickFree(b, desiredNick, line)
	})
	b.HandleFunc(replyNowOff, func(conn *client.Conn, line *client.Line) {
		handleNickFree(b, desiredNick, line)
	})

	go pollNick(b, desiredNick)
}
