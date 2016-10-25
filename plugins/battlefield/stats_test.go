package battlefield

import (
	"io/ioutil"
	"reflect"
	"testing"
	"time"
)

func TestParseStats(t *testing.T) {
	for _, tt := range []struct {
		file    string
		name    string
		id      uint64
		want    *Stats
		wantErr error
	}{
		{
			file:    "testdata/0.json",
			wantErr: errNotFound,
		},
		{
			file: "testdata/384000070.json",
			name: "r3skp",
			id:   384000070,
			want: &Stats{
				ID:   384000070,
				Name: "r3skp",
				BF1: BFStats{
					ID:          384000070,
					Game:        "bf1",
					TimePlayed:  time.Second * 71596,
					Kills:       1012,
					Deaths:      690,
					Wins:        34,
					Losses:      26,
					Rank:        34,
					ScorePerMin: 927.59,
				},
				BF4: BFStats{
					ID:          384000070,
					Game:        "bf4",
					TimePlayed:  time.Second * 2183110,
					Kills:       46278,
					Deaths:      23539,
					Wins:        1712,
					Losses:      1055,
					Rank:        140,
					ScorePerMin: 1083.86,
				},
			},
		},
	} {
		js, err := ioutil.ReadFile(tt.file)
		if err != nil {
			t.Errorf("%s: read: %v", tt.file, err)
			continue
		}
		got, err := parseStats(tt.id, tt.name, js)
		if err != tt.wantErr {
			t.Errorf("%s error: got %#v\nwant %#v", tt.file, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%s: got %#v\nwant %#v", tt.file, got, tt.want)
		}
	}
}
