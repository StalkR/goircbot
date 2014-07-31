package battleroyale

import (
	"fmt"
	"sync"
)

type score struct {
	wins, losses, kills int
}

func (s score) String() string {
	return fmt.Sprintf("wins (%d), losses (%d), kills (%d)",
		s.wins, s.losses, s.kills)
}

type Scoreboard struct {
	players      map[string]string // name -> ID
	sync.RWMutex                   // protects below
	scores       map[string]score  // ID -> score
}

func NewScoreboard(players map[string]string) *Scoreboard {
	return &Scoreboard{players: players}
}

func (s *Scoreboard) Refresh() error {
	m, err := get()
	if err != nil {
		return err
	}
	s.Lock()
	s.scores = m
	s.Unlock()
	return nil
}

func (s *Scoreboard) Players() []string {
	var names []string
	for name := range s.players {
		names = append(names, name)
	}
	return names
}

func (s *Scoreboard) Get(name string) (score, error) {
	id, ok := s.players[name]
	if !ok {
		return score{}, fmt.Errorf("no such player")
	}
	s.RLock()
	defer s.RUnlock()
	r, ok := s.scores[id]
	if !ok {
		return score{}, fmt.Errorf("not found")
	}
	return r, nil
}

func (s *Scoreboard) Status() map[string]score {
	s.RLock()
	defer s.RUnlock()
	m := make(map[string]score)
	for name, id := range s.players {
		if r, ok := s.scores[id]; ok {
			m[name] = r
		}
	}
	return m
}
