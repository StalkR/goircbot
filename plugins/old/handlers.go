// Package old implements a plugin to remember URLs and announce duplicates.
package old

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/duration"
	"github.com/StalkR/goircbot/lib/nohl"
	"github.com/fluffle/goirc/client"
)

var (
	linkRE    = regexp.MustCompile(`(?:^|\s)(https?://[^#\s]+)`)
	backlogRE = regexp.MustCompile("<[+%@&~]?[a-zA-Z0-9_`^\\[\\]-]+>")
)

func readURLs(b *bot.Bot, line *client.Line, o *Old, ignore map[string]bool) {
	target := line.Args[0]
	if !strings.HasPrefix(target, "#") {
		return
	}
	if _, ignore := ignore[line.Nick]; ignore {
		return
	}
	text := line.Args[1]
	if backlogRE.MatchString(text) {
		return
	}

	matches := linkRE.FindAllStringSubmatch(text, -1)
	if matches == nil {
		return
	}
	for _, submatches := range matches {
		url := submatches[1]
		i, err := o.Old(url)
		if err != nil {
			if err = o.Add(url, target, line.Nick); err != nil {
				log.Print(err)
			}
			return
		}
		ago := duration.Format(time.Since(i.Time))
		nick := nohl.Nick(b, target, i.Nick)
		b.Conn.Privmsg(target, fmt.Sprintf("old! first shared by %v %v ago", nick, ago))
	}
}

// Register registers the plugin with a bot.
func Register(b *bot.Bot, oldfile string, ignore []string) {
	ignoremap := make(map[string]bool)
	for _, nick := range ignore {
		ignoremap[nick] = true
	}

	o := load(oldfile)

	b.Conn.HandleFunc("privmsg",
		func(conn *client.Conn, line *client.Line) { readURLs(b, line, o, ignoremap) })

	if len(oldfile) > 0 {
		b.AddCron("old-save", bot.Cron{
			Handler:  func(b *bot.Bot) { save(oldfile, o) },
			Duration: time.Minute})
	}

	// Every day, clean URLs older than a year so it does not grow infinitely.
	b.AddCron("old-clean", bot.Cron{
		Handler:  func(b *bot.Bot) { o.Clean(time.Hour * 24 * 365) },
		Duration: time.Hour * 24})
}
