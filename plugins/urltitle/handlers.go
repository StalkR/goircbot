// Package urltitle implements a plugin to watch web URLs, fetch and display title.
package urltitle

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/fluffle/goirc/client"
)

func watchLine(b *bot.Bot, line *client.Line, ignoremap map[string]bool) {
	target := line.Args[0]
	if !strings.HasPrefix(target, "#") {
		return
	}
	if _, ignore := ignoremap[line.Nick]; ignore {
		return
	}
	text := line.Args[1]
	if m, err := regexp.Match(silenceRegexp, []byte(text)); err != nil || m {
		return
	}
	r, err := regexp.Compile("(?:^|\\s)(https?://[^\\s]+)")
	if err != nil {
		return
	}
	matches := r.FindSubmatch([]byte(text))
	if len(matches) < 2 {
		return
	}
	url := string(matches[1])
	title, err := Title(url)
	if err != nil {
		log.Println("urltitle:", err)
		return
	}
	if len(title) > 200 {
		title = title[:200]
	}
	b.Conn.Privmsg(target, fmt.Sprintf("%s :: %s", url, title))
}

// Register registers the plugin with a bot.
func Register(b *bot.Bot, ignore []string) {
	ignoremap := make(map[string]bool)
	for _, nick := range ignore {
		ignoremap[nick] = true
	}

	b.Conn.HandleFunc("privmsg",
		func(conn *client.Conn, line *client.Line) {
			watchLine(b, line, ignoremap)
		})
}
