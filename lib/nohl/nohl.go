// Package nohl implements methods to avoid highlighting nicks on a channel.
package nohl

import (
	"strings"

	"github.com/StalkR/goircbot/bot"
)

// Nick sanitizes a nick if it is present on the channel, otherwise unchanged.
func Nick(b bot.Bot, channel, nick string) string {
	_, on := b.Conn().StateTracker().IsOn(channel, nick)
	if !on {
		return nick
	}
	return nick[:len(nick)-1] + "*"
}

// String sanitizes words in string in case any is a nick present on the channel.
func String(b bot.Bot, channel, s string) string {
	ch := b.Conn().StateTracker().GetChannel(channel)
	if ch == nil {
		return s // channel not found (e.g. bot not on it)
	}
	for nick := range ch.Nicks {
		s = strings.Replace(s, nick, nick[:len(nick)-1]+"*", -1)
	}
	return s
}
