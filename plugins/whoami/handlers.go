// Package whoami implements a plugin to report user and host the bot runs as.
package whoami

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/StalkR/goircbot/bot"
)

func whoami(e *bot.Event) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("whoami: hostname error", err)
		hostname = "<error>"
	}
	var username string
	if current, err := user.Current(); err != nil {
		log.Println("whoami: user error", err)
		username = "<error>"
	} else {
		username = current.Username
	}
	reply := fmt.Sprintf("is running as %s@%s", username, hostname)
	e.Bot.Action(e.Target, reply)
}

// Register registers the plugin with a bot.
func Register(b bot.Bot) {
	b.Commands().Add("whoami", bot.Command{
		Help:    "report user@host of the bot",
		Handler: whoami,
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
