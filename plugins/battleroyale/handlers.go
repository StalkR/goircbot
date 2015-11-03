// Package battleroyale implements a plugin to view players on Battle Royale leaderboard.
package battleroyale

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/nohl"
)

func br(e *bot.Event, players map[string]uint64) {
	name := strings.TrimSpace(e.Args)
	if len(name) == 0 {
		brAll(e, players)
		return
	}
	steamID, ok := players[name]
	if !ok {
		parsed, err := strconv.ParseUint(name, 10, 64)
		if err != nil {
			e.Bot.Privmsg(e.Target, "unknown player name / invalid steam ID")
			return
		}
		steamID = parsed
	}
	p, err := viewPlayer(steamID)
	if err != nil {
		e.Bot.Privmsg(e.Target, err.Error())
		return
	}
	name = nohl.Nick(e.Bot, e.Target, name)
	e.Bot.Privmsg(e.Target, fmt.Sprintf("%s: %s", name, p.String()))
}

func brAll(e *bot.Event, players map[string]uint64) {
	var s []player
	for name, steamID := range players {
		p, err := viewPlayer(steamID)
		if err != nil {
			if err == errNotFound { // player has no stats yet
				s = append(s, player{Name: name})
				continue
			}
			e.Bot.Privmsg(e.Target, err.Error())
			return
		}
		p.Name = name // override returned player name with our own
		s = append(s, *p)
	}
	if len(s) == 0 {
		e.Bot.Privmsg(e.Target, "no player has stats yet")
		return
	}
	sort.Sort(sort.Reverse(byGlobalRank(s)))
	var o []string
	for _, p := range s {
		name := nohl.Nick(e.Bot, e.Target, p.Name)
		details := "n/a"
		if p.GlobalRank != 0 {
			details = p.Short()
		}
		o = append(o, fmt.Sprintf("%s (%s)", name, details))
	}
	e.Bot.Privmsg(e.Target, strings.Join(o, ", "))
}

// Register registers the plugin with a bot.
func Register(b bot.Bot, players map[string]uint64) {
	b.Commands().Add("br", bot.Command{
		Help:    "see player ranking on Battle Royale leaderboard",
		Handler: func(e *bot.Event) { br(e, players) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
