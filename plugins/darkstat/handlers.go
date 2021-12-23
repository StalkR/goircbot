// Package darkstat implements a plugin to see stats from darkstat.
package darkstat

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/darkstat"
	"github.com/StalkR/goircbot/lib/size"
)

// bw handles command to display either all darkstat or a single one by name.
func bw(e *bot.Event, c map[string]*darkstat.Conn) {
	name := strings.TrimSpace(e.Args)
	var stats string
	var err error
	if ds, ok := c[name]; ok {
		stats, err = single(name, ds)
	} else {
		stats, err = multiple(c)
	}
	if err != nil {
		log.Print("darkstat: error: ", err)
		e.Bot.Privmsg(e.Target, "darkstat: failed (see logs)")
		return
	}
	e.Bot.Privmsg(e.Target, stats)
}

// single formats bandwidth stats for a single darkstat.
func single(name string, ds *darkstat.Conn) (string, error) {
	i, o, err := ds.Bandwidth()
	if err != nil {
		return "", err
	}
	in, out := size.Byte(i).String(), size.Byte(o).String()
	return fmt.Sprintf("%v: %v/s in, %v/s out", name, in, out), nil
}

// multiple formats bandwidth stats for multiple darkstat.
func multiple(c map[string]*darkstat.Conn) (string, error) {
	var stats []string
	for name, ds := range c {
		stat, err := single(name, ds)
		if err != nil {
			log.Printf("darkstat: error %v: %v", name, err)
			continue
		}
		stats = append(stats, stat)
	}
	if len(stats) == 0 {
		return "", errors.New("darkstat: all failed")
	}
	sort.Strings(stats)
	return strings.Join(stats, " - "), nil
}

// Register registers the plugin with a bot.
func Register(b bot.Bot, URLs map[string]string) {
	c := make(map[string]*darkstat.Conn)
	for name, url := range URLs {
		ds, err := darkstat.New(url)
		if err != nil {
			panic(err)
		}
		c[name] = ds
	}

	b.Commands().Add("bw", bot.Command{
		Help:    "See current bandwidth stats from darkstat",
		Handler: func(e *bot.Event) { bw(e, c) },
		Pub:     true,
		Priv:    false,
		Hidden:  false})
}
