// Go IRC Bot example.
package main

import (
	"flag"
	bot "github.com/StalkR/goircbot"
	"github.com/StalkR/goircbot/plugins/admin"
	"github.com/StalkR/goircbot/plugins/dns"
	"github.com/StalkR/goircbot/plugins/failotron"
	"github.com/StalkR/goircbot/plugins/imdb"
	"github.com/StalkR/goircbot/plugins/scores"
	"github.com/StalkR/goircbot/plugins/up"
	"github.com/StalkR/goircbot/plugins/urbandictionary"
	"github.com/StalkR/goircbot/plugins/whoami"
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
	dns.Register(b)
	failotron.Register(b, []string{"bot"})
	imdb.Register(b)
	scores.Register(b, "/tmp/scores")
	up.Register(b)
	urbandictionary.Register(b)
	whoami.Register(b)
	b.Run()
}
