// Package ping implements a plugin to ping a host or IP using system ping.
// It is made for ping from iputils (options and output parsing).
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
func Ping(host string) (string, error) {
	matched, err := regexp.Match("^[\\w._-]+$", []byte(host))
	if err != nil {
		return "", err
	}
	if !matched {
		return "", errors.New("invalid host/IP")
	}
	// -c: packet count, -w: timeout in seconds
	out, err := exec.Command("ping", "-c", "1", "-w", "3", "--", host).Output()
	if err != nil {
		if fmt.Sprintf("%s", err) == "exit status 2" {
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
		return "", errors.New("timeout")
	}
	return string(line), nil
}

func PingHandler(b *bot.Bot, e *bot.Event) {
	arg := strings.TrimSpace(e.Args)
	if len(arg) == 0 {
		return
	}
	s, err := Ping(arg)
	if err != nil {
		b.Conn.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	b.Conn.Privmsg(e.Target, s)
}

// Register registers the plugin with a bot.
func Register(b *bot.Bot) {
	b.AddCommand("ping", bot.Command{
		Help:    "ping a host/IP",
		Handler: PingHandler,
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
