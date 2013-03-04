// Package sed implements a plugin to replace pattern in sentences.
// When someone says s/pattern/replace/ bot replaces that someone's last line.
package sed

import (
	"container/list"
	"fmt"
	bot "github.com/StalkR/goircbot"
	irc "github.com/fluffle/goirc/client"
	"regexp"
	"strings"
	"sync"
	"time"
)

var expiration = 5 * time.Minute
var maxLines = 4

// Entry is a backlog entry: text line and time it happened.
type Entry struct {
	Line string
	Time time.Time
}

// Expired tells if an entry is expired based on time it happened and expiration.
func (e *Entry) Expired() bool {
	return time.Now().After(e.Time.Add(expiration))
}

// Backlog is an accessible map of channels to nick to entries.
type Backlog struct {
	sync.Mutex
	M map[string]map[string]*list.List
}

// Clean maintains a backlog clean by expiring old entries and ensuring maximum of lines.
func (bl *Backlog) Clean() {
	bl.Lock()
	for channel, cl := range bl.M {
		for nick, l := range cl {
			var rm []*list.Element
			i := 0
			for e := l.Front(); e != nil; e = e.Next() {
				entry := e.Value.(Entry)
				if entry.Expired() || i >= maxLines {
					rm = append(rm, e)
				} else {
					i++
				}
			}
			for _, e := range rm {
				l.Remove(e)
			}
			if l.Len() == 0 {
				delete(cl, nick)
			}
		}
		if len(cl) == 0 {
			delete(bl.M, channel)
		}
	}
	bl.Unlock()
}

// Store saves a line from a channel/nick into backlog.
func (bl *Backlog) Store(channel, nick, line string) {
	bl.Lock()
	e := Entry{line, time.Now()}
	if bl.M == nil {
		bl.M = map[string]map[string]*list.List{}
	}
	if _, p := bl.M[channel]; p {
		if _, q := bl.M[channel][nick]; q {
			// most recent line is in first position
			bl.M[channel][nick].PushFront(e)
		} else {
			l := list.New()
			l.PushBack(e)
			bl.M[channel][nick] = l
		}
	} else {
		l := list.New()
		l.PushBack(e)
		bl.M[channel] = map[string]*list.List{nick: l}
	}
	bl.Unlock()
	bl.Clean()
}

// Search iterates through backlog lines of a channel/nick.
func (bl *Backlog) Search(channel, nick string) chan string {
	c := make(chan string)
	go func() {
		bl.Lock()
		if _, p := bl.M[channel]; p {
			if _, q := bl.M[channel][nick]; q {
				l := bl.M[channel][nick]
				for e := l.Front(); e != nil; e = e.Next() {
					entry := e.Value.(Entry)
					c <- entry.Line
				}
			}
		}
		close(c)
		bl.Unlock()
	}()
	return c
}

// Sed attempts to replace a pattern in a backlog for channel/nick.
func (bl *Backlog) Sed(channel, nick, pattern, replace string) string {
	if len(pattern) > 80 {
		pattern = pattern[:80]
	}
	if len(replace) > 80 {
		replace = replace[:80]
	}
	for line := range bl.Search(channel, nick) {
		if strings.Contains(line, pattern) {
			r := strings.Replace(line, pattern, replace, 1)
			if len(r) > 160 {
				r = r[:160]
			}
			return r
		}
	}
	return ""
}

func watchLine(b *bot.Bot, line *irc.Line, bl *Backlog) {
	channel := line.Args[0]
	nick := line.Nick
	text := line.Args[1]
	if !strings.HasPrefix(channel, "#") {
		return
	}
	r, err := regexp.Compile("^s/([^/]+)/([^/]+)(?:/g?)?")
	if err != nil {
		b.Conn.Privmsg(channel, fmt.Sprintf("error: %s", err))
		return
	}
	m := r.FindSubmatch([]byte(text))
	if m == nil {
		bl.Store(channel, nick, text)
		return
	}
	meant := bl.Sed(channel, nick, string(m[1]), string(m[2]))
	if meant == "" {
		return
	}
	b.Conn.Privmsg(channel, fmt.Sprintf("%s meant: %s", nick, meant))
	bl.Store(channel, nick, meant)
}

// Register registers the plugin with a bot.
func Register(b *bot.Bot) {
	bl := &Backlog{}
	b.Conn.AddHandler("privmsg",
		func(conn *irc.Conn, line *irc.Line) { watchLine(b, line, bl) })
}
