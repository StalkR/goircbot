// Package asm is a plugin to explain x86/x86-64 assembly instructions.
package asm

import (
	"fmt"
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/x86"
)

func explain(e *bot.Event, mnemonics map[string]string) {
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
func Register(b bot.Bot) error {
	ref, err := x86.New()
	if err != nil {
		return err
	}
	mnemonics := ref.Mnemonics()

	b.Commands().Add("asm", bot.Command{
		Help:    "explain an x86/x86-64 assembly instruction",
		Handler: func(e *bot.Event) { explain(e, mnemonics) },
		Pub:     true,
		Priv:    false,
		Hidden:  false})
	return nil
}
