// Package travisci implements a library to query Travis-CI builds using their JSON API.
package travisci

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/StalkR/goircbot/lib/transport"
)

// Build represents a Travis-CI build.
// Compared to the JSON:
// 1) Started/Finished replace StartedAt/FinishedAt and are of type time.Time
// 2) Success replaces Result and is of type bool
// When build state is in progress (e.g. not "finished"), Started/Finished are
// not set (zero value of time.Time) because Travis-CI does not provide them.
type Build struct {
	Id, RepositoryId, Number int
	State                    string
	Success                  bool
	Started, Finished        time.Time
	Duration                 int
	Commit, Branch, Message  string
	EventType                string
	BuildURL, CommitURL      string
}

func (b *Build) String() string {
	var status string
	if b.State == "finished" {
		status = "passed"
		if !b.Success {
			status = "errored"
		}
	} else {
		status = "in progress"
	}
	if b.Finished.IsZero() {
		return fmt.Sprintf("Build #%v: %v %v %v %v", b.Number, status,
			b.BuildURL, b.CommitURL, b.Message)
	}
	return fmt.Sprintf("Build #%v: %v (%v) %v %v %v", b.Number, status,
		b.Finished.Format("2006-01-02 15:04:05 UTC"), b.BuildURL, b.CommitURL, b.Message)
}

// buildJSON represents the builds.json replied by Travis-CI API.
// Some fields are pointers because they can be null, and we want to differentiate
// from value 0.
type buildJSON struct {
	Id           int    `json:"id"`
	RepositoryId int    `json:"repository_id"`
	Number       int    `json:"number,string"`
	State        string `json:"state"`
	Result       *int   `json:"result"`
	StartedAt    string `json:"started_at"`
	FinishedAt   string `json:"finished_at"`
	Duration     int    `json:"duration"`
	Commit       string `json:"commit"`
	Branch       string `json:"branch"`
	Message      string `json:"message"`
	EventType    string `json:"event_type"`
}

func Builds(user, repo string) ([]Build, error) {
	url := fmt.Sprintf("https://api.travis-ci.org/repos/%s/%s/builds.json", user, repo)
	client, err := transport.Client(url)
	if err != nil {
		return nil, err
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	js, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var bjs []buildJSON
	if err := json.Unmarshal(js, &bjs); err != nil {
		return nil, err
	}
	var builds []Build
	for _, b := range bjs {
		var startedAt, finishedAt time.Time
		if b.StartedAt != "" {
			startedAt, err = time.Parse("2006-01-02T15:04:05Z", b.StartedAt)
			if err != nil {
				return nil, err
			}
		}
		if b.FinishedAt != "" {
			finishedAt, err = time.Parse("2006-01-02T15:04:05Z", b.FinishedAt)
			if err != nil {
				return nil, err
			}
		}

		success := false
		if b.Result != nil && *b.Result == 0 {
			success = true
		}

		builds = append(builds, Build{
			Id:           b.Id,
			RepositoryId: b.RepositoryId,
			Number:       b.Number,
			State:        b.State,
			Success:      success,
			Started:      startedAt,
			Finished:     finishedAt,
			Duration:     b.Duration,
			Commit:       b.Commit,
			Branch:       b.Branch,
			Message:      b.Message,
			EventType:    b.EventType,
			BuildURL:     fmt.Sprintf("https://travis-ci.org/%v/%v/builds/%v", user, repo, b.Id),
			CommitURL:    fmt.Sprintf("https://github.com/%v/%v/commit/%v", user, repo, b.Commit),
		})
	}
	return builds, nil
}
