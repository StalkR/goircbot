// Package stock implements a plugin to get stock quotes.
package stock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/StalkR/goircbot/bot"
)

type quote struct {
	Status  string `json:"status"` // "error"
	Code    int    `json:"code"`   // 400
	Message string `json:"message"`

	Symbol         string `json:"symbol"`
	Name           string `json:"name"`
	Exchange       string `json:"exchange"`
	MICCode        string `json:"mic_code"`
	Currency       string `json:"currency"`
	Datetime       string `json:"datetime"`
	Timestamp      uint64 `json:"timestamp"`
	OpenS          string `json:"open"`
	Open           float64
	HighS          string `json:"high"`
	High           float64
	LowS           string `json:"low"`
	Low            float64
	CloseS         string `json:"close"`
	Close          float64
	VolumeS        string `json:"volume"`
	Volume         uint64
	PreviousCloseS string `json:"previous_close"`
	PreviousClose  float64
	ChangeS        string `json:"change"`
	Change         float64
	PercentChangeS string `json:"percent_change"`
	PercentChange  float64
	AverageVolumeS string `json:"average_volume"`
	AverageVolume  uint64
	IsMarketOpen   bool `json:"is_market_open"`
}

func (q *quote) parse() {
	if v, err := strconv.ParseFloat(q.OpenS, 64); err == nil {
		q.Open = v
	}
	if v, err := strconv.ParseFloat(q.HighS, 64); err == nil {
		q.High = v
	}
	if v, err := strconv.ParseFloat(q.LowS, 64); err == nil {
		q.Low = v
	}
	if v, err := strconv.ParseFloat(q.CloseS, 64); err == nil {
		q.Close = v
	}
	if v, err := strconv.ParseUint(q.VolumeS, 10, 64); err == nil {
		q.Volume = v
	}
	if v, err := strconv.ParseFloat(q.PreviousCloseS, 64); err == nil {
		q.PreviousClose = v
	}
	if v, err := strconv.ParseFloat(q.ChangeS, 64); err == nil {
		q.Change = v
	}
	if v, err := strconv.ParseFloat(q.PercentChangeS, 64); err == nil {
		q.PercentChange = v
	}
	if v, err := strconv.ParseUint(q.AverageVolumeS, 10, 64); err == nil {
		q.AverageVolume = v
	}
}

func (q *quote) String() string {
	plus := ""
	if q.Change > 0 {
		plus = "+"
	}
	return fmt.Sprintf("%s (%s): %.2f %s (%s%.2f, %s%.2f%%), %s volume - https://finance.yahoo.com/quote/%s",
		q.Symbol, q.Name, q.Close, q.Currency, plus, q.Change, plus, q.PercentChange, humanize(q.Volume), q.Symbol)
}

func humanize(i uint64) string {
	format := func(i, unit uint64) string {
		f := "%.0f"
		if i < 10*unit {
			f = "%.2f"
		} else if i < 100*unit {
			f = "%.1f"
		}
		return fmt.Sprintf(f, float64(i)/float64(unit))
	}
	switch {
	case i >= 1e9: // Billion
		return fmt.Sprintf("%vB", format(i, 1e9))
	case i >= 1e6: // Million
		return fmt.Sprintf("%vM", format(i, 1e6))
	case i >= 1e3: // Thousand
		return fmt.Sprintf("%vK", format(i, 1e3))
	}
	return fmt.Sprintf("%v", i)
}

const apiURL = "https://api.twelvedata.com"

func stock(apiKey, symbol string) (*quote, error) {
	v := url.Values{}
	v.Set("symbol", symbol)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/quote?%s", apiURL, v.Encode()), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("apikey %v", apiKey))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%v", resp.Status)
	}
	var r quote
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	if r.Status == "error" {
		return nil, fmt.Errorf("error: %v", r.Message)
	}
	r.parse()
	return &r, nil
}

func handleStock(e *bot.Event, apiKey string) {
	symbol := strings.TrimSpace(e.Args)
	if len(symbol) == 0 {
		return
	}
	q, err := stock(apiKey, symbol)
	if err != nil {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	e.Bot.Privmsg(e.Target, q.String())
}

// Register registers the plugin with a bot.
func Register(b bot.Bot, apiKey string) {
	b.Commands().Add("stock", bot.Command{
		Help:    "get trading stock information like price",
		Handler: func(e *bot.Event) { handleStock(e, apiKey) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
