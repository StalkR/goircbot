// Package idle implements a plugin to see idle times on channels.
package idle

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/StalkR/goircbot/bot"
	"github.com/fluffle/goirc/client"
)

func topidle(e *bot.Event, ignore []string) {
	st := e.Bot.Conn().StateTracker()
	ch := st.GetChannel(e.Target)
	if ch == nil {
		return
	}
	ignoremap := make(map[string]bool)
	for _, nick := range ignore {
		ignoremap[nick] = true
	}
	var nicks []string
	for nick := range ch.Nicks {
		if n := st.GetNick(nick); n == nil || n.Modes.Bot {
			continue
		}
		if _, present := ignoremap[nick]; present {
			continue
		}
		nicks = append(nicks, nick)
	}
	if len(nicks) == 0 {
		return
	}

	done := make(chan struct{})
	m := make(map[string]time.Duration)
	defer e.Bot.Conn().HandleFunc("317",
		func(conn *client.Conn, line *client.Line) {
			r := strings.Split(line.Raw, " ")
			if len(r) < 5 {
				return
			}
			seconds, err := strconv.Atoi(r[4])
			if err != nil {
				return
			}
			m[r[3]] = time.Duration(seconds) * time.Second
			if len(m) == len(nicks) {
				close(done)
			}
		}).Remove()

	for _, nick := range nicks {
		e.Bot.Conn().Whois(nick)
	}

	<-done
	if len(m) == 0 {
		return
	}

	var idlers []idler
	for nick, idle := range m {
		idlers = append(idlers, idler{Nick: nick, Idle: idle})
	}
	sort.Sort(byIdle(idlers))
	top := 3
	if len(idlers) < top {
		top = len(idlers)
	}
	for i, v := range idlers[:top] {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("%d - %s idle %s", i+1, v.Nick, v.Idle.String()))
	}
}

type idler struct {
	Nick string
	Idle time.Duration
}

type byIdle []idler

func (a byIdle) Len() int           { return len(a) }
func (a byIdle) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byIdle) Less(i, j int) bool { return a[i].Idle > a[j].Idle }

// Register registers the plugin with a bot.
// Use ignore as a list of nicks to ignore.
func Register(b bot.Bot, ignore []string) {
	b.Commands().Add("topidle", bot.Command{
		Help:    "top idlers on the channel",
		Handler: func(e *bot.Event) { topidle(e, ignore) },
		Pub:     true,
		Priv:    false,
		Hidden:  false})
}
