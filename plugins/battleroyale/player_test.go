package battleroyale

import (
	"io/ioutil"
	"reflect"
	"testing"
	"time"
)

func TestParseViewPlayer(t *testing.T) {
	for _, tt := range []struct {
		file string
		want *player
	}{
		{
			file: "testdata/viewplayer_76561198077511634.html",
			want: &player{
				SteamID:             76561198077511634,
				Name:                "WeTarDidSmurf",
				RankText:            "Seasoned",
				GlobalRank:          170,
				TotalPlayTime:       18*time.Hour + 39*time.Minute + 17*time.Second,
				AverageTimeSurvived: 41*time.Minute + 52*time.Second,
				AveragePlacement:    7,
				Wins:                7,
				Losses:              17,
				Kills:               48,
				TotalDistanceMoved:  "229.37 km",
				WinPoints:           1448.88,
				KillPoints:          1356.54,
				TotalPoints:         1720.19,
			},
		},
	} {
		b, err := ioutil.ReadFile(tt.file)
		if err != nil {
			t.Errorf("%s: read: %v", tt.file, err)
			continue
		}
		got, err := parseViewPlayer(string(b))
		if err != nil {
			t.Errorf("%s: parse: %v", tt.file, err)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%s: got %#v\nwant %#v", tt.file, got, tt.want)
		}
	}
}
