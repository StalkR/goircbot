// Package untappd implements a plugin to get info from Untappd users.
package untappd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/duration"
	"github.com/StalkR/goircbot/lib/transport"
)

var (
	indexRE   = regexp.MustCompile(`<title>Untappd`)
	nameRE    = regexp.MustCompile(`<title>(.*) on Untappd</title>`)
	totalRE   = regexp.MustCompile(`(?s)<p><strong>Total</strong></p>\s*<span>(\d+)`)
	detailsRE = regexp.MustCompile(`(?s)<div class="details">(.*)`)
	beerRE    = regexp.MustCompile(`(?s)<a href=./b/[^>]+>([^<]+)`)
	whenRE    = regexp.MustCompile(`(?s)<li class="time[^>]+>([^<]+)</li>`)
)

// Info represents an Untappd user info with last activity and total.
type Info struct {
	Name  string
	Beer  string
	When  time.Time
	Total int
}

// String formats an Info to fit on a single IRC line.
func (i Info) String() string {
	if i.Total == 0 {
		return fmt.Sprintf("%v doesn't drink, booo!", i.Name)
	}
	return fmt.Sprintf("%v is drinking %v (%v ago) - Total beers: %v",
		i.Name, i.Beer, duration.Format(time.Since(i.When)), i.Total)
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

// parsePage parses an Untappd user page and return its info.
func parsePage(page []byte) (i Info, err error) {
	if indexRE.Match(page) {
		return i, errors.New("no such user")
	}

	match := nameRE.FindSubmatch(page)
	if match == nil {
		return i, errors.New("untappd: could not parse name")
	}
	i.Name = strings.TrimSpace(string(match[1]))

	match = totalRE.FindSubmatch(page)
	if match == nil {
		return i, errors.New("untappd: could not parse total")
	}
	i.Total, err = strconv.Atoi(strings.TrimSpace(string(match[1])))
	if err != nil {
		return i, fmt.Errorf("untappd: error parsing total: %v", err)
	}
	if i.Total == 0 {
		return i, nil
	}

	match = detailsRE.FindSubmatch(page)
	if match == nil {
		return i, errors.New("untappd: could not parse details")
	}
	details := match[1]

	match = beerRE.FindSubmatch(details)
	if match == nil {
		return i, errors.New("untappd: could not parse beer")
	}
	i.Beer = strings.TrimSpace(string(match[1]))

	match = whenRE.FindSubmatch(details)
	if match == nil {
		return i, errors.New("untappd: could not parse when")
	}
	when, err := time.Parse(time.RFC1123Z, strings.TrimSpace(string(match[1])))
	if err != nil {
		return i, fmt.Errorf("untappd: error parsing time: %v", err)
	}
	i.When = when

	return
}

// UserInfo obtains info for a user by querying its web page and parsing it.
func UserInfo(user string) (Info, error) {
	page, err := userPage(user)
	if err != nil {
		return Info{}, fmt.Errorf("untappd: error getting page: %v", err)
	}
	return parsePage(page)
}

// untappd handles an IRC command to print an Untappd user info.
func untappd(e *bot.Event) {
	user := strings.TrimSpace(e.Args)
	if len(user) == 0 {
		return
	}
	info, err := UserInfo(user)
	if err != nil {
		e.Bot.Privmsg(e.Target, err.Error())
		return
	}
	e.Bot.Privmsg(e.Target, info.String())
}

// Register registers the plugin with a bot.
func Register(b bot.Bot) {
	b.Commands().Add("untappd", bot.Command{
		Help:    "get stats and activity from an Untappd user",
		Handler: untappd,
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
