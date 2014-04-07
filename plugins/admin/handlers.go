// Package admin implements a plugin in which registered admins can instruct the
// bot to perform commands (say, act, notice, op, deop, voice, devoice, quit).
package admin

import (
	"fmt"
	"log"
	"strings"

	"github.com/StalkR/goircbot/bot"
)

// Admins allowed to use commands in the form nick!ident@host
var Admins []string

func authorized(e *bot.Event) bool {
	for _, admin := range Admins {
		if e.Line.Src == admin {
			return true
		}
	}
	log.Printf("admin: %s not authorized", e.Line.Src)
	return false
}

func extractArgs(args string) (target, text string, err bool) {
	words := strings.SplitN(args, " ", 2)
	if len(words) < 2 {
		err = false
		return
	}
	target, text = words[0], words[1]
	return
}

func say(e *bot.Event) {
	if !authorized(e) {
		return
	}
	if target, text, err := extractArgs(e.Args); !err {
		e.Bot.Privmsg(target, text)
	}
}

func act(e *bot.Event) {
	if !authorized(e) {
		return
	}
	if target, text, err := extractArgs(e.Args); !err {
		e.Bot.Action(target, text)
	}
}

func notice(e *bot.Event) {
	if !authorized(e) {
		return
	}
	if target, text, err := extractArgs(e.Args); !err {
		e.Bot.Notice(target, text)
	}
}

func doMode(e *bot.Event, sign, mode string) {
	if !authorized(e) {
		return
	}
	var channel string
	args := strings.Split(e.Args, " ")
	if strings.HasPrefix(args[0], "#") {
		// Explicit request with a channel argument.
		channel = args[0]
		args = args[1:]
	} else if strings.HasPrefix(e.Line.Args[0], "#") {
		// Implicit request on a channel.
		channel = e.Line.Args[0]
	} else {
		// Private query, not applicable.
		return
	}
	if len(args) <= 0 || args[0] == "" {
		args = []string{e.Line.Nick}
	}
	e.Bot.Mode(channel, fmt.Sprintf("%s%s %s", sign,
		strings.Repeat(mode, len(args)), strings.Join(args, " ")))
}

func nick(e *bot.Event) {
	if !authorized(e) || e.Args == "" {
		return
	}
	e.Bot.Conn().Nick(e.Args)
}

func quit(e *bot.Event) {
	if !authorized(e) {
		return
	}
	e.Bot.Quit(e.Args)
}

// Register registers the plugin with a bot.
func Register(b bot.Bot, admins []string) {
	Admins = admins

	b.Commands().Add("say", bot.Command{
		Help:    "say <target> <text>",
		Handler: say,
		Pub:     true,
		Priv:    true,
		Hidden:  true,
	})
	b.Commands().Add("act", bot.Command{
		Help:    "act <target> <text>",
		Handler: act,
		Pub:     true,
		Priv:    true,
		Hidden:  true,
	})
	b.Commands().Add("notice", bot.Command{
		Help:    "notice <target> <text>",
		Handler: notice,
		Pub:     true,
		Priv:    true,
		Hidden:  true,
	})
	b.Commands().Add("op", bot.Command{
		Help:    "op [<target>]",
		Handler: func(e *bot.Event) { doMode(e, "+", "o") },
		Pub:     true,
		Priv:    true,
		Hidden:  true,
	})
	b.Commands().Add("deop", bot.Command{
		Help:    "deop [<target>]",
		Handler: func(e *bot.Event) { doMode(e, "-", "o") },
		Pub:     true,
		Priv:    true,
		Hidden:  true,
	})
	b.Commands().Add("voice", bot.Command{
		Help:    "voice [<target>]",
		Handler: func(e *bot.Event) { doMode(e, "+", "v") },
		Pub:     true,
		Priv:    true,
		Hidden:  true,
	})
	b.Commands().Add("devoice", bot.Command{
		Help:    "devoice [<target>]",
		Handler: func(e *bot.Event) { doMode(e, "-", "v") },
		Pub:     true,
		Priv:    true,
		Hidden:  true,
	})
	b.Commands().Add("nick", bot.Command{
		Help:    "nick <nickname>",
		Handler: nick,
		Pub:     true,
		Priv:    true,
		Hidden:  true,
	})
	b.Commands().Add("quit", bot.Command{
		Help:    "quit [msg]",
		Handler: quit,
		Pub:     true,
		Priv:    true,
		Hidden:  true,
	})
}
