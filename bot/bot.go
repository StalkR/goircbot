// Package bot implements an IRC bot with plugins.
package bot

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/StalkR/goircbot/lib/tls"
	"github.com/fluffle/goirc/client"
	"github.com/fluffle/goirc/state"
)

// Bot represents an IRC bot, with IRC client object, settings, commands.
type Bot interface {
	Run()                                                     // Run bot, reconnect if disconnect.
	Quit(msg string)                                          // Quit bot from IRC with a msg.
	Commands() *Commands                                      // For plugins to Add/Del commands.
	Action(t, msg string)                                     // Shortcut to Conn().Action()
	Connected() bool                                          // Shortcut to Conn().Connected()
	HandleFunc(n string, h client.HandlerFunc) client.Remover // Shortcut to Conn().HandleFunc()
	Invite(nick, channel string)                              // Shortcut to Conn().Nick()
	Join(c string)                                            // Shortcut to Conn().Join()
	Me() *state.Nick                                          // Shortcut to Conn().Me()
	Mode(t string, m ...string)                               // Shortcut to Conn().Mode()
	Nick(nick string)                                         // Shortcut to Conn().Nick()
	Notice(t, msg string)                                     // Shortcut to Conn().Notice()
	Part(c string, m ...string)                               // Shortcut to Conn().Part()
	Privmsg(t, msg string)                                    // Shortcut to Conn().Privmsg()
	Conn() *client.Conn                                       // Conn returns the underlying goirc client connection.
	Channels() []string                                       // Channels returns list of channels which bot has joined.
	CommandPrefix() string                                    // Return's the bot's prefix used to specify commands
}

// NewBotOptions creates a new Bot implementation with options.
func NewBotOptions(options ...Option) (Bot, error) {
	cfg := &client.Config{
		Me:          &state.Nick{Nick: "goircbot"},
		NewNick:     func(s string) string { return s + "_" },
		PingFreq:    3 * time.Minute,
		QuitMessage: "I have to go.",
		Server:      "",
		SSL:         false,
		SSLConfig:   nil,
		Version:     "Powered by GoIRCBot",
		Recover:     (*client.Conn).LogPanic,
		SplitLen:    450, // default is the same as goirc's unexported defaultLimit
		Proxy:       "",
		Pass:        "",
	}

	b := &botImpl{
		config:        cfg,
		reconnect:     true,
		quit:          make(chan bool),
		commands:      NewCommands(),
		commandPrefix: "!",
	}

	for _, option := range options {
		option(b)
	}

	if b.config.Server == "" {
		return nil, fmt.Errorf("Host option not given")
	}

	if b.config.SSL {
		hostPort := strings.SplitN(b.config.Server, ":", 2)
		b.config.SSLConfig = tls.Config(hostPort[0])
	}

	if b.config.Me.Ident == "" {
		b.config.Me.Ident = b.config.Me.Nick
	}

	if b.config.Me.Name == "" {
		b.config.Me.Name = b.config.Me.Nick
	}

	b.setup()
	return b, nil
}

// Option represents a configurable option.
type Option func(*botImpl)

// Host is an option to configure the IRC server host.
func Host(host string) func(*botImpl) {
	return func(b *botImpl) { b.config.Server = host }
}

// Nick is an option to configure the bot nick.
func Nick(nick string) func(*botImpl) {
	return func(b *botImpl) { b.config.Me.Nick = nick }
}

// SSL is an option to configure SSL (TLS) on the connection.
func SSL(ssl bool) func(*botImpl) {
	return func(b *botImpl) { b.config.SSL = ssl }
}

// Ident is an option to configure the bot ident.
func Ident(ident string) func(*botImpl) {
	return func(b *botImpl) { b.config.Me.Ident = ident }
}

// RealName is an option to configure the bot real name.
func RealName(realName string) func(*botImpl) {
	return func(b *botImpl) { b.config.Me.Name = realName }
}

// Proxy is an option to configure a proxy to connect to the IRC server.
func Proxy(proxy string) func(*botImpl) {
	return func(b *botImpl) { b.config.Proxy = proxy }
}

// Password is an option to configure the IRC server password.
func Password(password string) func(*botImpl) {
	return func(b *botImpl) { b.config.Pass = password }
}

// Channels is an option to configure the channels for the bot to join.
func Channels(channels []string) func(*botImpl) {
	return func(b *botImpl) { b.channels = channels }
}

// CommandPrefix is an option to configure the prefix of commands (default is !).
func CommandPrefix(commandPrefix string) func(impl *botImpl) {
	return func(b *botImpl) { b.commandPrefix = commandPrefix }
}

// SplitLen is an option to configure the split length for some IRC commands (default is 450).
func SplitLen(n int) func(impl *botImpl) {
	return func(b *botImpl) { b.config.SplitLen = n }
}

// NewBot creates a new Bot implementation with a set of parameters.
//
// Deprecated: use NewBotOptions instead.
func NewBot(host string, ssl bool, nick, ident string, channels []string) Bot {
	b, _ := NewBotOptions(Host(host), Nick(nick), SSL(ssl), Ident(ident), Channels(channels))
	return b
}

// NewBotWithProxy creates a new Bot implementation with a set of parameters including a proxy.
//
// Deprecated: use NewBotOptions instead.
func NewBotWithProxy(host string, ssl bool, nick, ident string, channels []string, proxy string) Bot {
	b, _ := NewBotOptions(Host(host), Nick(nick), SSL(ssl), Ident(ident), Proxy(proxy), Channels(channels))
	return b
}

// botImpl implements Bot.
type botImpl struct {
	config        *client.Config
	conn          *client.Conn
	reconnect     bool
	quit          chan bool
	commands      *Commands
	channels      []string
	commandPrefix string
}

// Run starts the Bot by connecting it to IRC. It automatically reconnects.
func (b *botImpl) Run() {
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
func (b *botImpl) Quit(msg string) {
	b.reconnect = false
	b.conn.Quit(msg)
}

// Commands returns the underlying Commands.
func (b *botImpl) Commands() *Commands { return b.commands }

// Shortcuts to b.Conn to ease mocking of Bot interface.

// Action is a shortcut to goirc connection Action().
func (b *botImpl) Action(t, msg string) { b.Conn().Action(t, msg) }

// Connected is a shortcut to goirc connection Connected().
func (b *botImpl) Connected() bool { return b.Conn().Connected() }

// HandleFunc is a shortcut to goirc connection HandleFunc().
func (b *botImpl) HandleFunc(n string, h client.HandlerFunc) client.Remover {
	return b.Conn().HandleFunc(n, h)
}

// Invite is a shortcut to goirc connection Invite().
func (b *botImpl) Invite(nick, channel string) { b.Conn().Invite(nick, channel) }

// Join is a shortcut to goirc connection Join().
func (b *botImpl) Join(c string) { b.Conn().Join(c) }

// Me is a shortcut to goirc connection Me().
func (b *botImpl) Me() *state.Nick { return b.Conn().Me() }

// Mode is a shortcut to goirc connection Mode().
func (b *botImpl) Mode(t string, m ...string) { b.Conn().Mode(t, m...) }

// Nick is a shortcut to goirc connection Nick().
func (b *botImpl) Nick(nick string) { b.Conn().Nick(nick) }

// Notice is a shortcut to goirc connection Notice().
func (b *botImpl) Notice(t, msg string) { b.Conn().Notice(t, msg) }

// Part is a shortcut to goirc connection Part().
func (b *botImpl) Part(c string, m ...string) { b.Conn().Part(c, m...) }

// Privmsg is a shortcut to goirc connection Privmsg().
func (b *botImpl) Privmsg(t, msg string) { b.Conn().Privmsg(t, msg) }

// Conn returns the goirc connection.
func (b *botImpl) Conn() *client.Conn { return b.conn }

// Channels returns a sorted list of channels the bot is currently on.
func (b *botImpl) Channels() []string {
	chans := b.Conn().Me().Channels
	channels := make([]string, 0, len(chans))
	for name := range chans {
		channels = append(channels, name)
	}
	sort.Strings(channels)
	return channels
}

func (b *botImpl) setup() {
	b.conn = client.Client(b.config)

	b.conn.EnableStateTracking()
	// On connect, mark ourselves as bot first then join channels (can be long).
	b.conn.HandleFunc("connected",
		func(conn *client.Conn, line *client.Line) {
			conn.Mode(conn.Me().Nick, "+B")
			for _, channel := range b.channels {
				conn.Join(channel)
			}
		})

	// Signal disconnect to Bot.Run so it can reconnect.
	b.conn.HandleFunc("disconnected",
		func(conn *client.Conn, line *client.Line) {
			channels := b.Channels()
			// On the first disconnect, we will remember channels. Then state is
			// reinitialized and another connection is attempted. If that one fails,
			// we end up with an empty state, so no channels and we do not want to
			// save that, it does not make much sense. So ignore empty values.
			if len(channels) > 0 {
				b.channels = channels
			}
			b.quit <- true
		})

	// Print out any error from the server
	b.conn.HandleFunc("error",
		func(conn *client.Conn, line *client.Line) {
			log.Println(conn.Config().Server, line.Cmd, line.Text())
		})

	b.conn.HandleFunc("privmsg",
		func(conn *client.Conn, line *client.Line) { b.commands.Handle(b, line) })

	b.commands.Add("help", Command{
		Help:    "show commands or detailed help",
		Handler: func(e *Event) { b.commands.Help(e) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}

func (b *botImpl) CommandPrefix() string {
	return b.commandPrefix
}
