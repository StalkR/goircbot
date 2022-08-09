package scores

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// Scores is the main structure to hold scores in the plugin.
type Scores struct {
	sync.Mutex
	scores map[string]int
	dirty  bool
}

// NewScores returns a new initialized Scores.
func NewScores() *Scores {
	s := &Scores{}
	s.scores = make(map[string]int)
	return s
}

// Add adds n score points to a given thing.
func (s *Scores) Add(thing string, n int) {
	if n == 0 {
		return
	}
	s.Lock()
	defer s.Unlock()
	score := s.scores[thing] // If not present, default value 0.
	if score+n == 0 {
		delete(s.scores, thing)
	} else {
		s.scores[thing] = score + n
	}
	s.dirty = true
}

// Score returns the score of a given thing.
func (s *Scores) Score(thing string) int {
	s.Lock()
	defer s.Unlock()
	score := s.scores[thing] // If not present, default value 0.
	return score
}

// List sorts scores and returns an ordered ScoreList.
func (s *Scores) List() *ScoreList {
	s.Lock()
	defer s.Unlock()
	o := make(ScoreList, 0, len(s.scores))
	for name, value := range s.scores {
		o = append(o, &ScoreEntry{name, value})
	}
	sort.Sort(o)
	return &o
}

// String returns formatted top +/- scores and total.
func (s *Scores) String() string {
	s.Lock()
	if len(s.scores) == 0 {
		s.Unlock()
		return "no scores yet"
	}
	s.Unlock()
	l := *s.List()
	min := 3
	if len(l) < min {
		min = len(l)
	}
	plus := make([]string, 0, min)
	for i := 0; i < min; i++ {
		plus = append(plus, l[len(l)-1-i].String())
	}
	minus := make([]string, 0, min)
	for i := 0; i < min; i++ {
		minus = append(minus, l[i].String())
	}
	return fmt.Sprintf("High: %s; Low: %s; Total things scored: %d",
		strings.Join(plus, ", "), strings.Join(minus, ", "), len(l))
}

// ScoreList is a slice of score entries and implements sort.Interface.
// Not directly used by the plugin but used to calculate top scores.
type ScoreList []*ScoreEntry

func (l ScoreList) Len() int { return len(l) }

func (l ScoreList) Swap(i, j int) { l[i], l[j] = l[j], l[i] }

func (l ScoreList) Less(i, j int) bool { return l[i].Value < l[j].Value }

// ScoreEntry represents a single entry of a ScoreList.
type ScoreEntry struct {
	Name  string
	Value int
}

func (e *ScoreEntry) String() string {
	return fmt.Sprintf("%s (%d)", e.Name, e.Value)
}
