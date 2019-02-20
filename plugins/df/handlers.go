// Package df implements a plugin to monitor disk usage and notify when low.
package df

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/disk"
	"github.com/StalkR/goircbot/lib/size"
)

const delay = 5 * time.Minute

// An Alarm represents a path and disk size limit (in bytes) when to notify.
// Notification happens once when limit is crossed then is reset if
// it goes back above the limit.
type Alarm struct {
	Path     string
	Limit    uint64
	notified bool
}

// NewAlarm creates a new Alarm as specified.
func NewAlarm(path string, limit size.Byte) Alarm {
	return Alarm{Path: path, Limit: uint64(limit)}
}

// Monitor monitors a path and notifies when limit is crossed.
func (a *Alarm) Monitor(b bot.Bot) {
	for ; ; time.Sleep(delay) {
		total, free, err := disk.Space(a.Path)
		if err != nil {
			log.Printf("df: error: %v", err)
			continue
		}
		if free > a.Limit {
			a.notified = false
			continue
		}
		if !a.notified {
			a.Notify(b, total, free)
		}
	}
}

// status reports disk free/used status of a path.
func status(path string) (string, error) {
	total, free, err := disk.Space(path)
	if err != nil {
		return "", err
	}
	percent := 100 * (total - free) / total
	return fmt.Sprintf("%v has %v free (%v%% of %v used)",
		path, size.Byte(free), percent, size.Byte(total)), nil
}

// Notify notifies disk usage on all channels.
func (a *Alarm) Notify(b bot.Bot, total, free uint64) {
	if !b.Connected() {
		return
	}
	a.notified = true
	percent := 100 * (total - free) / total
	line := fmt.Sprintf("Warning: %v has %v free (%v%% of %v used)",
		a.Path, size.Byte(free), percent, size.Byte(total))
	for _, channel := range b.Channels() {
		b.Privmsg(channel, line)
	}
}

func df(e *bot.Event, paths map[string]bool) {
	if len(paths) == 0 {
		e.Bot.Privmsg(e.Target, "error: no paths configured")
		return
	}
	var pathlist []string
	for p := range paths {
		pathlist = append(pathlist, p)
	}

	path := strings.TrimSpace(e.Args)
	if path == "" {
		if len(paths) > 1 {
			e.Bot.Privmsg(e.Target, fmt.Sprintf("usage: df <path> (%v)", strings.Join(pathlist, ", ")))
			return
		}
		path = pathlist[0]
	}
	if !paths[path] {
		e.Bot.Privmsg(e.Target, "error: path not configured")
		return
	}
	line, err := status(path)
	if err != nil {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("error: %v", err))
		return
	}
	e.Bot.Privmsg(e.Target, line)
}

// Register registers the plugin with a bot.
func Register(b bot.Bot, alarms ...Alarm) {
	paths := map[string]bool{}
	for _, a := range alarms {
		paths[a.Path] = true
		go a.Monitor(b)
	}

	b.Commands().Add("df", bot.Command{
		Help:    "See disk usage",
		Handler: func(e *bot.Event) { df(e, paths) },
		Pub:     true,
		Priv:    false,
		Hidden:  false})
}
