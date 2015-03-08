// Package failotron implements a plugin in which users of a channel can ask the
// bot to randomly select a human (non-bot) on the channel for the next fail.
package failotron

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/StalkR/goircbot/bot"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func failotron(e *bot.Event, ignore []string) {
	st := e.Bot.Conn().StateTracker()
	ch := st.GetChannel(e.Target)
	if ch == nil {
		return
	}
	ignoremap := make(map[string]bool)
	for _, nick := range ignore {
		ignoremap[nick] = true
	}
	var humans []string
	for nick := range ch.Nicks {
		if n := st.GetNick(nick); n == nil || n.Modes.Bot {
			continue
		}
		if _, present := ignoremap[nick]; present {
			continue
		}
		humans = append(humans, nick)
	}
	if len(humans) == 0 {
		return
	}
	target := humans[rand.Intn(len(humans))]
	e.Bot.Privmsg(e.Target, fmt.Sprintf("FAIL-O-TRON ===> %s <=== FAIL-O-TRON", target))
}

// Register registers the plugin with a bot.
// Use ignore as a list of nicks to ignore.
func Register(b bot.Bot, ignore []string) {
	b.Commands().Add("failotron", bot.Command{
		Help:    "find who is going to have the next fail",
		Handler: func(e *bot.Event) { failotron(e, ignore) },
		Pub:     true,
		Priv:    false,
		Hidden:  false})
}
