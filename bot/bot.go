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
type Bot interface {
	Run()                       // Run bot, reconnect if disconnect.
	Quit(msg string)            // Quit bot from IRC with a msg.
	Commands() *Commands        // For plugins to Add/Del commands.
	Action(t, msg string)       // Shortcut to Conn().Action()
	Connected() bool            // Shortcut to Conn().Connected()
	Me() *state.Nick            // Shortcut to Conn().Me()
	Mode(t string, m ...string) // Shortcut to Conn().Mode()
	Notice(t, msg string)       // Shortcut to Conn().Notice()
	Privmsg(t, msg string)      // Shortcut to Conn().Privmsg()
	Conn() *client.Conn         // Conn returns the underlying goirc client connection.
}

// NewBot creates a new Bot implementation with a set of parameters.
func NewBot(host string, ssl bool, nick, ident string, channels []string) Bot {
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
	conn.EnableStateTracking()
	b := &BotImpl{
		conn:      conn,
		reconnect: true,
		quit:      make(chan bool),
		commands:  NewCommands(),
		channels:  channels,
	}

	// Join channels on connect and mark ourselves as a Bot.
	conn.HandleFunc("connected",
		func(conn *client.Conn, line *client.Line) {
			for _, channel := range b.channels {
				conn.Join(channel)
			}
			conn.Mode(conn.Me().Nick, "+B")
		})

	// Signal disconnect to Bot.Run so it can reconnect.
	conn.HandleFunc("disconnected",
		func(conn *client.Conn, line *client.Line) {
			channels := conn.Me().Channels()
			b.channels = make([]string, 0, len(channels))
			for _, ch := range channels {
				b.channels = append(b.channels, ch.Name)
			}
			b.quit <- true
		})

	conn.HandleFunc("privmsg",
		func(conn *client.Conn, line *client.Line) { b.commands.Handle(b, line) })

	b.commands.Add("help", Command{
		Help:    "show commands or detailed help",
		Handler: func(e *Event) { b.commands.Help(e) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})

	return b
}

// BotImpl implements Bot.
type BotImpl struct {
	conn      *client.Conn
	reconnect bool
	quit      chan bool
	commands  *Commands
	channels  []string
}

// Run starts the Bot by connecting it to IRC. It automatically reconnects.
func (b *BotImpl) Run() {
	for b.reconnect {
		if err := b.Conn().Connect(); err != nil {
			log.Println("Connection error:", err, "- reconnecting in 1 minute")
			time.Sleep(time.Minute)
			continue
		}

		// Wait on quit channel for a disconnect event.
		<-b.quit
	}
}

// Quit quits the bot from IRC (and no reconnect).
func (b *BotImpl) Quit(msg string) {
	b.reconnect = false
	b.conn.Quit(msg)
}

// Commands returns the underlying Commands.
func (b *BotImpl) Commands() *Commands { return b.commands }

// Shortcuts to b.Conn to ease mocking of Bot interface.
func (b *BotImpl) Action(t, msg string)       { b.Conn().Action(t, msg) }
func (b *BotImpl) Connected() bool            { return b.Conn().Connected() }
func (b *BotImpl) Me() *state.Nick            { return b.Conn().Me() }
func (b *BotImpl) Mode(t string, m ...string) { b.Conn().Mode(t, m...) }
func (b *BotImpl) Notice(t, msg string)       { b.Conn().Notice(t, msg) }
func (b *BotImpl) Privmsg(t, msg string)      { b.Conn().Privmsg(t, msg) }
func (b *BotImpl) Conn() *client.Conn         { return b.conn }
