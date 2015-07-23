package git

import (
	"log"
	"time"

	"github.com/StalkR/goircbot/bot"
)

func watch(repo string, duration time.Duration, notify func(string)) {
	last := ""
	for ; ; time.Sleep(duration) {
		commit, err := lastLog(repo)
		if err != nil {
			log.Printf("git: error watching %s: %s", repo, err)
			continue
		}
		if commit.Msg != last && last != "" {
			notify(commit.String())
		}
		last = commit.Msg
	}
}

func notify(b bot.Bot, line string) {
	if !b.Connected() {
		return
	}
	for _, channel := range b.Channels() {
		b.Privmsg(channel, line)
	}
}

// Watch registers a watcher of git repo commit log with a bot.
// Repo is a URL to a cgit log page such as
// https://git.kernel.org/cgit/linux/kernel/git/torvalds/linux.git/log/
func Watch(b bot.Bot, repo string, duration time.Duration) {
	go watch(repo, duration, func(line string) {
		notify(b, line)
	})
}
