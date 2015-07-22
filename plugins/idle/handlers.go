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

const timeout = 10 * time.Second

func getIdle(e *bot.Event, nicks []string) []idler {
	nickSet := make(map[string]struct{})
	for _, nick := range nicks {
		// IRC nicks are case insensitive
		nickSet[strings.ToLower(nick)] = struct{}{}
	}

	infos := make(map[string]time.Duration)
	received := make(map[string]struct{})
	done := make(chan struct{})
	defer e.Bot.Conn().HandleFunc("317", // Idle & sign on time
		func(conn *client.Conn, line *client.Line) {
			r := strings.Split(line.Raw, " ")
			if len(r) < 5 {
				return
			}
			seconds, err := strconv.Atoi(r[4])
			if err != nil {
				return
			}
			infos[r[3]] = time.Duration(seconds) * time.Second
		}).Remove()
	defer e.Bot.Conn().HandleFunc("318", // End of whois
		func(conn *client.Conn, line *client.Line) {
			r := strings.Split(line.Raw, " ")
			if len(r) < 4 {
				return
			}
			received[r[3]] = struct{}{}
			if len(received) == len(nickSet) {
				close(done)
			}
		}).Remove()

	// Temporarily disable flood protection or it's too slow
	e.Bot.Conn().Config().Flood = true
	defer func() {
		e.Bot.Conn().Config().Flood = false
	}()

	for nick := range nickSet {
		e.Bot.Conn().Whois(nick)
	}
	select {
	case <-time.After(timeout):
	case <-done:
	}

	var idlers []idler
	for nick, idle := range infos {
		// check if we did request idle for that nick otherwise ignore
		if _, ok := nickSet[strings.ToLower(nick)]; !ok {
			continue
		}
		idlers = append(idlers, idler{Nick: nick, Idle: idle})
	}
	sort.Sort(byIdle(idlers))
	return idlers
}

func handleIdle(e *bot.Event) {
	nick := strings.TrimSpace(e.Args)
	if len(nick) == 0 || strings.Contains(nick, " ") {
		return
	}
	idlers := getIdle(e, []string{nick})
	if len(idlers) != 1 {
		return
	}
	e.Bot.Privmsg(e.Target, fmt.Sprintf("%s idle %s", idlers[0].Nick, idlers[0].Idle.String()))
}

func handleTopidle(e *bot.Event, ignore []string) {
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
	idlers := getIdle(e, nicks)
	if len(idlers) == 0 {
		return
	}
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
	b.Commands().Add("idle", bot.Command{
		Help:    "get idle time of a user",
		Handler: func(e *bot.Event) { handleIdle(e) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})
	b.Commands().Add("topidle", bot.Command{
		Help:    "top idlers on the channel",
		Handler: func(e *bot.Event) { handleTopidle(e, ignore) },
		Pub:     true,
		Priv:    false,
		Hidden:  false})
}
