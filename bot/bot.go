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

// NewBotWithProxyPassword creates a new Bot implementation with options.
func NewBotOptions(options ...func(*BotImpl)) (Bot, error) {
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
		SplitLen:    450,
		Proxy:       "",
		Pass:        "",
	}

	b := &BotImpl{
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

func Host(host string) func(*BotImpl) {
	return func(b *BotImpl) { b.config.Server = host }
}

func Nick(nick string) func(*BotImpl) {
	return func(b *BotImpl) { b.config.Me.Nick = nick }
}

func SSL(ssl bool) func(*BotImpl) {
	return func(b *BotImpl) { b.config.SSL = ssl }
}

func Ident(ident string) func(*BotImpl) {
	return func(b *BotImpl) { b.config.Me.Ident = ident }
}

func RealName(realName string) func(*BotImpl) {
	return func(b *BotImpl) { b.config.Me.Name = realName }
}

func Proxy(proxy string) func(*BotImpl) {
	return func(b *BotImpl) { b.config.Proxy = proxy }
}

func Password(password string) func(*BotImpl) {
	return func(b *BotImpl) { b.config.Pass = password }
}

func Channels(channels []string) func(*BotImpl) {
	return func(b *BotImpl) { b.channels = channels }
}

func CommandPrefix(commandPrefix string) func(impl *BotImpl) {
	return func(b *BotImpl) { b.commandPrefix = commandPrefix }
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

// BotImpl implements Bot.
type BotImpl struct {
	config        *client.Config
	conn          *client.Conn
	reconnect     bool
	quit          chan bool
	commands      *Commands
	channels      []string
	commandPrefix string
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
func (b *BotImpl) Action(t, msg string) { b.Conn().Action(t, msg) }
func (b *BotImpl) Connected() bool      { return b.Conn().Connected() }
func (b *BotImpl) HandleFunc(n string, h client.HandlerFunc) client.Remover {
	return b.Conn().HandleFunc(n, h)
}
func (b *BotImpl) Invite(nick, channel string) { b.Conn().Invite(nick, channel) }
func (b *BotImpl) Join(c string)               { b.Conn().Join(c) }
func (b *BotImpl) Me() *state.Nick             { return b.Conn().Me() }
func (b *BotImpl) Mode(t string, m ...string)  { b.Conn().Mode(t, m...) }
func (b *BotImpl) Nick(nick string)            { b.Conn().Nick(nick) }
func (b *BotImpl) Notice(t, msg string)        { b.Conn().Notice(t, msg) }
func (b *BotImpl) Part(c string, m ...string)  { b.Conn().Part(c, m...) }
func (b *BotImpl) Privmsg(t, msg string)       { b.Conn().Privmsg(t, msg) }
func (b *BotImpl) Conn() *client.Conn          { return b.conn }
func (b *BotImpl) Channels() []string {
	var channels []string
	for name := range b.Conn().Me().Channels {
		channels = append(channels, name)
	}
	sort.Strings(channels)
	return channels
}

func (b *BotImpl) setup() {
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
				b.channels = b.Channels()
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

func (b *BotImpl) CommandPrefix() string {
	return b.commandPrefix
}
