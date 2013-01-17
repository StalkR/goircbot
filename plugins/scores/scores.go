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
	Map map[string]int
}

// NewScores returns a new initialized Scores.
func NewScores() *Scores {
	var s Scores
	s.Map = make(map[string]int)
	return &s
}

// ScoreOf returns a formatted string with the score of a given thing.
func (s *Scores) ScoreOf(thing string) string {
	score, present := s.Map[thing]
	if !present {
		score = 0
	}
	return fmt.Sprintf("%s is %d", thing, score)
}

// List sorts scores and returns an ordered ScoreList.
// It assumes the lock has already been taken.
func (s *Scores) List() *ScoreList {
	o := make(ScoreList, 0, len(s.Map))
	for name, value := range s.Map {
		o = append(o, &ScoreEntry{name, value})
	}
	sort.Sort(o)
	return &o
}

// String returns formatted top +/- scores and total.
// It assumes the lock has already been taken.
func (s *Scores) String() string {
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

// Score represents a single entry of a ScoreList.
type ScoreEntry struct {
	Name  string
	Value int
}

func (e *ScoreEntry) String() string {
	return fmt.Sprintf("%s (%d)", e.Name, e.Value)
}
