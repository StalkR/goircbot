// Package dl implements a plugin to see downloads stats from mldonkey/transmission.
package dl

import (
	"log"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/mldonkey"
	"github.com/StalkR/goircbot/lib/transmission"
)

func mldonkeyStats(url string) (string, error) {
	s, err := mldonkey.Stats(url)
	if err != nil {
		return "", nil
	}
	return s.String(), nil
}

func transmissionStats(url string) (string, error) {
	s, err := transmission.Stats(url)
	if err != nil {
		return "", nil
	}
	return s.String(), nil
}

func dl(b *bot.Bot, e *bot.Event, mldonkeyURL, transmissionURL string) {
	if mldonkeyURL != "" {
		if stats, err := mldonkeyStats(mldonkeyURL); err != nil {
			log.Println("dl: mldonkey:", err)
		} else {
			b.Conn.Privmsg(e.Target, "MLDonkey: "+stats)
		}
	}
	if transmissionURL != "" {
		if stats, err := transmissionStats(transmissionURL); err != nil {
			log.Println("dl: transmission:", err)
		} else {
			b.Conn.Privmsg(e.Target, "Transmission: "+stats)
		}
	}
}

// Register registers the plugin with a bot.
func Register(b *bot.Bot, mldonkeyURL, transmissionURL string) {

	b.AddCommand("dl", bot.Command{
		Help:    "See downloads status",
		Handler: func(b *bot.Bot, e *bot.Event) { dl(b, e, mldonkeyURL, transmissionURL) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
