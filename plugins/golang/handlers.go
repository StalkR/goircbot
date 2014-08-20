// Package golang implements a plugin to run Go snippets with http://play.golang.org.
package golang

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/transport"
)

// Register registers the plugin with a bot.
func Register(b bot.Bot) {
	b.Commands().Add("go", bot.Command{
		Help:    "run go code on http://play.golang.org",
		Handler: goCmd,
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}

func goCmd(e *bot.Event) {
	snippet := strings.TrimSpace(e.Args)
	if len(snippet) == 0 {
		return
	}
	// note: goirc takes care of cutting new lines
	out, err := run(snippet)
	if err != nil {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	e.Bot.Privmsg(e.Target, out)
}

const (
	fmtURL = "http://play.golang.org/fmt"
	runURL = "http://play.golang.org/compile"
)

type runResult struct {
	Errors string
	Events []event
}

type event struct {
	Message string
	Delay   int
}

func run(snippet string) (string, error) {
	code, err := goFmt(fmt.Sprintf("package main\nfunc main() {\n%s\n}", snippet))
	if err != nil {
		return "", err
	}
	c, err := transport.Client(runURL)
	if err != nil {
		return "", err
	}
	resp, err := c.PostForm(runURL, url.Values{
		"body":    {code},
		"version": {"2"},
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var r runResult
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", err
	}
	if r.Errors != "" {
		return "", fmt.Errorf("%s", r.Errors)
	}
	var out []string
	for _, e := range r.Events {
		out = append(out, e.Message)
	}
	return strings.Join(out, "\n"), nil
}

type fmtResult struct {
	Body, Error string
}

func goFmt(code string) (string, error) {
	c, err := transport.Client(fmtURL)
	if err != nil {
		return "", err
	}
	resp, err := c.PostForm(fmtURL, url.Values{
		"body":    {code},
		"imports": {"true"},
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var r fmtResult
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", err
	}
	if r.Error != "" {
		return "", fmt.Errorf("%s", r.Error)
	}
	return r.Body, nil
}
