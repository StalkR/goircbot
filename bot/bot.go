// Package bot implements an IRC bot with plugins.
package bot

import (
	"log"
	"strings"
	"time"

	"github.com/StalkR/goircbot/lib/tls"
	"github.com/fluffle/goirc/client"
	"github.com/fluffle/goirc/state"
)

// Bot represents an IRC bot, with IRC client object, settings, commands.
type Bot struct {
	Host      string
	Conn      *client.Conn
	Quit      chan bool
	Reconnect bool
	commands  map[string]Command
}

// NewBot creates a new Bot with a set of parameters.
func NewBot(host string, ssl bool, nick, ident string, channels []string) *Bot {
	hostPort := strings.SplitN(host, ":", 2)
	cfg := &client.Config{
		Me:          state.NewNick(nick),
		NewNick:     func(s string) string { return s + "_" },
		PingFreq:    3 * time.Minute,
		QuitMessage: "I have to go.",
		Server:      host,
		SSL:         ssl,
		SSLConfig:   tls.Config(hostPort[0]),
		Version:     "Powered by GoIRCBot",
		Recover:     (*client.Conn).LogPanic,
		SplitLen:    450,
	}
	cfg.Me.Ident = ident
	cfg.Me.Name = nick

	conn := client.Client(cfg)
	b := &Bot{
		Host:      host,
		Conn:      conn,
		Quit:      make(chan bool),
		Reconnect: true,
		commands:  make(map[string]Command),
	}

	b.Conn.EnableStateTracking()

	// Join channels on connect and mark ourselves as a Bot.
	b.Conn.HandleFunc("connected",
		func(conn *client.Conn, line *client.Line) {
			for _, channel := range channels {
				conn.Join(channel)
			}
			conn.Mode(conn.Me().Nick, "+B")
		})

	// Signal disconnect to Bot.Run so it can reconnect.
	b.Conn.HandleFunc("disconnected",
		func(conn *client.Conn, line *client.Line) { b.Quit <- true })

	b.Conn.HandleFunc("privmsg",
		func(conn *client.Conn, line *client.Line) { handleCommand(b, line) })

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
		if err := b.Conn.Connect(); err != nil {
			log.Println("Connection error:", err, "- reconnecting in 1min")
			time.Sleep(time.Minute)
			continue
		}

		// Wait on quit channel for a disconnect event.
		<-b.Quit
	}
}
