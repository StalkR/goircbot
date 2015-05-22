// Package errors implements a plugin to get Linux & Windows error and status information.
package errors

//go:generate go run gen/build.go -f gen/winerrors.txt -v winerrors -p errors -o winerrors.go
//go:generate go run gen/build.go -f gen/ntstatus.txt -v ntstatus -p errors -o ntstatus.go
//go:generate go run gen/build.go -f gen/errnos.txt -v errnos -p errors -o errnos.go

import (
	"path/filepath"
	"strings"

	"github.com/StalkR/goircbot/bot"
)

func handle(e *bot.Event, table [][]string) {
	arg := strings.TrimSpace(e.Args)
	if len(arg) == 0 {
		return
	}
	e.Bot.Privmsg(e.Target, find(table, arg))
}

// Register registers the plugin with a bot.
func Register(b bot.Bot, dir string) error {
	winerrors, err := parse(filepath.Join(dir, "winerrors.txt"))
	if err != nil {
		return err
	}
	ntstatus, err := parse(filepath.Join(dir, "ntstatus.txt"))
	if err != nil {
		return err
	}
	errnos, err := parse(filepath.Join(dir, "errnos.txt"))
	if err != nil {
		return err
	}

	b.Commands().Add("error", bot.Command{
		Help:    "get Windows error code information",
		Handler: func(e *bot.Event) { handle(e, winerrors) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})
	b.Commands().Add("status", bot.Command{
		Help:    "get Windows status code information",
		Handler: func(e *bot.Event) { handle(e, ntstatus) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})
	b.Commands().Add("errno", bot.Command{
		Help:    "get Windows status code information",
		Handler: func(e *bot.Event) { handle(e, errnos) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})
	return nil
}
