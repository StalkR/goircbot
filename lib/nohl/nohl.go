// Package nohl implements methods to avoid highlighting nicks on a channel.
package nohl

import (
	"strings"

	"github.com/StalkR/goircbot/bot"
)

// Nick sanitizes a nick if it is present on the channel, otherwise unchanged.
func Nick(b bot.Bot, channel, nick string) string {
	// Is the bot on this channel? get state.Channel object.
	ch, on := b.Me().IsOnStr(channel)
	if !on {
		return nick
	}
	// Is this nick on this channel?
	if _, on := ch.IsOnStr(nick); !on {
		return nick
	}
	return nick[:len(nick)-1] + "*"
}

// String sanitizes words in string in case any is a nick present on the channel.
func String(b bot.Bot, channel, s string) string {
	a := strings.Split(s, " ")
	for i, w := range a {
		a[i] = Nick(b, channel, w)
	}
	return strings.Join(a, " ")
}
