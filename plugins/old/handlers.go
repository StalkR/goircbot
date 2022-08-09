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
	linkRE    = regexp.MustCompile(`(?:^|\s)(https?://[^\s]+)`)
	backlogRE = regexp.MustCompile("<[+%@&~ ]?[a-zA-Z0-9_`^\\[\\]-]+>")
	silenceRE = regexp.MustCompile(`(^|\s)tg(\)|\s|$)`) // Line ignored if matched.
)

func readURLs(b bot.Bot, line *client.Line, o *Old, ignore map[string]bool) {
	target := line.Args[0]
	if !strings.HasPrefix(target, "#") {
		return
	}
	if _, ignore := ignore[line.Nick]; ignore {
		return
	}
	text := line.Args[1]
	if backlogRE.MatchString(text) || silenceRE.MatchString(text) {
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
		b.Privmsg(target, fmt.Sprintf("old! first shared by %v %v ago", nick, ago))
	}
}

// Register registers the plugin with a bot.
func Register(b bot.Bot, oldfile string, ignore []string) {
	ignoremap := make(map[string]bool)
	for _, nick := range ignore {
		ignoremap[nick] = true
	}

	o := load(oldfile)

	b.Conn().HandleFunc("privmsg",
		func(conn *client.Conn, line *client.Line) { readURLs(b, line, o, ignoremap) })

	// Every minute, save to file.
	if len(oldfile) > 0 {
		go func() {
			for range time.Tick(time.Minute) {
				save(oldfile, o)
			}
		}()
	}

	// Every day, clean URLs older than a year so it does not grow infinitely.
	go func() {
		for range time.Tick(time.Hour * 24) {
			o.Clean(time.Hour * 24 * 365)
		}
	}()
}
