// Package mac implements a plugin to find the orgainzation behind a MAC
// address using IEEE public OUI.
//go:generate go run generate/ieeeoui.go -p mac -v ieeeoui -o ieeeoui.go
//go:generate gofmt -w ieeeoui.go
package mac

import (
	"errors"
	"strconv"
	"strings"

	"github.com/StalkR/goircbot/bot"
)

var ignoreChars = strings.NewReplacer("-", "", ":", "")

func find(address string) (string, error) {
	hex := ignoreChars.Replace(address)
	if len(hex) < 6 {
		return "", errors.New("need at least 3 bytes of address")
	}
	oui, err := strconv.ParseUint(hex[:6], 16, 0)
	if err != nil {
		return "", errors.New("invalid address")
	}
	org, ok := ieeeoui[oui]
	if !ok {
		return "", errors.New("not found")
	}
	return org, nil
}

func mac(e *bot.Event) {
	arg := strings.TrimSpace(e.Args)
	if len(arg) == 0 {
		return
	}
	org, err := find(arg)
	if err != nil {
		e.Bot.Privmsg(e.Target, err.Error())
		return
	}
	e.Bot.Privmsg(e.Target, org)
}

// Register registers the plugin with a bot.
func Register(b bot.Bot) {
	b.Commands().Add("mac", bot.Command{
		Help:    "find manufacturer of a MAC address (IEEE public OUI)",
		Handler: mac,
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
