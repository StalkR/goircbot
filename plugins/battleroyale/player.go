package battleroyale

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/StalkR/goircbot/lib/transport"
)

type player struct {
	SteamID       uint64
	Name          string
	RankText      string
	GlobalRank    int
	TotalPlayTime time.Duration
	// RoundsPlayed = Wins+Losses
	AverageTimeSurvived time.Duration
	AveragePlacement    int
	Wins                int
	Losses              int
	Kills               int
	TotalDistanceMoved  string
	// WinPercent = Wins/RoundsPlayed
	WinPoints   float64
	KillPoints  float64
	TotalPoints float64
	// AverageKillsPerMatch = Kills/RoundsPlayed
}

func (p player) String() string {
	kd := float64(p.Kills) / float64(p.Losses)
	wl := float64(p.Kills) / float64(p.Wins+p.Losses)
	return fmt.Sprintf("#%d %s, %d wins, %d kills, %d losses, K/D %.2f, W/L %.2f, %s play time %s",
		p.GlobalRank, p.RankText, p.Wins, p.Kills, p.Losses, kd, wl, p.TotalPlayTime, p.URL())
}

func (p player) Short() string {
	return fmt.Sprintf("#%d %s %dW %dK %dL",
		p.GlobalRank, p.RankText, p.Wins, p.Kills, p.Losses)
}

const leaderboardURL = "http://battleroyalegames.com/leaderboard/index.php"

func (p player) URL() string {
	v := url.Values{}
	v.Set("page", "viewplayer")
	v.Set("id", fmt.Sprintf("%d", p.SteamID))
	return fmt.Sprintf("%s?%s", leaderboardURL, v.Encode())
}

type byGlobalRank []player

func (a byGlobalRank) Len() int      { return len(a) }
func (a byGlobalRank) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byGlobalRank) Less(i, j int) bool {
	if a[i].GlobalRank == 0 {
		return false
	}
	if a[j].GlobalRank == 0 {
		return true
	}
	return a[i].GlobalRank < a[j].GlobalRank
}

func viewPlayer(steamID uint64) (*player, error) {
	viewPlayerURL := player{SteamID: steamID}.URL()
	c, err := transport.Client(viewPlayerURL)
	if err != nil {
		return nil, err
	}
	resp, err := c.Get(viewPlayerURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return parseViewPlayer(string(b))
}

var (
	steamIDRE  = regexp.MustCompile(`"http://steamcommunity\.com/profiles/(\d+)"`)
	nameRE     = regexp.MustCompile(`<h4>([^<]+)</h4>`)
	ranktextRE = regexp.MustCompile(`<span class="ranktext"[^>]*>([^<]+)</span>`)
	brRankRE   = regexp.MustCompile(`<span class="br_rank">(\d+)</span>`)
	statboxRE  = regexp.MustCompile(`<div class="statbox2">\s*([^<]+)<br />([^<]+)<br />([^<]+)<br />([^<]+)<br />([^<]+)<br />([^<]+)<br />([^<]+)<br />([^<]+)<br />([^<]+)<br />([^<]+)<br />([^<]+)<br />([^<]+)<br />([^<]+)<br/>\s*</div>`)
)

var errNotFound = errors.New("battleroyale: player not found")

func parseViewPlayer(page string) (*player, error) {
	if strings.Contains(page, "No player profile found") {
		return nil, errNotFound
	}
	p := &player{}
	var err error
	{
		m := steamIDRE.FindStringSubmatch(page)
		if m == nil {
			return nil, fmt.Errorf("battleroyale: steam ID not found")
		}
		p.SteamID, err = strconv.ParseUint(m[1], 10, 64)
		if err != nil {
			return nil, err
		}
	}
	{
		m := nameRE.FindStringSubmatch(page)
		if m == nil {
			return nil, fmt.Errorf("battleroyale: name not found")
		}
		p.Name = m[1]
	}
	{
		m := ranktextRE.FindStringSubmatch(page)
		if m == nil {
			return nil, fmt.Errorf("battleroyale: rank text not found")
		}
		p.RankText = m[1]
	}
	{
		m := brRankRE.FindStringSubmatch(page)
		if m == nil {
			return nil, fmt.Errorf("battleroyale: br rank not found")
		}
		p.GlobalRank, err = strconv.Atoi(m[1])
		if err != nil {
			return nil, err
		}
	}
	m := statboxRE.FindStringSubmatch(page)
	if m == nil {
		return nil, fmt.Errorf("battleroyale: stat box not found")
	}
	// 1 Total Playtime
	// 2 Rounds Played
	// 3 Average Time Survived
	// 4 Average Placement
	// 5 Wins
	// 6 Losses
	// 7 Kills
	// 8 Total Distance Moved
	// 9 Win Percent
	// 10 Win Points
	// 11 Kill Points
	// 12 Total Points
	// 13 Average Kills / Match
	p.TotalPlayTime, err = time.ParseDuration(strings.Replace(m[1], " ", "", -1))
	if err != nil {
		return nil, err
	}
	p.AverageTimeSurvived, err = time.ParseDuration(strings.Replace(m[3], " ", "", -1))
	if err != nil {
		return nil, err
	}
	p.AveragePlacement, err = strconv.Atoi(m[4])
	if err != nil {
		return nil, err
	}
	p.Wins, err = strconv.Atoi(m[5])
	if err != nil {
		return nil, err
	}
	p.Losses, err = strconv.Atoi(m[6])
	if err != nil {
		return nil, err
	}
	p.Kills, err = strconv.Atoi(m[7])
	if err != nil {
		return nil, err
	}
	p.TotalDistanceMoved = m[8]
	p.WinPoints, err = strconv.ParseFloat(m[10], 64)
	if err != nil {
		return nil, err
	}
	p.KillPoints, err = strconv.ParseFloat(m[11], 64)
	if err != nil {
		return nil, err
	}
	p.TotalPoints, err = strconv.ParseFloat(m[12], 64)
	if err != nil {
		return nil, err
	}
	return p, nil
}
