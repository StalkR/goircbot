// Go IRC Bot example.
package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/size"
	"github.com/StalkR/goircbot/plugins/admin"
	"github.com/StalkR/goircbot/plugins/asm"
	"github.com/StalkR/goircbot/plugins/battlefield"
	"github.com/StalkR/goircbot/plugins/cdecl"
	"github.com/StalkR/goircbot/plugins/darkstat"
	"github.com/StalkR/goircbot/plugins/df"
	"github.com/StalkR/goircbot/plugins/dl"
	"github.com/StalkR/goircbot/plugins/dns"
	"github.com/StalkR/goircbot/plugins/errors"
	"github.com/StalkR/goircbot/plugins/failotron"
	"github.com/StalkR/goircbot/plugins/geo"
	"github.com/StalkR/goircbot/plugins/git"
	"github.com/StalkR/goircbot/plugins/golang"
	"github.com/StalkR/goircbot/plugins/hots"
	"github.com/StalkR/goircbot/plugins/idle"
	"github.com/StalkR/goircbot/plugins/imdb"
	"github.com/StalkR/goircbot/plugins/invite"
	"github.com/StalkR/goircbot/plugins/mac"
	"github.com/StalkR/goircbot/plugins/metal"
	"github.com/StalkR/goircbot/plugins/old"
	"github.com/StalkR/goircbot/plugins/ping"
	"github.com/StalkR/goircbot/plugins/quotes"
	"github.com/StalkR/goircbot/plugins/renick"
	"github.com/StalkR/goircbot/plugins/scores"
	"github.com/StalkR/goircbot/plugins/search"
	"github.com/StalkR/goircbot/plugins/sed"
	"github.com/StalkR/goircbot/plugins/tail"
	timep "github.com/StalkR/goircbot/plugins/time"
	"github.com/StalkR/goircbot/plugins/tor"
	"github.com/StalkR/goircbot/plugins/translate"
	"github.com/StalkR/goircbot/plugins/travisci"
	"github.com/StalkR/goircbot/plugins/up"
	"github.com/StalkR/goircbot/plugins/urban"
	"github.com/StalkR/goircbot/plugins/urltitle"
	"github.com/StalkR/goircbot/plugins/whoami"
	"github.com/StalkR/goircbot/plugins/wunderground"
	"github.com/fluffle/goirc/logging/glog"
)

var (
	host     = flag.String("host", "irc.example.com", "Server host[:port]")
	ssl      = flag.Bool("ssl", true, "Enable SSL")
	nick     = flag.String("nick", "goircbot", "Bot nick")
	ident    = flag.String("ident", "goircbot", "Bot ident")
	channels = flag.String("channels", "", "Channels to join (separated by comma)")

	ignore = []string{"bot"}
)

func main() {
	flag.Parse()
	glog.Init()
	b := bot.NewBot(*host, *ssl, *nick, *ident, strings.Split(*channels, ","))
	admin.Register(b, []string{"nick!ident@host"})
	asm.Register(b)
	if err := battlefield.Register(b, "account@ea.com", "password",
		map[string]uint64{
			"playerName": 123456789, // persona ID
		}); err != nil {
		log.Fatal(err)
	}
	cdecl.Register(b)
	darkstat.Register(b, map[string]string{
		"public":  "http://darkstat.public.com",
		"private": "https://user:pass@darkstat.private.com",
	})
	df.Register(b, df.NewAlarm(`/`, 10*size.GB))
	dl.Register(b, "", "")
	dns.Register(b)
	errors.Register(b)
	failotron.Register(b, ignore)
	hots.Register(b, map[string]int{
		"playerName": 1234, // Player ID
	})
	geo.Register(b)
	git.Register(b, map[string]string{
		"linux": "https://git.kernel.org/cgit/linux/kernel/git/torvalds/linux.git/log/",
	})
	golang.Register(b)
	idle.Register(b, ignore)
	imdb.Register(b)
	invite.Register(b)
	mac.Register(b)
	metal.Register(b)
	old.Register(b, "/tmp/old", ignore)
	ping.Register(b)
	quotes.Register(b, "/tmp/quotes")
	renick.Register(b, *nick)
	scores.Register(b, "/tmp/scores")
	search.Register(b, "<key>", "<cx>")
	sed.Register(b)
	tail.Register(b, []string{"/etc/passwd"})
	timep.Register(b, "<key>")
	tor.Register(b, "127.0.0.1:9051", "secret")
	translate.Register(b, "<key>")
	travisci.Register(b)
	travisci.Watch(b, []string{"StalkR/goircbot"}, 5*time.Minute)
	up.Register(b)
	urban.Register(b)
	urltitle.Register(b, ignore)
	whoami.Register(b)
	wunderground.Register(b, "<key>")
	b.Run()
}
