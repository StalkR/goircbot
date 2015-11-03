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
				RankText:            "Champion",
				GlobalRank:          1,
				TotalPlayTime:       33*time.Hour + 4*time.Minute + 20*time.Second,
				AverageTimeSurvived: 40*time.Minute + 29*time.Second,
				AveragePlacement:    5,
				Wins:                21,
				Losses:              28,
				Kills:               100,
				TotalDistanceMoved:  "478.78 km",
				WinPoints:           1538.35,
				KillPoints:          1404.2,
				TotalPoints:         1819.19,
				FavouriteMatchType:  "100 % Regular",
			},
		},
		{
			file: "testdata/viewplayer_76561198134388393.html",
			want: &player{
				SteamID:             76561198134388393,
				Name:                "NeKIT",
				RankText:            "Beginner",
				GlobalRank:          859,
				TotalPlayTime:       9*time.Hour + 37*time.Minute + 36*time.Second,
				AverageTimeSurvived: 14*time.Minute + 5*time.Second,
				AveragePlacement:    28,
				Wins:                1,
				Losses:              40,
				Kills:               18,
				TotalDistanceMoved:  "125.81 km",
				WinPoints:           1112.33,
				KillPoints:          1062.45,
				TotalPoints:         1324.82,
				FavouriteMatchType:  "85 % Regular",
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
