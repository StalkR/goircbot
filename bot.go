// Package goircbot implements an IRC bot using goirc library
// http://go.pkgdoc.org/github.com/fluffle/goirc
package goircbot

import (
	irc "github.com/fluffle/goirc/client"
	"log"
)

// Bot represents an IRC bot, with IRC client object, settings, commands and crons.
type Bot struct {
	Host      string
	Conn      *irc.Conn
	Quit      chan bool
	Reconnect bool
	commands  map[string]Command
	crons     map[string]Cron
}

// NewBot creates a new Bot with a set of parameters.
func NewBot(host string, ssl bool, nick, ident string, channels []string) (b *Bot) {
	b = &Bot{
		Host:      host,
		Conn:      irc.SimpleClient(nick, ident),
		Quit:      make(chan bool),
		Reconnect: true,
		commands:  make(map[string]Command),
		crons:     make(map[string]Cron),
	}
	if ssl {
		b.Conn.SSL = true
	}
	b.Conn.EnableStateTracking()

	// Join channels on connect and mark ourselves as a Bot.
	b.Conn.AddHandler("connected",
		func(conn *irc.Conn, line *irc.Line) {
			for _, channel := range channels {
				conn.Join(channel)
			}
			conn.Mode(conn.Me.Nick, "+B")
		})

	// Signal disconnect to Bot.Run so it can reconnect.
	b.Conn.AddHandler("disconnected",
		func(conn *irc.Conn, line *irc.Line) { b.Quit <- true })

	b.Conn.AddHandler("privmsg",
		func(conn *irc.Conn, line *irc.Line) { handleCommand(b, line) })

	b.AddCommand("help", Command{
		Help:    "show commands and detailed help",
		Handler: handleHelp,
		Pub:     true,
		Priv:    true,
		Hidden:  false})

	return b
}

// Run starts a configured Bot.
func (b *Bot) Run() {
	// Reconnect loop, unless we want to exit.
	for b.Reconnect {
		if err := b.Conn.Connect(b.Host); err != nil {
			log.Println("Connection error:", err)
		}

		// Wait on quit channel for a disconnect event.
		<-b.Quit
	}
}
