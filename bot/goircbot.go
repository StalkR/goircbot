// Go IRC Bot example.
package main

import (
	"flag"
	bot "goircbot"
	"goircbot/plugins/admin"
	"goircbot/plugins/failotron"
	"strings"
)

var host *string = flag.String("host", "irc.example.com", "Server host[:port]")
var ssl *bool = flag.Bool("ssl", true, "Enable SSL")
var nick *string = flag.String("nick", "goircbot", "Bot nick")
var ident *string = flag.String("ident", "goircbot", "Bot ident")
var channels *string = flag.String("channels", "", "Channels to join (separated by comma)")

func main() {
	flag.Parse()
	b := bot.NewBot(*host, *ssl, *nick, *ident, strings.Split(*channels, ","))
	admin.Register(b, []string{"nick!ident@host"})
	failotron.Register(b)
	b.Run()
}
