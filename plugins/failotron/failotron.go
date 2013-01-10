// Package failotron implements a plugin in which users of a channel can ask the
// bot to randomly select a human (non-bot) on the channel for the next fail.
package failotron

import (
	"fmt"
	bot "github.com/StalkR/goircbot"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func failotron(b *bot.Bot, e *bot.Event, ignore []string) {
	ch, on := b.Conn.Me.IsOnStr(e.Target)
	if !on {
		return
	}
	ignoremap := make(map[string]bool)
	for _, nick := range ignore {
		ignoremap[nick] = true
	}
	nicks := ch.Nicks()
	humans := make([]string, 0, len(nicks))
	for _, nick := range nicks {
		if nick.Modes.Bot {
			continue
		}
		if _, present := ignoremap[nick.Nick]; present {
			continue
		}
		humans = append(humans, nick.Nick)
	}
	if len(humans) == 0 {
		return
	}
	target := humans[rand.Intn(len(humans))]
	b.Conn.Privmsg(e.Target, fmt.Sprintf("FAIL-O-TRON ===> %s <=== FAIL-O-TRON", target))
}

// Register registers the plugin with a bot.
// Use ignore as a list of nicks to ignore.
func Register(b *bot.Bot, ignore []string) {
	b.AddCommand("failotron", bot.Command{
		Help:    "find who is going to have the next fail",
		Handler: func(b *bot.Bot, e *bot.Event) { failotron(b, e, ignore) },
		Pub:     true,
		Priv:    false,
		Hidden:  false})
}
