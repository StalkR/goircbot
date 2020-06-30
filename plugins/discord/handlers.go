// Package discord implements a plugin to bridge IRC with discord.
package discord

import (
	"fmt"
	"log"
	"strings"

	bridge "github.com/StalkR/discordgo-bridge"
	"github.com/StalkR/goircbot/bot"
	"github.com/fluffle/goirc/client"
)

// A Channel represents a discord channel.
type Channel struct {
	Channel string
	Webhook string
}

// Register registers the plugin with a bot.
// Config maps IRC channel to Discord channel.
func Register(b bot.Bot, token string, config map[string]*Channel) {
	channels := map[string]*bridge.Channel{}
	for chanirc, discord := range config {
		channels[chanirc] = bridge.NewChannel(discord.Channel, discord.Webhook, func(nick, text string) {
			toIRC(b, chanirc, nick, text)
		})
	}
	var list []*bridge.Channel
	for _, v := range channels {
		list = append(list, v)
	}

	d := bridge.NewBot(token, list...)
	if err := d.Start(); err != nil {
		log.Fatal(err)
	}
	// never closed

	b.Conn().HandleFunc("privmsg",
		func(conn *client.Conn, line *client.Line) { toDiscord(d, line, channels) })
}

func toIRC(b bot.Bot, channel, nick, text string) {
	if strings.HasPrefix(text, b.CommandPrefix()) {
		b.Privmsg(channel, fmt.Sprintf("Command sent from discord by %v", nick))
		b.Privmsg(channel, text)
		return
	}
	b.Privmsg(channel, fmt.Sprintf("<%v> %v", nick, text))
}

func toDiscord(d *bridge.Bot, line *client.Line, channels map[string]*bridge.Channel) {
	channel := line.Args[0]
	nick := line.Nick
	text := line.Args[1]
	c, ok := channels[channel]
	if !ok {
		return
	}
	c.Send(nick, text)
}
