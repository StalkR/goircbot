// Package metal implements a plugin to get metal band information.
package metal

import (
	"fmt"
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/spiritofmetal"
)

func metal(e *bot.Event) {
	name := strings.TrimSpace(e.Args)
	if len(name) == 0 {
		return
	}
	bands, err := spiritofmetal.Search(name)
	if err != nil {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	if len(bands) == 0 {
		e.Bot.Privmsg(e.Target, "no band found")
		return
	}
	var found []string
	for _, band := range bands {
		found = append(found, band.String())
	}
	e.Bot.Privmsg(e.Target, strings.Join(found, ", "))
}

// Register registers the plugin with a bot.
func Register(b bot.Bot) {
	b.Commands().Add("metal", bot.Command{
		Help:    "get metal band information",
		Handler: metal,
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
