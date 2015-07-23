// Package git implements a plugin to view and watch git commits.
package git

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/transport"
)

type commitMsg struct {
	Msg    string
	Age    string
	Author string
}

func (c *commitMsg) String() string {
	return fmt.Sprintf("%s (%s ago by %s)", c.Msg, c.Age, c.Author)
}

var (
	logRE = regexp.MustCompile(`<span class='age-[^']+'>([^<]+)</span></td>` +
		`<td><a href='([^']+)'>([^<]+)</a>.*?</span></td>.*?>([^<]+)</td>`)
	hrefRE = regexp.MustCompile(`/commit/\?id=[0-9a-f]{40}$`)
)

func lastLog(url string) (*commitMsg, error) {
	c, err := transport.Client(url)
	if err != nil {
		return nil, err
	}
	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	for _, match := range logRE.FindAllStringSubmatch(string(b), -1) {
		age, href, msg, author := match[1], match[2], match[3], match[4]
		if hrefRE.MatchString(href) {
			return &commitMsg{Msg: msg, Age: age, Author: author}, nil
		}
	}
	return nil, fmt.Errorf("git: last log not found")
}

func handleGit(e *bot.Event, repos map[string]string) {
	repo, ok := repos[strings.TrimSpace(e.Args)]
	if !ok {
		e.Bot.Privmsg(e.Target, "not found")
		return
	}
	commit, err := lastLog(repo)
	if err != nil {
		log.Printf("git: last log %s: %v", repo, err)
		e.Bot.Privmsg(e.Target, "error")
		return
	}
	e.Bot.Privmsg(e.Target, commit.String())
}

// Register registers the plugin with a bot.
// It shows the last last commit of a repo identified by its short name.
// Repos is a map of short name to URL to a cgit log page such as
// https://git.kernel.org/cgit/linux/kernel/git/torvalds/linux.git/log/
func Register(b bot.Bot, repos map[string]string) {
	b.Commands().Add("git", bot.Command{
		Help:    "get last commit of a repo",
		Handler: func(e *bot.Event) { handleGit(e, repos) },
		Pub:     true,
		Priv:    false,
		Hidden:  false})
}
