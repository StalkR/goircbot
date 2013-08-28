// Package untappd implements a plugin to get info from Untappd users.
package untappd

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/transport"
)

var (
	detailsRE = regexp.MustCompile(`(?s)<div class="details">(.*)`)
	userRE    = regexp.MustCompile(`(?s)<a href=./user[^>]+>([^<]+)`)
	beerRE    = regexp.MustCompile(`(?s)<a href=./b/[^>]+>([^<]+)`)
	whenRE    = regexp.MustCompile(`(?s)<li class="time[^>]+>([^<]+)</li>`)
	totalRE   = regexp.MustCompile(`(?s)<p><strong>Total</strong></p>\s*<span>(\d+)`)
)

// Info represents an Untappd user info with last activity and total.
type Info struct {
	User  string
	Beer  string
	When  time.Time
	Total int
}

// String formats an Info to fit on a single IRC line.
func (i Info) String() string {
	d := time.Since(i.When)
	d = d / time.Second * time.Second
	return fmt.Sprintf("%v is drinking %v (%v ago) - Total beers: %v",
		i.User, i.Beer, d, i.Total)
}

// userPage fetches an Untappd user page.
func userPage(user string) ([]byte, error) {
	url := fmt.Sprintf("https://untappd.com/user/%s", user)
	client, err := transport.Client(url)
	if err != nil {
		return nil, err
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	page, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return page, nil
}

// UserInfo fetches then parses an Untappd user page to return its info.
func UserInfo(user string) (i Info, err error) {
	page, err := userPage(user)
	if err != nil {
		return
	}
	match := detailsRE.FindSubmatch(page)
	if match == nil {
		return
	}
	details := match[1]

	match = userRE.FindSubmatch(details)
	if match == nil {
		return
	}
	i.User = strings.TrimSpace(string(match[1]))

	match = beerRE.FindSubmatch(details)
	if match == nil {
		return
	}
	i.Beer = strings.TrimSpace(string(match[1]))

	match = whenRE.FindSubmatch(details)
	if match == nil {
		return
	}
	when, err := time.Parse(time.RFC1123Z, strings.TrimSpace(string(match[1])))
	if err != nil {
		return
	}
	i.When = when

	match = totalRE.FindSubmatch(page)
	if match == nil {
		return
	}
	i.Total, err = strconv.Atoi(strings.TrimSpace(string(match[1])))
	if err != nil {
		return
	}

	return
}

// untappd handles an IRC command to print an Untappd user info.
func untappd(b *bot.Bot, e *bot.Event) {
	user := strings.TrimSpace(e.Args)
	if len(user) == 0 {
		return
	}
	info, err := UserInfo(user)
	if err != nil {
		b.Conn.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	b.Conn.Privmsg(e.Target, info.String())
}

// Register registers the plugin with a bot.
func Register(b *bot.Bot) {
	b.AddCommand("untappd", bot.Command{
		Help:    "get stats and activity from an Untappd user",
		Handler: untappd,
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
