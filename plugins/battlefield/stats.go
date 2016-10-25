package battlefield

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type BFStats struct {
	ID          uint64
	Game        string // "bf1", "bf4"
	TimePlayed  time.Duration
	Kills       uint
	Deaths      uint
	Wins        uint
	Losses      uint
	Rank        uint
	ScorePerMin float64
}

func (s *BFStats) URL() string {
	return fmt.Sprintf("https://www.battlefield.com/companion/career/%d/%s", s.ID, s.Game)
}

func (s *BFStats) String() string {
	kd := float64(s.Kills) / float64(s.Deaths)
	wl := 100 * s.Wins / (s.Wins + s.Losses)
	return fmt.Sprintf("rank %d, %d%% wins, %d kills, %d deaths, KD %.2f, SPM %.2f, %s play time %s",
		s.Rank, wl, s.Kills, s.Deaths, kd, s.ScorePerMin, s.TimePlayed, s.URL())
}

func (s *BFStats) Short() string {
	kd := float64(s.Kills) / float64(s.Deaths)
	return fmt.Sprintf("rank %d KD %.2f SPM %.2f",
		s.Rank, kd, s.ScorePerMin)
}

type Stats struct {
	ID   uint64
	Name string // TODO get name automatically
	BF1  BFStats
	BF4  BFStats
}

type byRankBF1 []Stats

func (a byRankBF1) Len() int      { return len(a) }
func (a byRankBF1) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byRankBF1) Less(i, j int) bool {
	if a[i].BF1.Rank == a[j].BF1.Rank {
		return a[i].BF1.ScorePerMin < a[j].BF1.ScorePerMin
	}
	return a[i].BF1.Rank < a[j].BF1.Rank
}

type byRankBF4 []Stats

func (a byRankBF4) Len() int      { return len(a) }
func (a byRankBF4) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byRankBF4) Less(i, j int) bool {
	if a[i].BF4.Rank == a[j].BF4.Rank {
		return a[i].BF4.ScorePerMin < a[j].BF4.ScorePerMin
	}
	return a[i].BF4.Rank < a[j].BF4.Rank
}

var errNotFound = errors.New("battlefield: player not found")

func parseStats(personaID uint64, name string, js []byte) (*Stats, error) {
	var result struct {
		Result struct {
			GameStats struct {
				Tunguska struct { // bf1
					TimePlayed uint
					Kills      uint
					Deaths     uint
					Wins       uint
					Losses     uint
					Rank       struct {
						Number uint
					}
					Spm float64
				}
				Bf4 struct {
					TimePlayed uint
					Kills      uint
					Deaths     uint
					Wins       uint
					Losses     uint
					Rank       struct {
						Number uint
					}
					Spm float64
				}
			}
		}
		Error struct {
			Message string
		}
	}
	if err := json.Unmarshal(js, &result); err != nil {
		return nil, err
	}
	if result.Error.Message == "Internal Error: java.util.NoSuchElementException" {
		return nil, errNotFound
	}
	if result.Error.Message != "" {
		return nil, fmt.Errorf("battlefield: server error: %s", result.Error.Message)
	}
	bf1 := result.Result.GameStats.Tunguska
	bf4 := result.Result.GameStats.Bf4
	return &Stats{
		ID:   personaID,
		Name: name,
		BF1: BFStats{
			ID:          personaID,
			Game:        "bf1",
			TimePlayed:  time.Duration(bf1.TimePlayed) * time.Second,
			Kills:       bf1.Kills,
			Deaths:      bf1.Deaths,
			Wins:        bf1.Wins,
			Losses:      bf1.Losses,
			Rank:        bf1.Rank.Number,
			ScorePerMin: bf1.Spm,
		},
		BF4: BFStats{
			ID:          personaID,
			Game:        "bf4",
			TimePlayed:  time.Duration(bf4.TimePlayed) * time.Second,
			Kills:       bf4.Kills,
			Deaths:      bf4.Deaths,
			Wins:        bf4.Wins,
			Losses:      bf4.Losses,
			Rank:        bf4.Rank.Number,
			ScorePerMin: bf4.Spm,
		},
	}, nil
}
