// Package dnssec implements a plugin to analyzes a domain with Verisign Labs DNSSEC analyzer.
package dnssec

import (
	"fmt"
	"strings"

	"github.com/StalkR/dnssec-analyzer/dnssec"
	"github.com/StalkR/goircbot/bot"
)

func handle(e *bot.Event) {
	site := strings.TrimSpace(e.Args)
	if len(site) == 0 {
		return
	}
	analysis, err := dnssec.Analyze(site)
	if err != nil {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("error: %v", err))
		return
	}

	var all []string
	for _, domain := range analysis {
		warnings, errors := 0, 0
		for _, result := range domain.Results {
			switch result.Status {
			case dnssec.WARNING:
				warnings++
			case dnssec.ERROR:
				errors++
			}
		}
		var status []string
		if warnings > 0 {
			status = append(status, fmt.Sprintf("warnings: %v", warnings))
		}
		if errors > 0 {
			status = append(status, fmt.Sprintf("errors: %v", errors))
		}
		if len(status) == 0 {
			all = append(all, fmt.Sprintf("%v (OK)", domain.Name))
		} else {
			all = append(all, fmt.Sprintf("%v (%v)", domain.Name, strings.Join(status, ", ")))
		}
	}

	e.Bot.Privmsg(e.Target, fmt.Sprintf("%v - %v%v", strings.Join(all, ", "), dnssec.URL, site))
}

// Register registers the plugin with a bot.
func Register(b bot.Bot) {
	b.Commands().Add("dnssec", bot.Command{
		Help:    "analyzes a domain with Verisign Labs DNSSEC analyzer",
		Handler: handle,
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
