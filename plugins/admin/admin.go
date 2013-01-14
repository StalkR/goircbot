// Package admin implements a plugin in which registered admins can instruct the
// bot to perform commands (say, act, notice, op, deop, voice, devoice, quit).
package admin

import (
	"fmt"
	bot "github.com/StalkR/goircbot"
	"log"
	"strings"
)

// Admins allowed to use commands in the form nick!ident@host
var Admins []string

func Authorized(e *bot.Event) bool {
	for _, admin := range Admins {
		if e.Line.Src == admin {
			return true
		}
	}
	log.Println("admin: not authorized", e.Line.Src)
	return false
}

func ExtractArgs(args string) (target, text string, err bool) {
	words := strings.SplitN(args, " ", 2)
	if len(words) < 2 {
		err = false
		return
	}
	target, text = words[0], words[1]
	return
}

func Say(b *bot.Bot, e *bot.Event) {
	if !Authorized(e) {
		return
	}
	if target, text, err := ExtractArgs(e.Args); !err {
		b.Conn.Privmsg(target, text)
	}
}

func Act(b *bot.Bot, e *bot.Event) {
	if !Authorized(e) {
		return
	}
	if target, text, err := ExtractArgs(e.Args); !err {
		b.Conn.Action(target, text)
	}
}

func Notice(b *bot.Bot, e *bot.Event) {
	if !Authorized(e) {
		return
	}
	if target, text, err := ExtractArgs(e.Args); !err {
		b.Conn.Notice(target, text)
	}
}

func DoMode(b *bot.Bot, e *bot.Event, sign, mode string) {
	if !Authorized(e) {
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

func Quit(b *bot.Bot, e *bot.Event) {
	if !Authorized(e) {
		return
	}
	b.Reconnect = false
	b.Conn.Quit(e.Args)
}

// Register registers the plugin with a bot.
func Register(b *bot.Bot, admins []string) {
	Admins = admins

	b.AddCommand("say", bot.Command{"say <target> <text>", Say, true, true, true})
	b.AddCommand("act", bot.Command{"act <target> <text>", Act, true, true, true})
	b.AddCommand("notice", bot.Command{"notice <target> <text>", Notice, true, true, true})

	b.AddCommand("op", bot.Command{"op [<target>]",
		func(b *bot.Bot, e *bot.Event) { DoMode(b, e, "+", "o") },
		true, true, true})
	b.AddCommand("deop", bot.Command{"deop [<target>]",
		func(b *bot.Bot, e *bot.Event) { DoMode(b, e, "-", "o") },
		true, true, true})
	b.AddCommand("voice", bot.Command{"voice [<target>]",
		func(b *bot.Bot, e *bot.Event) { DoMode(b, e, "+", "v") },
		true, true, true})
	b.AddCommand("devoice", bot.Command{"devoice [<target>]",
		func(b *bot.Bot, e *bot.Event) { DoMode(b, e, "-", "v") },
		true, true, true})

	b.AddCommand("quit", bot.Command{"quit [msg]", Quit, true, true, true})
}
