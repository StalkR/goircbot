// Package admin implements a plugin in which registered admins can instruct the
// bot to perform commands (say, act, notice, op, deop, voice, devoice, quit).
package admin

import (
	"fmt"
	bot "goircbot"
	"log"
	"strings"
)

// Admins allowed to use commands in the form nick!ident@host
var Admins []string

func authorized(e *bot.Event) bool {
	for _, admin := range Admins {
		if e.Line.Src == admin {
			return true
		}
	}
	log.Println("admin: not authorized", e.Line.Src)
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

func say(b *bot.Bot, e *bot.Event) {
	if !authorized(e) {
		return
	}
	if target, text, err := extractArgs(e.Args); !err {
		b.Conn.Privmsg(target, text)
	}
}

func act(b *bot.Bot, e *bot.Event) {
	if !authorized(e) {
		return
	}
	if target, text, err := extractArgs(e.Args); !err {
		b.Conn.Action(target, text)
	}
}

func notice(b *bot.Bot, e *bot.Event) {
	if !authorized(e) {
		return
	}
	if target, text, err := extractArgs(e.Args); !err {
		b.Conn.Notice(target, text)
	}
}

func doMode(b *bot.Bot, e *bot.Event, sign, mode string) {
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
	b.Conn.Mode(channel, fmt.Sprintf("%s%s %s", sign,
		strings.Repeat(mode, len(args)), strings.Join(args, " ")))
}

func quit(b *bot.Bot, e *bot.Event) {
	if !authorized(e) {
		return
	}
	b.Reconnect = false
	b.Conn.Quit(e.Args)
}

// Register registers the plugin with a bot.
func Register(b *bot.Bot, admins []string) {
	Admins = admins

	b.AddCommand("say", bot.Command{"say <target> <text>", say, true, true, true})
	b.AddCommand("act", bot.Command{"act <target> <text>", act, true, true, true})
	b.AddCommand("notice", bot.Command{"notice <target> <text>", notice, true, true, true})

	b.AddCommand("op", bot.Command{"op [<target>]",
		func(b *bot.Bot, e *bot.Event) { doMode(b, e, "+", "o") },
		true, true, true})
	b.AddCommand("deop", bot.Command{"deop [<target>]",
		func(b *bot.Bot, e *bot.Event) { doMode(b, e, "-", "o") },
		true, true, true})
	b.AddCommand("voice", bot.Command{"voice [<target>]",
		func(b *bot.Bot, e *bot.Event) { doMode(b, e, "+", "v") },
		true, true, true})
	b.AddCommand("devoice", bot.Command{"devoice [<target>]",
		func(b *bot.Bot, e *bot.Event) { doMode(b, e, "-", "v") },
		true, true, true})

	b.AddCommand("quit", bot.Command{"quit [msg]", quit, true, true, true})
}
