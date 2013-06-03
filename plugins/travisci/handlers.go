// Package travisci implements a plugin to get and monitor Travis-CI builds.
package travisci

import (
	"fmt"
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/travisci"
)

const noResult = "no such user/repo or no build yet: https://www.travis-ci.org/%v/%v"

func travis(b *bot.Bot, e *bot.Event) {
	userRepo := strings.SplitN(e.Args, "/", 2)
	if len(userRepo) != 2 {
		b.Conn.Privmsg(e.Target, "usage: travis <user>/<repo>")
		return
	}
	user, repo := userRepo[0], userRepo[1]
	builds, err := travisci.Builds(user, repo)
	if err != nil {
		b.Conn.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	if len(builds) == 0 {
		b.Conn.Privmsg(e.Target, fmt.Sprintf(noResult, user, repo))
		return
	}
	last := builds[0]
	var status string
	if last.State == "finished" {
		status = "passed"
		if !last.Success {
			status = "errored"
		}
	} else {
		status = "in progress"
	}
	if last.Finished.IsZero() {
		b.Conn.Privmsg(e.Target, fmt.Sprintf("Build #%v: %v, %v (%v/%v)\n",
			last.Number, status, last.Message, last.Branch, last.Commit))
	} else {
		b.Conn.Privmsg(e.Target, fmt.Sprintf("Build #%v: %v (%v) %v (%v/%v)\n",
			last.Number, status, last.Finished, last.Message, last.Branch, last.Commit))
	}

}

// Register registers the plugin with a bot.
func Register(b *bot.Bot) {
	b.AddCommand("travis", bot.Command{
		Help:    "get build status of a user/repo on http://travis-ci.org",
		Handler: travis,
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
