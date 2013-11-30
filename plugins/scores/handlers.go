// Package scores implements a plugin to score things on channels.
// One can do X++ (or X--) to give (or take) points to X.
package scores

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/StalkR/goircbot/bot"
	"github.com/fluffle/goirc/client"
)

func parseScore(b *bot.Bot, line *client.Line, s *Scores) {
	text := strings.TrimSpace(line.Args[1])
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
	thing := sanitize(text[:len(text)-2])
	match, err := regexp.Match(`^[-_a-zA-Z0-9/ '":;\\`+"`]+$", []byte(thing))
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
	s.Add(thing, modifier)
	b.Conn.Privmsg(target, fmt.Sprintf("%s is now %d", thing, s.Score(thing)))
}

func sanitize(text string) string {
	clean := removeChars(text, " ", "` ", `\`, `"`, "'", ":", ";")
	if len(clean) > 128 {
		return clean[:128]
	}
	return clean
}

func removeChars(s string, chars ...string) string {
	for _, c := range chars {
		s = strings.Replace(s, c, "", -1)
	}
	return s
}

func showScore(b *bot.Bot, e *bot.Event, s *Scores) {
	thing := strings.TrimSpace(e.Args)
	if len(thing) == 0 {
		return
	}
	b.Conn.Privmsg(e.Target, fmt.Sprintf("%s is %d", thing, s.Score(thing)))
}

func topScores(b *bot.Bot, e *bot.Event, s *Scores) {
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
	s := load(scoresfile)

	b.Conn.HandleFunc("privmsg",
		func(conn *client.Conn, line *client.Line) { parseScore(b, line, s) })

	b.AddCommand("score", bot.Command{
		Help:    "score <thing> - show score of something",
		Handler: func(b *bot.Bot, e *bot.Event) { showScore(b, e, s) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})

	b.AddCommand("scores", bot.Command{
		Help:    "show top +/- scores",
		Handler: func(b *bot.Bot, e *bot.Event) { topScores(b, e, s) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})

	// Every minute, save to file.
	if len(scoresfile) > 0 {
		go func() {
			for _ = range time.Tick(time.Minute) {
				save(scoresfile, s)
			}
		}()
	}
}
