// Package timezone implements a plugin to get local time at a location.
// It uses geonames.org for location search and timezone.
package timezone

import (
  "fmt"
  "strings"
  "time"

  "github.com/StalkR/goircbot/bot"
  "github.com/StalkR/goircbot/lib/geonames"
)

func showTime(e *bot.Event, username string) {
  q := strings.TrimSpace(e.Args)
  if len(q) == 0 {
    return
  }
  loc, err := geonames.Search(username, q)
  if err != nil {
    e.Bot.Privmsg(e.Target, fmt.Sprintf("error: %v", err))
    return
  }
  tz, err := geonames.Timezone(username, loc.Latitude, loc.Longitude)
  if err != nil {
    e.Bot.Privmsg(e.Target, fmt.Sprintf("error: %v", err))
    return
  }
  e.Bot.Privmsg(e.Target, fmt.Sprintf("%v, %v: %v - https://www.timeanddate.com/worldclock/@%v,%v",
    loc.Name, loc.Country,
    time.Now().In(tz).Format(time.RFC1123),
    loc.Latitude, loc.Longitude))
}

// Register registers the plugin with a bot.
func Register(b bot.Bot, username string) {
  b.Commands().Add("time", bot.Command{
    Help:    "get time from geonames.org",
    Handler: func(e *bot.Event) { showTime(e, username) },
    Pub:     true,
    Priv:    true,
    Hidden:  false})
}
