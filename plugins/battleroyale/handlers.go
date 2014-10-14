// Package battleroyale implements a plugin to see user ranking on Battle Royale game.
package battleroyale

import (
	"fmt"
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/nohl"
)

func br(e *bot.Event, players map[string]string) {
	name := strings.TrimSpace(e.Args)
	if len(name) == 0 {
		brAll(e, players)
		return
	}
	var p *playerInfo
	var err error
	if uid, ok := players[name]; ok {
		p, err = scoreByUID(uid)
	} else {
		p, err = scoreByName(name)
	}
	if err != nil {
		e.Bot.Privmsg(e.Target, err.Error())
		return
	}
	name = nohl.Nick(e.Bot, e.Target, name)
	e.Bot.Privmsg(e.Target, fmt.Sprintf("%s: %d wins, %d kills, %d loss, %d points, K/D %.2f, W/L %.2f, max kill distance %.2f",
		name, p.Wins, p.Kills, p.Loss, p.Points, p.KillDeathRatio, p.WinRate, p.MaxKillDistance))
}

func brAll(e *bot.Event, players map[string]string) {
	var o []string
	for name, uid := range players {
		p, err := scoreByUID(uid)
		if err != nil {
			e.Bot.Privmsg(e.Target, err.Error())
			return
		}
		name = nohl.Nick(e.Bot, e.Target, name)
		o = append(o, fmt.Sprintf("%s (%d W, %d K, %d L, %d pts)",
			name, p.Wins, p.Kills, p.Loss, p.Points))
	}
	e.Bot.Privmsg(e.Target, strings.Join(o, ", "))
}

// Register registers the plugin with a bot.
func Register(b bot.Bot, players map[string]string) {
	b.Commands().Add("br", bot.Command{
		Help:    "see player ranking on Battle Royale",
		Handler: func(e *bot.Event) { br(e, players) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
