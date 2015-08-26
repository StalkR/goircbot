package hots

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"regexp"
	"strconv"

	"github.com/StalkR/goircbot/lib/transport"
)

type Score struct {
	Medal string
	Score int
	MMR   int
}

type Stats struct {
	PlayerID   int
	TeamLeague Score
	HeroLeague Score
	QuickMatch Score
}

func (s *Stats) String() string {
	tl := fmt.Sprintf("Team League: %s %d (MMR %d)", s.TeamLeague.Medal, s.TeamLeague.Score, s.TeamLeague.MMR)
	if s.TeamLeague.Medal == "" {
		tl = "Team League: n/a"
	}
	hl := fmt.Sprintf("Hero League: %s %d (MMR %d)", s.HeroLeague.Medal, s.HeroLeague.Score, s.HeroLeague.MMR)
	if s.HeroLeague.Medal == "" {
		hl = "Hero League: n/a"
	}
	qm := fmt.Sprintf("Quick Match: %s %d (MMR %d)", s.QuickMatch.Medal, s.QuickMatch.Score, s.QuickMatch.MMR)
	if s.QuickMatch.Medal == "" {
		qm = "Quick Match: n/a"
	}
	return fmt.Sprintf("%s, %s, %s - %s?PlayerID=%d", tl, hl, qm, statsURL, s.PlayerID)
}

func NewStats(id int) (*Stats, error) {
	b, err := get(id)
	if err != nil {
		return nil, err
	}
	return parse(string(b))
}

const statsURL = "https://www.hotslogs.com/Player/Profile"

func get(id int) ([]byte, error) {
	c, err := transport.Client(statsURL)
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Add("PlayerID", fmt.Sprintf("%d", id))
	resp, err := c.Get(fmt.Sprintf("%s?%s", statsURL, v.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

var (
	playerIDRE   = regexp.MustCompile(`<link href="https://www\.hotslogs\.com/Player/Profile\?PlayerID=(\d+)"`)
	teamLeagueRE = regexp.MustCompile(`<th>Team League</th><td><img[^>]*>[^<]*<span>(\w+) (\d+) \(Current MMR: (\d+)\)</span>`)
	heroLeagueRE = regexp.MustCompile(`<th>Hero League</th><td><img[^>]*>[^<]*<span>(\w+) (\d+) \(Current MMR: (\d+)\)</span>`)
	quickMatchRE = regexp.MustCompile(`<th>Quick Match</th><td><img[^>]*>[^<]*<span>(\w+) (\d+) \(Current MMR: (\d+)\)</span>`)
)

func parse(page string) (*Stats, error) {
	var s Stats
	if m := playerIDRE.FindStringSubmatch(page); m != nil {
		playerID, err := strconv.Atoi(m[1])
		if err != nil {
			return nil, err
		}
		s.PlayerID = playerID
	} else {
		return nil, errors.New("hots: could not find player ID")
	}
	if m := teamLeagueRE.FindStringSubmatch(page); m != nil {
		s.TeamLeague.Medal = m[1]
		score, err := strconv.Atoi(m[2])
		if err != nil {
			return nil, err
		}
		s.TeamLeague.Score = score
		mmr, err := strconv.Atoi(m[3])
		if err != nil {
			return nil, err
		}
		s.TeamLeague.MMR = mmr
	}
	if m := heroLeagueRE.FindStringSubmatch(page); m != nil {
		s.HeroLeague.Medal = m[1]
		score, err := strconv.Atoi(m[2])
		if err != nil {
			return nil, err
		}
		s.HeroLeague.Score = score
		mmr, err := strconv.Atoi(m[3])
		if err != nil {
			return nil, err
		}
		s.HeroLeague.MMR = mmr
	}
	if m := quickMatchRE.FindStringSubmatch(page); m != nil {
		s.QuickMatch.Medal = m[1]
		score, err := strconv.Atoi(m[2])
		if err != nil {
			return nil, err
		}
		s.QuickMatch.Score = score
		mmr, err := strconv.Atoi(m[3])
		if err != nil {
			return nil, err
		}
		s.QuickMatch.MMR = mmr
	}
	return &s, nil
}
