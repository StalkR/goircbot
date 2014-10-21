package battleroyale

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/StalkR/goircbot/lib/transport"
)

const scoreURL = "http://battleroyale.gamingdeluxe.co.uk:8080/api/Players"

type playerInfo struct {
	ID                  int
	UID                 string
	ImageURL            string `json:"ImgUrl"`
	Name                string
	Wins, Kills, Loss   int
	KillDeathRatio      float64
	WinRate             float64
	AverageKillDistance float64
	MaxKillDistance     float64
	Points              int `json:"PTS"`
}

func (p playerInfo) String() string {
	return fmt.Sprintf("%d wins, %d kills, %d loss, %d points, K/D %.2f, W/L %.2f, max kill distance %.2f",
		p.Wins, p.Kills, p.Loss, p.Points, p.KillDeathRatio, p.WinRate, p.MaxKillDistance)
}

func (p playerInfo) Short() string {
	return fmt.Sprintf("%d W, %d L, %d K, %d pts", p.Wins, p.Loss, p.Kills, p.Points)
}

type byPoints []*playerInfo

func (b byPoints) Len() int           { return len(b) }
func (b byPoints) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byPoints) Less(i, j int) bool { return b[i].Points < b[j].Points }

func scoreByName(name string) (*playerInfo, error) {
	return getPlayerInfo(url.Values{"Count": {"1"}, "Name": {name}})
}

func scoreByUID(uid string) (*playerInfo, error) {
	p, err := getPlayerInfo(url.Values{"Count": {"1"}, "Id": {uid}})
	if err != nil {
		return nil, err
	}
	if p.UID != uid {
		return nil, fmt.Errorf("got stats for UID %v; want %v", p.UID, uid)
	}
	return p, nil
}

var errNotFound = errors.New("not found")

func getPlayerInfo(u url.Values) (*playerInfo, error) {
	c, err := transport.Client(scoreURL)
	if err != nil {
		return nil, err
	}
	resp, err := c.Get(fmt.Sprintf("%s?%s", scoreURL, u.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var p []*playerInfo
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return nil, err
	}
	if len(p) != 1 {
		return nil, errNotFound
	}
	return p[0], nil
}
