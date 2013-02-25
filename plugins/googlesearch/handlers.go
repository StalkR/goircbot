// Package googlesearch implements a plugin to search on Google Custom Search.
package googlesearch

import (
	bot "github.com/StalkR/goircbot"
	"github.com/StalkR/misc/google/customsearch"
	"log"
	"strings"
)

func Search(b *bot.Bot, e *bot.Event, key, cx string) {
	term := strings.TrimSpace(e.Args)
	if len(term) == 0 {
		return
	}
	r, err := customsearch.Search(term, key, cx)
	if err != nil {
		log.Println("googlesearch:", err)
		return
	}
	if len(r.Items) == 0 {
		b.Conn.Privmsg(e.Target, "No result.")
		return
	}
	b.Conn.Privmsg(e.Target, r.Items[0].String())
}

// Register registers the plugin with a bot.
func Register(b *bot.Bot, key, cx string) {
	b.AddCommand("search", bot.Command{
		Help:    "search the Web with Google",
		Handler: func(b *bot.Bot, e *bot.Event) { Search(b, e, key, cx) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
