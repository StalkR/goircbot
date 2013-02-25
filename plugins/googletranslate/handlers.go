// Package googletranslate implements a plugin to translate with Google Translate.
package googletranslate

import (
	"fmt"
	bot "github.com/StalkR/goircbot"
	"github.com/StalkR/misc/google/translate"
	"log"
	"strings"
)

func Supported(target, key string) ([]string, error) {
	languages, err := translate.Languages(target, key)
	if err != nil {
		return nil, err
	}
	var langs []string
	for _, lang := range languages {
		langs = append(langs, lang.Language)
	}
	return langs, nil
}

func Translate(b *bot.Bot, e *bot.Event, key string) {
	line := strings.TrimSpace(e.Args)
	args := strings.SplitN(line, " ", 3)
	var reply string
	switch {
	case len(line) == 0:
		langs, err := Supported("", key)
		if err != nil {
			log.Println("googletranslate:", err)
			return
		}
		reply = fmt.Sprintf("Supported languages: %s", strings.Join(langs, ", "))

	case len(args) == 1:
		langs, err := Supported(line, key)
		if err != nil {
			log.Println("googletranslate:", err)
			return
		}
		reply = fmt.Sprintf("Supported languages for %s: %s", line, strings.Join(langs, ", "))

	case len(args) == 2:
		return

	default:
		source, target, text := args[0], args[1], args[2]
		if source == "-" {
			source = ""
		}
		t, err := translate.Translate(source, target, text, key)
		if err != nil {
			log.Println("googletranslate:", err)
			return
		}
		if source == "" {
			source = t.DetectedSourceLanguage
		}
		reply = fmt.Sprintf("%s->%s: %s\n", source, target, t.TranslatedText)
	}
	b.Conn.Privmsg(e.Target, reply)
}

// Register registers the plugin with a bot.
func Register(b *bot.Bot, key string) {
	b.AddCommand("tr", bot.Command{
		Help:    "translate <from|-> <to> <text> using Google Translate",
		Handler: func(b *bot.Bot, e *bot.Event) { Translate(b, e, key) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
