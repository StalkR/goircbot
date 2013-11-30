// Package travisci implements a plugin to get and monitor Travis-CI builds.
package travisci

import (
	"fmt"
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/travisci"
)

const noResult = "no such user/repo or no build yet: https://www.travis-ci.org/%v/%v"

func travis(e *bot.Event) {
	userRepo := strings.SplitN(e.Args, "/", 2)
	if len(userRepo) != 2 {
		e.Bot.Privmsg(e.Target, "usage: travis <user>/<repo>")
		return
	}
	user, repo := userRepo[0], userRepo[1]
	builds, err := travisci.Builds(user, repo)
	if err != nil {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	if len(builds) == 0 {
		e.Bot.Privmsg(e.Target, fmt.Sprintf(noResult, user, repo))
		return
	}
	last := builds[0]
	e.Bot.Privmsg(e.Target, last.String())
}

// Register registers the plugin with a bot.
func Register(b bot.Bot) {
	b.Commands().Add("travis", bot.Command{
		Help:    "get build status of a user/repo on http://travis-ci.org",
		Handler: travis,
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
