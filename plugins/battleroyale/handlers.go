// Package battleroyale implements a plugin to see user ranking on Battle Royale game.
package battleroyale

import (
	"fmt"
	"sort"
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
	e.Bot.Privmsg(e.Target, fmt.Sprintf("%s: %s", name, p.String()))
}

func brAll(e *bot.Event, players map[string]string) {
	var l []*playerInfo
	for name, uid := range players {
		p, err := scoreByUID(uid)
		if err != nil {
			if err == errNotFound {
				continue // player has no stats yet
			}
			e.Bot.Privmsg(e.Target, err.Error())
			return
		}
		p.Name = name // override steam name with player name
		l = append(l, p)
	}
	if len(l) == 0 {
		e.Bot.Privmsg(e.Target, "no player has stats yet")
		return
	}
	sort.Sort(sort.Reverse(byPoints(l)))
	var o []string
	for _, p := range l {
		name := nohl.Nick(e.Bot, e.Target, p.Name)
		o = append(o, fmt.Sprintf("%s (%s)", name, p.Short()))
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
