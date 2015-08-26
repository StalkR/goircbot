// Package hots implements a plugin to show Heroes of the Storm (HotS) player stats.
package hots

import (
	"fmt"
	"sort"
	"strings"

	"github.com/StalkR/goircbot/bot"
)

// Register registers the plugin with a bot.
func Register(b bot.Bot, players map[string]int) {
	b.Commands().Add("hots", bot.Command{
		Help:    "show Heroes of the Storm (HotS) stats of a player",
		Handler: func(e *bot.Event) { handleHots(e, players) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}

func handleHots(e *bot.Event, players map[string]int) {
	player := strings.ToLower(strings.TrimSpace(e.Args))
	if len(player) == 0 {
		return
	}
	// case-insensitive find
	id, ok := func() (int, bool) {
		for p, id := range players {
			if strings.ToLower(p) == player {
				return id, true
			}
		}
		return 0, false
	}()
	if !ok {
		var s []string
		for player := range players {
			s = append(s, player)
		}
		sort.Strings(s)
		e.Bot.Privmsg(e.Target, fmt.Sprintf("not found - known players: %s", strings.Join(s, ", ")))
		return
	}
	stats, err := NewStats(id)
	if err != nil {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	e.Bot.Privmsg(e.Target, stats.String())
}
