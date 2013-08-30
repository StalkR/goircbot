// Package df implements a plugin to monitor disk usage and notify when low.
package df

import (
	"fmt"
	"log"
	"syscall"
	"time"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/size"
)

const delay = 5 * time.Minute

// An Alarm represents a path and disk size limit (in bytes) when to notify.
// Notification happens once when limit is crossed then is reset if
// it goes back above the limit.
type Alarm struct {
	Path     string
	Limit    int
	notified bool
}

// NewAlarm creates a new Alarm as specified.
func NewAlarm(path string, limit size.Byte) Alarm {
	return Alarm{Path: path, Limit: int(limit)}
}

// Monitor monitors a path and notifies when limit is crossed.
func (a *Alarm) Monitor(b *bot.Bot) {
	s := syscall.Statfs_t{}
	for ; ; time.Sleep(delay) {
		err := syscall.Statfs(a.Path, &s)
		if err != nil {
			log.Printf("diskusage: statfs error: %v", err)
			continue
		}
		total := int(s.Bsize) * int(s.Blocks)
		free := int(s.Bsize) * int(s.Bfree)
		if free > a.Limit {
			a.notified = false
			continue
		}
		if !a.notified {
			a.Notify(b, total, free)
		}
	}
}

// Notify notifies disk usage on all channels.
func (a *Alarm) Notify(b *bot.Bot, total, free int) {
	if !b.Conn.Connected() {
		return
	}
	a.notified = true
	percent := 100 * (total - free) / total
	totalFmt := size.Byte(total).String()
	freeFmt := size.Byte(free).String()
	line := fmt.Sprintf("Warning: %v has %v free (%v%% of %v used)",
		a.Path, freeFmt, percent, totalFmt)
	for _, channel := range b.Conn.Me().Channels() {
		b.Conn.Privmsg(channel.Name, line)
	}
}

// Register registers the plugin with a bot.
func Register(b *bot.Bot, alarms []Alarm) {
	for _, a := range alarms {
		go a.Monitor(b)
	}
}
