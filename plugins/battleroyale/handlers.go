// Package battleroyale implements a plugin to see user ranking on Battle Royale game.
package battleroyale

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/nohl"
)

func br(e *bot.Event, s *Scoreboard) {
	name := strings.TrimSpace(e.Args)
	if len(name) == 0 {
		var n []string
		for _, name := range s.Players() {
			n = append(n, nohl.Nick(e.Bot, e.Target, name))
		}
		e.Bot.Privmsg(e.Target, fmt.Sprintf("Players: %s", strings.Join(n, ", ")))
		return
	}
	r, err := s.Get(name)
	if err != nil {
		e.Bot.Privmsg(e.Target, err.Error())
		return
	}
	e.Bot.Privmsg(e.Target, fmt.Sprintf("%s: %s",
		nohl.Nick(e.Bot, e.Target, name), r.String()))
}

// Register registers the plugin with a bot.
func Register(b bot.Bot, players map[string]string) {
	s := NewScoreboard(players)
	go refresh(b, s)
	b.Commands().Add("br", bot.Command{
		Help:    "see player ranking on Battle Royale",
		Handler: func(e *bot.Event) { br(e, s) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}

func refresh(b bot.Bot, s *Scoreboard) {
	var prev map[string]score
	for ; ; <-time.Tick(time.Hour) {
		if err := s.Refresh(); err != nil {
			log.Printf("battleroyale: refresh error: %v", err)
		}
		log.Print("battleroyale: scores refreshed")
		if prev == nil {
			prev = s.Status()
			continue
		}
		current := s.Status()
		for name, newscore := range current {
			oldscore, ok := prev[name]
			if !ok || !reflect.DeepEqual(oldscore, newscore) {
				notify(b, fmt.Sprintf("New score for %s: %s", name, newscore.String()))
			}
		}
		prev = current
	}
}

func notify(b bot.Bot, line string) {
	if !b.Connected() {
		return
	}
	for _, channel := range b.Me().Channels() {
		b.Privmsg(channel.Name, line)
	}
}
