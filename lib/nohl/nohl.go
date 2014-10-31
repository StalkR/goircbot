// Package nohl implements methods to avoid highlighting nicks on a channel.
package nohl

import (
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/fluffle/goirc/state"
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
	var c *state.Channel
	for _, e := range b.Me().Channels() {
		if e.Name == channel {
			c = e
		}
	}
	if c == nil {
		return s // channel not found (e.g. bot not on it)
	}
	for _, n := range c.Nicks() {
		s = strings.Replace(s, n.Nick, n.Nick[:len(n.Nick)-1]+"*", -1)
	}
	return s
}
