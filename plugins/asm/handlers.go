// Package asm is a plugin to explain x86/x86-64 assembly instructions.
//go:generate go run generate/mnemonics.go -p asm -v mnemonics -o mnemonics.go
//go:generate gofmt -w mnemonics.go
package asm

import (
	"fmt"
	"strings"

	"github.com/StalkR/goircbot/bot"
)

func explain(e *bot.Event) {
	instr := strings.ToUpper(strings.TrimSpace(e.Args))
	if len(instr) == 0 {
		return
	}
	desc, ok := mnemonics[instr]
	if !ok {
		e.Bot.Privmsg(e.Target, "not found")
		return
	}
	e.Bot.Privmsg(e.Target, fmt.Sprintf("%s: %s", instr, desc))
}

// Register registers the plugin with a bot.
func Register(b bot.Bot) {
	b.Commands().Add("asm", bot.Command{
		Help:    "explain an x86/x86-64 assembly instruction",
		Handler: explain,
		Pub:     true,
		Priv:    false,
		Hidden:  false})
}
