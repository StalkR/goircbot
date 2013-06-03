// Go IRC Bot example.
package main

import (
	"flag"
	"strings"
	"time"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/plugins/admin"
	"github.com/StalkR/goircbot/plugins/dl"
	"github.com/StalkR/goircbot/plugins/dns"
	"github.com/StalkR/goircbot/plugins/failotron"
	"github.com/StalkR/goircbot/plugins/geo"
	"github.com/StalkR/goircbot/plugins/imdb"
	"github.com/StalkR/goircbot/plugins/mac"
	"github.com/StalkR/goircbot/plugins/ping"
	"github.com/StalkR/goircbot/plugins/scores"
	"github.com/StalkR/goircbot/plugins/search"
	"github.com/StalkR/goircbot/plugins/sed"
	"github.com/StalkR/goircbot/plugins/tail"
	"github.com/StalkR/goircbot/plugins/translate"
	"github.com/StalkR/goircbot/plugins/travisci"
	"github.com/StalkR/goircbot/plugins/up"
	"github.com/StalkR/goircbot/plugins/urban"
	"github.com/StalkR/goircbot/plugins/urltitle"
	"github.com/StalkR/goircbot/plugins/whoami"
)

var (
	host     *string = flag.String("host", "irc.example.com", "Server host[:port]")
	ssl      *bool   = flag.Bool("ssl", true, "Enable SSL")
	nick     *string = flag.String("nick", "goircbot", "Bot nick")
	ident    *string = flag.String("ident", "goircbot", "Bot ident")
	channels *string = flag.String("channels", "", "Channels to join (separated by comma)")

	ignore = []string{"bot"}
)

func main() {
	flag.Parse()
	b := bot.NewBot(*host, *ssl, *nick, *ident, strings.Split(*channels, ","))
	admin.Register(b, []string{"nick!ident@host"})
	dl.Register(b, "", "")
	dns.Register(b)
	failotron.Register(b, ignore)
	geo.Register(b)
	imdb.Register(b)
	mac.Register(b)
	ping.Register(b)
	scores.Register(b, "/tmp/scores")
	search.Register(b, "<key>", "<cx>")
	sed.Register(b)
	tail.Register(b, []string{"/etc/passwd"})
	translate.Register(b, "<key>")
	travisci.Register(b)
	travisci.Watch(b, []string{"StalkR/goircbot"}, 5*time.Minute)
	up.Register(b)
	urban.Register(b)
	urltitle.Register(b, ignore)
	whoami.Register(b)
	b.Run()
}
