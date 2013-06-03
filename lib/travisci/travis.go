// Package travisci implements a library to query Travis-CI builds using their JSON API.
package travisci

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/StalkR/goircbot/lib/tls"
)

// buildsURL is a format string taking user and repo to make the JSON builds URL.
const buildsURL = "https://api.travis-ci.org/repos/%s/%s/builds.json"

// Build represents a Travis-CI build.
// Compared to the JSON:
// 1) Started/Finished replace StartedAt/FinishedAt and are of type time.Time
// 2) Success replaces Result and is of type bool
// When build state is in progress (e.g. not "finished"), Started/Finished are
// not set (zero value of time.Time) because Travis-CI does not provide them.
type Build struct {
	Id           int
	RepositoryId int
	Number       int
	State        string
	Success      bool
	Started      time.Time
	Finished     time.Time
	Duration     int
	Commit       string
	Branch       string
	Message      string
	EventType    string
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

func timeoutDialer(d time.Duration) func(net, addr string) (net.Conn, error) {
	return func(netw, addr string) (net.Conn, error) {
		return net.DialTimeout(netw, addr, d)
	}
}

func httpClient(rawurl string) *http.Client {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}
	return &http.Client{
		Transport: &http.Transport{
			Dial:            timeoutDialer(5 * time.Second),
			TLSClientConfig: tls.Config(u.Host),
		},
	}
}

func Builds(user, repo string) ([]Build, error) {
	url := fmt.Sprintf(buildsURL, user, repo)
	resp, err := httpClient(url).Get(url)
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
		if b.Result != nil {
			if *b.Result != 0 {
				panic(fmt.Sprintf("build result has non-zero value %d", *b.Result))
			}
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
		})
	}
	return builds, nil
}
