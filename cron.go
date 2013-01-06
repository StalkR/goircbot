package goircbot

import (
	"log"
	"time"
)

// CronHandler is a callback used by plugins to have scheduled tasks.
type CronHandler func(b *Bot)

// Cron is a scheduled task that plugins can add.
type Cron struct {
	Handler  CronHandler
	Duration time.Duration
}

// AddCron is a convenient way for plugins to add scheduled tasks.
// Task starts immediately with a first run, then run at every duration.
// Task does not execute when disconnected but keeps running.
func (b *Bot) AddCron(name string, c Cron) {
	if _, present := b.crons[name]; present {
		log.Println("AddCron: already defined", name)
		return
	}
	b.crons[name] = c
	// Start scheduled task in a goroutine that detects when stopped.
	go func() {
		for {
			if cron, present := b.crons[name]; present {
				if b.Conn.Connected {
					cron.Handler(b)
				}
				time.Sleep(cron.Duration)
			} else {
				return
			}
		}
	}()
}

// DelCron is to remove scheduled tasks added via AddCron.
// Task is not guaranteed to stop immediately but before the next run.
func (b *Bot) DelCron(name string) {
	if _, present := b.crons[name]; !present {
		log.Println("AddCron: not defined", name)
		return
	}
	delete(b.crons, name)
}
