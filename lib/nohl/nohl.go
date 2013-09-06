// Package nohl implements methods to avoid highlighting nicks on a channel.
package nohl

import (
	"github.com/StalkR/goircbot/bot"
)

// Nick sanitizes a nick if it is present on the channel, otherwise unchanged.
func Nick(b *bot.Bot, channel, nick string) string {
	// Is the bot on this channel? get state.Channel object.
	ch, on := b.Conn.Me().IsOnStr(channel)
	if !on {
		return nick
	}
	// Is this nick on this channel?
	if _, on := ch.IsOnStr(nick); !on {
		return nick
	}
	return nick[:len(nick)-1] + "*"
}
