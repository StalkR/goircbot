// Package coin implements a plugin to show coin prices from coinlayer API.
package coin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/transport"
)

const queryURL = "http://api.coinlayer.com/api/live"

type result struct {
	Success bool               `json:"success"`
	Rates   map[string]float64 `json:"rates"`
	Error   struct {
		Info string `json:"info"`
	} `json:"error"`
}

var symbolRE = regexp.MustCompile(`^[A-Z]{3}$`)

func rate(apiKey, symbol string) (string, error) {
	if !symbolRE.MatchString(symbol) {
		return "", errors.New("invalid symbol")
	}
	u := fmt.Sprintf("%s?%s", queryURL, url.Values{"access_key": []string{apiKey}}.Encode())
	c, err := transport.Client(u)
	if err != nil {
		return "", err
	}
	resp, err := c.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var v result
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return "", err
	}
	if !v.Success {
		return "", errors.New(v.Error.Info)
	}
	r, ok := v.Rates[symbol]
	if !ok {
		return "", fmt.Errorf("could not find rate for %v", symbol)
	}
	return fmt.Sprintf("%v", r), nil
}

func coin(e *bot.Event, apiKey string) {
	symbol := strings.ToUpper(strings.TrimSpace(e.Args))
	if len(symbol) == 0 {
		e.Bot.Privmsg(e.Target, "Usage: coin <symbol, e.g. BTC>")
		return
	}
	r, err := rate(apiKey, symbol)
	if err != nil {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("error: %v", err))
		return
	}
	more := fmt.Sprintf("https://www.coinbase.com/price/%v", strings.ToLower(symbol))
	e.Bot.Privmsg(e.Target, fmt.Sprintf("1 %v = %v USD - %v", symbol, r, more))
}

// Register registers the plugin with a bot.
func Register(b bot.Bot, apiKey string) {
	b.Commands().Add("coin", bot.Command{
		Help:    "show rate of a symbol from https://coinlayer.com/",
		Handler: func(e *bot.Event) { coin(e, apiKey) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
