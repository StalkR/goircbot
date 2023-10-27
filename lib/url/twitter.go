package url

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	twitter "github.com/g8rswimmer/go-twitter"
)

var (
	twitterRE = regexp.MustCompile(`^https?://twitter\.com/.*?/status/(\d+)(?:/photo/\d+)?(?:#|$)`)
)

var TwitterAPIToken string

func handleTwitter(target string) (string, error) {
	m := twitterRE.FindStringSubmatch(target)
	if len(m) != 2 || TwitterAPIToken == "" {
		return "", errSkip
	}
	id := m[1]
	return lookup(id)
}

func lookup(id string) (string, error) {
	tweet := &twitter.Tweet{
		Authorizer: authorize{},
		Client:     http.DefaultClient,
		Host:       "https://api.twitter.com",
	}

	lookups, err := tweet.Lookup(context.Background(), []string{id}, twitter.TweetFieldOptions{})
	if err != nil {
		return "", err
	}
	if len(lookups) != 1 {
		return "", fmt.Errorf("expected 1 tweet, got %v", len(lookups))
	}
	for _, v := range lookups {
		return fmt.Sprintf("<%v> %v", v.Tweet.AuthorID, v.Tweet.Text), nil
	}
	panic("unreachable")
}

type authorize struct{}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", TwitterAPIToken))
}
