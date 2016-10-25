// Package battlefield implements a plugin to view Battlefield 1 player stats.
package battlefield

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/nohl"
)

func handle(e *bot.Event, s *Session, players map[string]uint64, game string) {
	name := strings.ToLower(strings.TrimSpace(e.Args))
	if len(name) == 0 {
		showAll(e, s, players, game)
		return
	}
	personaID, err := func() (uint64, error) {
		// case-insensitive find
		for p, id := range players {
			if strings.ToLower(p) == name {
				return id, nil
			}
		}
		// not found, try to parse id
		parsed, err := strconv.ParseUint(name, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("unknown player name / invalid persona ID")
		}
		return parsed, nil
	}()
	if err != nil {
		e.Bot.Privmsg(e.Target, err.Error())
	}
	// if possible take name from map instead of arg
	for p, id := range players {
		if personaID == id {
			name = p
			break
		}
	}
	stats, err := s.Stats(personaID, name)
	if err != nil {
		e.Bot.Privmsg(e.Target, err.Error())
		return
	}
	name = nohl.Nick(e.Bot, e.Target, name)
	switch game {
	case "bf1":
		e.Bot.Privmsg(e.Target, fmt.Sprintf("%s: %s", name, stats.BF1.String()))
	case "bf4":
		e.Bot.Privmsg(e.Target, fmt.Sprintf("%s: %s", name, stats.BF4.String()))
	}
}

func showAll(e *bot.Event, s *Session, players map[string]uint64, game string) {
	var allStats []Stats
	for name, personaID := range players {
		stats, err := s.Stats(personaID, name)
		if err != nil {
			if err == errNotFound { // player has no stats yet
				continue
			}
			e.Bot.Privmsg(e.Target, err.Error())
			return
		}
		allStats = append(allStats, *stats)
	}
	if len(allStats) == 0 {
		e.Bot.Privmsg(e.Target, "no player has stats yet")
		return
	}
	switch game {
	case "bf1":
		sort.Sort(sort.Reverse(byRankBF1(allStats)))
	case "bf4":
		sort.Sort(sort.Reverse(byRankBF4(allStats)))
	}
	var o []string
	for _, stats := range allStats {
		name := nohl.Nick(e.Bot, e.Target, stats.Name)
		switch game {
		case "bf1":
			o = append(o, fmt.Sprintf("%s (%s)", name, stats.BF1.Short()))
		case "bf4":
			o = append(o, fmt.Sprintf("%s (%s)", name, stats.BF4.Short()))
		}
	}
	e.Bot.Privmsg(e.Target, strings.Join(o, ", "))
}

// Register registers the plugin with a bot.
func Register(b bot.Bot, email, password string, players map[string]uint64) error {
	s := NewSession(email, password)
	if err := s.Login(); err != nil {
		return err
	}

	b.Commands().Add("bf1", bot.Command{
		Help:    "see Battlefield 1 player stats",
		Handler: func(e *bot.Event) { handle(e, s, players, "bf1") },
		Pub:     true,
		Priv:    true,
		Hidden:  false})

	b.Commands().Add("bf4", bot.Command{
		Help:    "see Battlefield 4 player stats",
		Handler: func(e *bot.Event) { handle(e, s, players, "bf4") },
		Pub:     true,
		Priv:    true,
		Hidden:  false})

	return nil
}
