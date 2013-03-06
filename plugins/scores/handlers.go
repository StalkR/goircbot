// Package scores implements a plugin to score things on channels.
// One can do X++ (or X--) to give (or take) points to X.
package scores

import (
	"fmt"
	bot "github.com/StalkR/goircbot"
	"github.com/fluffle/goirc/client"
	"log"
	"regexp"
	"strings"
	"time"
)

func ParseScore(b *bot.Bot, line *client.Line, s *Scores) {
	text := line.Args[1]
	var modifier int
	switch {
	case !strings.HasPrefix(line.Args[0], "#") || len(text) < 3:
		return
	case strings.HasSuffix(text, "++"):
		modifier = 1
	case strings.HasSuffix(text, "--"):
		// We allow - (not +) in thing but not at the end to avoid x---.
		if text[len(text)-3] == '-' {
			return
		}
		modifier = -1
	default:
		return
	}
	target := line.Args[0]
	thing := Sanitize(text[:len(text)-2])
	match, err := regexp.Match("^[-_a-zA-Z0-9/ `'\":;\\\\]+$", []byte(thing))
	if err != nil {
		log.Println("scores: regexp error", err)
		return
	}
	if !match {
		return
	}
	if thing == line.Nick && modifier == 1 {
		modifier = -1
		reply := fmt.Sprintf("Scoring for yourself? %s--", thing)
		b.Conn.Privmsg(target, reply)
	}
	s.Lock()
	defer s.Unlock()
	score, present := s.Map[thing]
	if !present {
		score = 0
	}
	newScore := score + modifier
	s.Map[thing] = newScore
	if newScore == 0 {
		delete(s.Map, thing)
	}
	reply := fmt.Sprintf("%s is now %d", thing, newScore)
	b.Conn.Privmsg(target, reply)
}

func Sanitize(text string) string {
	clean := RemoveChars(text, " ", "` ", "\\", "\"", "'", ":", ";")
	if len(clean) > 128 {
		return clean[:128]
	}
	return clean
}

func RemoveChars(s string, chars ...string) string {
	for _, c := range chars {
		s = strings.Replace(s, c, "", -1)
	}
	return s
}

func ShowScore(b *bot.Bot, e *bot.Event, s *Scores) {
	thing := strings.TrimSpace(e.Args)
	if len(thing) == 0 {
		return
	}
	s.Lock()
	defer s.Unlock()
	b.Conn.Privmsg(e.Target, s.ScoreOf(thing))
}

func TopScores(b *bot.Bot, e *bot.Event, s *Scores) {
	s.Lock()
	defer s.Unlock()
	if len(s.Map) == 0 {
		b.Conn.Privmsg(e.Target, "no scores yet")
		return
	}
	b.Conn.Privmsg(e.Target, s.String())
}

// Register registers the plugin with a bot.
func Register(b *bot.Bot, scoresfile string) {
	s := Load(scoresfile)

	b.Conn.HandleFunc("privmsg",
		func(conn *client.Conn, line *client.Line) { ParseScore(b, line, s) })

	b.AddCommand("score", bot.Command{
		Help:    "score <thing> - show score of something",
		Handler: func(b *bot.Bot, e *bot.Event) { ShowScore(b, e, s) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})

	b.AddCommand("scores", bot.Command{
		Help:    "show top +/- scores",
		Handler: func(b *bot.Bot, e *bot.Event) { TopScores(b, e, s) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})

	if len(scoresfile) > 0 {
		b.AddCron("name", bot.Cron{
			Handler:  func(b *bot.Bot) { Save(scoresfile, s) },
			Duration: time.Minute})
	}
}
