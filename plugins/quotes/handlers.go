// Package quotes implements a plugin to save and replay quotes.
package quotes

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/nohl"
)

func add(e *bot.Event, q *Quotes) {
	text := strings.TrimSpace(e.Args)
	if len(text) == 0 {
		return
	}
	a := q.Add(e.Line.Nick, text)
	e.Bot.Privmsg(e.Target, fmt.Sprintf("added #%d", a.ID))
}

func del(e *bot.Event, q *Quotes) {
	id, err := strconv.Atoi(e.Args)
	if err != nil {
		e.Bot.Privmsg(e.Target, err.Error())
		return
	}
	if !q.Delete(id) {
		e.Bot.Privmsg(e.Target, "not found")
		return
	}
	e.Bot.Privmsg(e.Target, fmt.Sprintf("deleted #%d", id))
}

func search(e *bot.Event, q *Quotes) {
	term := strings.TrimSpace(e.Args)
	if len(term) == 0 {
		if q.Empty() {
			e.Bot.Privmsg(e.Target, "no quotes yet")
			return
		}
		e.Bot.Privmsg(e.Target, nohl.String(e.Bot, e.Target, q.Random().String()))
		return
	}
	if id, err := strconv.Atoi(term[1:]); strings.HasPrefix(term, "#") && err == nil {
		r, found := q.Get(id)
		if !found {
			e.Bot.Privmsg(e.Target, "not found")
			return
		}
		e.Bot.Privmsg(e.Target, nohl.String(e.Bot, e.Target, r.String()))
		return
	}
	lastResults = q.Search(term)
	lastPos = 0
	if len(lastResults) == 0 {
		e.Bot.Privmsg(e.Target, "not found")
		return
	}
	next(e)
}

var (
	lastResults []Quote
	lastPos     int
)

func next(e *bot.Event) {
	if lastPos >= len(lastResults) {
		e.Bot.Privmsg(e.Target, "no more")
		return
	}
	s := fmt.Sprintf("%d/%d: %s", lastPos+1, len(lastResults), lastResults[lastPos].String())
	e.Bot.Privmsg(e.Target, nohl.String(e.Bot, e.Target, s))
	lastPos++
}

// Register registers the plugin with a bot.
func Register(b bot.Bot, f string) {
	q := load(f)

	b.Commands().Add("addquote", bot.Command{
		Help:    "add a quote",
		Handler: func(e *bot.Event) { add(e, q) },
		Pub:     true,
		Priv:    false,
		Hidden:  false})

	b.Commands().Add("delquote", bot.Command{
		Help:    "del a quote",
		Handler: func(e *bot.Event) { del(e, q) },
		Pub:     true,
		Priv:    false,
		Hidden:  false})

	b.Commands().Add("quote", bot.Command{
		Help:    "search or show a random quote",
		Handler: func(e *bot.Event) { search(e, q) },
		Pub:     true,
		Priv:    false,
		Hidden:  false})

	b.Commands().Add("next", bot.Command{
		Help:    "show the next quote after a search",
		Handler: next,
		Pub:     true,
		Priv:    false,
		Hidden:  false})

	// Every minute, save to file (effective only if changed).
	if len(f) > 0 {
		go func() {
			for range time.Tick(10 * time.Second) {
				save(f, q)
			}
		}()
	}
}
