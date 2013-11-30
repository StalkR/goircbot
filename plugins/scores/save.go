package scores

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func load(scoresfile string) *Scores {
	s := NewScores()
	if len(scoresfile) == 0 {
		return s
	}
	b, err := ioutil.ReadFile(scoresfile)
	if err != nil {
		log.Println("scores: unable to open scores file")
		return s
	}
	if err := json.Unmarshal(b, &s.scores); err != nil {
		log.Println("scores: unable to load scores")
		return s
	}
	log.Println("scores: loaded successfully")
	return s
}

func save(scoresfile string, s *Scores) {
	s.Lock()
	defer s.Unlock()
	if !s.dirty {
		return
	}
	b, err := json.Marshal(s.scores)
	if err != nil {
		log.Println("scores: unable to encode scores for saving")
		return
	}
	if err := ioutil.WriteFile(scoresfile, b, 0644); err != nil {
		log.Println("scores: unable to save scores")
		return
	}
	s.dirty = false
}
