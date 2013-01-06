// Package failotron implements a plugin in which users of a channel can ask the
// bot to randomly select a human (non-bot) on the channel for the next fail.
package failotron

import (
	"fmt"
	bot "goircbot"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func failotron(b *bot.Bot, e *bot.Event) {
	ch, on := b.Conn.Me.IsOnStr(e.Target)
	if !on {
		return
	}
	nicks := ch.Nicks()
	humans := make([]string, 0, len(nicks))
	for _, nick := range nicks {
		if !nick.Modes.Bot {
			humans = append(humans, nick.Name)
		}
	}
	if len(humans) == 0 {
		return
	}
	target := humans[rand.Intn(len(humans))]
	b.Conn.Privmsg(e.Target, fmt.Sprintf("FAIL-O-TRON ===> %s <=== FAIL-O-TRON", target))
}

// Register registers the plugin with a bot.
func Register(b *bot.Bot) {
	b.AddCommand("failotron", bot.Command{
		Help:    "find who is going to have the next fail",
		Handler: failotron,
		Pub:     true,
		Priv:    false,
		Hidden:  false})
}
