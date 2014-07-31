package battleroyale

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"

	"github.com/StalkR/goircbot/lib/transport"
)

const scoreboardURL = "http://battleroyale.unknownservers.net/leaderboard/"

func get() (map[string]score, error) {
	c, err := transport.Client(scoreboardURL)
	if err != nil {
		return nil, err
	}
	resp, err := c.Get(scoreboardURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return parse(string(b))
}

var scoreRE = regexp.MustCompile(`<tr>\s*` +
	`<td><a[^>]*>(\d+)</a></td>\s*` +
	`<td><center>(\d+)</center></td>\s*` +
	`<td><center>(\d+)</center></td>\s*` +
	`<td><center>(\d+)</center></td>\s*` +
	`</tr>`)

func parse(page string) (map[string]score, error) {
	m := make(map[string]score)
	s := scoreRE.FindAllStringSubmatch(page, -1)
	if s == nil {
		return nil, fmt.Errorf("battleroyale: no match")
	}
	for _, e := range s {
		wins, _ := strconv.Atoi(e[2])
		losses, _ := strconv.Atoi(e[3])
		kills, _ := strconv.Atoi(e[4])
		m[e[1]] = score{wins: wins, losses: losses, kills: kills}
	}
	return m, nil
}
