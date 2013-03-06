// Package ping implements a plugin to ping a host or IP (v4/v6) using system ping.
// It is made for ping/ping6 from iputils (options and output parsing).
package ping

import (
	"errors"
	"fmt"
	bot "github.com/StalkR/goircbot"
	"os/exec"
	"regexp"
	"strings"
)

// Ping runs ping against given host and returns its output.
func Ping(host string, ipv6 bool) (string, error) {
	matched, err := regexp.Match("^[\\w._-]+$", []byte(host))
	if err != nil {
		return "", err
	}
	if !matched {
		return "", errors.New("invalid host/IP")
	}
	six := ""
	if ipv6 {
		six = "6"
	}
	// -c: packet count, -w: timeout in seconds
	out, err := exec.Command("ping"+six, "-c", "1", "-w", "3", "--", host).Output()
	if err != nil {
		errs := fmt.Sprintf("%s", err)
		if errs == "exit status 1" {
			return "", errors.New("timeout")
		}
		if errs == "exit status 2" {
			return "", errors.New("unknown host")
		}
		return "", err
	}
	r, err := regexp.Compile("\\d+ bytes from .*")
	if err != nil {
		return "", err
	}
	line := r.Find(out)
	if line == nil {
		return "", errors.New("cannot parse ping output")
	}
	return string(line), nil
}

func PingHandler(b *bot.Bot, e *bot.Event, ipv6 bool) {
	arg := strings.TrimSpace(e.Args)
	if len(arg) == 0 {
		return
	}
	s, err := Ping(arg, ipv6)
	if err != nil {
		b.Conn.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	b.Conn.Privmsg(e.Target, s)
}

// Register registers the plugin with a bot.
func Register(b *bot.Bot) {
	b.AddCommand("ping", bot.Command{
		Help:    "ping a host/IPv4",
		Handler: func(b *bot.Bot, e *bot.Event) { PingHandler(b, e, false) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})
	b.AddCommand("ping6", bot.Command{
		Help:    "ping a host/IPv6",
		Handler: func(b *bot.Bot, e *bot.Event) { PingHandler(b, e, true) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
