// Package translate implements a plugin to translate with Google Translate.
package translate

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/google/translate"
)

func supported(target, key string) ([]string, error) {
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

func compactSpaces(s string) string {
	r, err := regexp.Compile(`\s\s+`)
	if err != nil {
		return s
	}
	return string(r.ReplaceAll([]byte(s), []byte(" ")))
}

func tr(e *bot.Event, key string) {
	line := strings.TrimSpace(compactSpaces(e.Args))
	args := strings.SplitN(line, " ", 3)
	var reply string
	switch {
	case len(line) == 0:
		langs, err := supported("", key)
		if err != nil {
			log.Println("translate:", err)
			return
		}
		reply = fmt.Sprintf("Supported languages: %s", strings.Join(langs, ", "))

	case len(args) == 1:
		langs, err := supported(line, key)
		if err != nil {
			log.Println("translate:", err)
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
			log.Println("translate:", err)
			return
		}
		if source == "" {
			source = t.DetectedSourceLanguage
		}
		reply = fmt.Sprintf("%s->%s: %s\n", source, target, t.TranslatedText)
	}
	if len(reply) > 300 {
		reply = reply[:300]
	}
	e.Bot.Privmsg(e.Target, reply)
}

// Register registers the plugin with a bot.
func Register(b bot.Bot, key string) {
	b.Commands().Add("tr", bot.Command{
		Help:    "translate <from|-> <to> <text> using Google Translate",
		Handler: func(e *bot.Event) { tr(e, key) },
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
