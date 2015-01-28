package travisci

import (
	"log"
	"strings"
	"time"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/travisci"
)

func watch(user, repo string, duration time.Duration, notify func(string)) {
	lastBuild := 0
	for ; ; time.Sleep(duration) {
		builds, err := travisci.Builds(user, repo)
		if err != nil {
			log.Printf("travisci: error watching %v/%v builds: %v\n", user, repo, err)
			continue
		}
		if len(builds) == 0 {
			log.Printf("travisci: %v/%v has no build yet\n", user, repo)
			continue
		}
		if lastBuild == 0 {
			lastBuild = builds[0].Number
			continue
		}
		for i, j := 0, len(builds)-1; i < j; i, j = i+1, j-1 {
			builds[i], builds[j] = builds[j], builds[i]
		}
		for _, b := range builds {
			if b.Number <= lastBuild || b.State != "finished" {
				continue
			}
			lastBuild = b.Number
			if b.Success {
				continue
			}
			notify(b.String())
		}
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

// Watch registers a watcher of user/repos with a bot.
// If any build fails, it will be announced on all channels.
func Watch(b bot.Bot, userRepos []string, duration time.Duration) {
	for _, arg := range userRepos {
		userRepo := strings.SplitN(arg, "/", 2)
		if len(userRepo) != 2 {
			panic("invalid user/repo: " + arg)
		}
		go watch(userRepo[0], userRepo[1], duration, func(line string) { notify(b, line) })
	}
}
