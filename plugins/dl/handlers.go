// Package dl implements a plugin to see downloads stats from mldonkey/transmission.
package dl

import (
	"log"
	"regexp"
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/mldonkey"
	"github.com/StalkR/goircbot/lib/transmission"
)

var (
	transmissionRE = regexp.MustCompile(`^(https?://|magnet:?)`)
	mldonkeyRE     = regexp.MustCompile(`^ed2k://`)
)

func handleMLDonkey(b *bot.Bot, e *bot.Event, url string) {
	c, err := mldonkey.New(url)
	if err != nil {
		log.Println("dl: mldonkey new:", err)
		return
	}
	args := strings.SplitN(strings.TrimSpace(e.Args), " ", 2)
	link := args[0]
	if link == "" {
		stats, err := c.Stats()
		if err != nil {
			log.Println("dl: mldonkey:", err)
			return
		}
		b.Conn.Privmsg(e.Target, "[MLDonkey] "+stats.String())
		return
	}
	if !mldonkeyRE.MatchString(link) {
		return
	}
	if err := c.Add(link); err != nil {
		log.Println("dl: mldonkey add:", err)
		return
	}
	b.Conn.Privmsg(e.Target, "[MLDonkey] Added")
}

func handleTransmission(b *bot.Bot, e *bot.Event, url string) {
	c, err := transmission.New(url)
	if err != nil {
		log.Println("dl: transmission new:", err)
		return
	}
	args := strings.SplitN(strings.TrimSpace(e.Args), " ", 2)
	link := args[0]
	if link == "" {
		stats, err := c.Stats()
		if err != nil {
			log.Println("dl: transmission stats:", err)
			return
		}
		b.Conn.Privmsg(e.Target, "[Transmission] "+stats.String())
		return
	}
	if !transmissionRE.MatchString(link) {
		return
	}
	name, err := c.Add(link)
	if err != nil {
		log.Println("dl: transmission add:", err)
		return
	}
	b.Conn.Privmsg(e.Target, "[Transmission] Added: "+name)
}

// Register registers the plugin with a bot.
func Register(b *bot.Bot, mldonkeyURL, transmissionURL string) {

	b.AddCommand("dl", bot.Command{
		Help: "See downloads status or add downloads",
		Handler: func(b *bot.Bot, e *bot.Event) {
			if mldonkeyURL != "" {
				handleMLDonkey(b, e, mldonkeyURL)
			}
			if transmissionURL != "" {
				handleTransmission(b, e, transmissionURL)
			}
		},
		Pub:    true,
		Priv:   false,
		Hidden: false})
}
