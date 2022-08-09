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

func parseScore(b bot.Bot, line *client.Line, s *Scores) {
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
		b.Privmsg(target, reply)
	}
	s.Add(thing, modifier)
	b.Privmsg(target, fmt.Sprintf("%s is now %d", thing, s.Score(thing)))
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

func showScore(e *bot.Event, s *Scores) {
	thing := strings.TrimSpace(e.Args)
	if len(thing) == 0 {
		return
	}
	e.Bot.Privmsg(e.Target, fmt.Sprintf("%s is %d", thing, s.Score(thing)))
}

// Register registers the plugin with a bot.
func Register(b bot.Bot, scoresfile string) {
	s := load(scoresfile)

	b.Conn().HandleFunc("privmsg",
		func(conn *client.Conn, line *client.Line) { parseScore(b, line, s) })

	b.Commands().Add("score", bot.Command{
		Help:    "score <thing> - show score of something",
		Handler: func(e *bot.Event) { showScore(e, s) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})

	b.Commands().Add("scores", bot.Command{
		Help:    "show top +/- scores",
		Handler: func(e *bot.Event) { e.Bot.Privmsg(e.Target, s.String()) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})

	// Every minute, save to file.
	if len(scoresfile) > 0 {
		go func() {
			for range time.Tick(time.Minute) {
				save(scoresfile, s)
			}
		}()
	}
}
