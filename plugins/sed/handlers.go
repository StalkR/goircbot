// Package sed implements a plugin to replace pattern in sentences.
// When someone says s/pattern/replace/ bot replaces that someone's last line.
package sed

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/fluffle/goirc/client"
)

func watchLine(b *bot.Bot, line *client.Line, bl *Backlog) {
	channel := line.Args[0]
	nick := line.Nick
	text := line.Args[1]
	if !strings.HasPrefix(channel, "#") {
		return
	}
	r, err := regexp.Compile("^s/([^/]+)/([^/]+)(?:/g?)?")
	if err != nil {
		b.Conn.Privmsg(channel, fmt.Sprintf("error: %s", err))
		return
	}
	m := r.FindSubmatch([]byte(text))
	if m == nil {
		bl.Store(channel, nick, text)
		return
	}
	meant := bl.Sed(channel, nick, string(m[1]), string(m[2]))
	if meant == "" {
		return
	}
	b.Conn.Privmsg(channel, fmt.Sprintf("%s meant: %s", nick, meant))
	bl.Store(channel, nick, meant)
}

// Register registers the plugin with a bot.
func Register(b *bot.Bot) {
	bl := &Backlog{}
	b.Conn.HandleFunc("privmsg",
		func(conn *client.Conn, line *client.Line) { watchLine(b, line, bl) })
}
