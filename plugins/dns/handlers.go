// Package dns implements a plugin to query DNS of a host/IP.
package dns

import (
	"fmt"
	"net"
	"strings"

	"github.com/StalkR/goircbot/bot"
)

// Resolve resolves a host/IP with DNS and returns a list of results.
func Resolve(h string) ([]string, error) {
	results := []string{}
	if net.ParseIP(h) != nil {
		name, err := net.LookupAddr(h)
		if err != nil {
			return nil, err
		}
		results = name
	} else {
		addrs, err := net.LookupHost(h)
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			results = append(results, strings.TrimRight(addr, "."))
		}
	}
	return results, nil
}

func dns(b *bot.Bot, e *bot.Event) {
	arg := strings.TrimSpace(e.Args)
	if len(arg) == 0 {
		return
	}
	results, err := Resolve(arg)
	if err != nil {
		b.Conn.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	if len(results) == 0 {
		b.Conn.Privmsg(e.Target, "Not found.")
		return
	}
	max := len(results)
	if max > 20 {
		max = 20
	}
	reply := fmt.Sprintf("%s: %s", arg, strings.Join(results[:max], ", "))
	b.Conn.Privmsg(e.Target, reply)
}

// Register registers the plugin with a bot.
func Register(b *bot.Bot) {
	b.AddCommand("dns", bot.Command{
		Help:    "DNS resolve a host or IP",
		Handler: dns,
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
