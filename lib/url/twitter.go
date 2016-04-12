package url

import (
	"errors"
	"html"
	"regexp"
	"strings"
)

var (
	twitterRE = regexp.MustCompile(`^https?://twitter\.com/.*?/status/\d+(/photo/\d+)?(#|$)`)
	tweetRE   = regexp.MustCompile(`(?s)<div class="permalink-inner[^"]*">.*?<p class="[^"]*?js-tweet-text[^"]*"[^>]*>(.*?)</p>`)
)

type Twitter struct{}

func (p *Twitter) Match(url string) bool {
	return twitterRE.MatchString(url)
}

func (p *Twitter) Parse(body string) (string, error) {
	text := tweetRE.FindStringSubmatch(body)
	if text == nil {
		return "", errors.New("url: twitter: cannot parse tweet")
	}
	s := StripTags(text[1])
	// unescape would replace &nbsp; by \u00a0 but we prefer normal space \u0020
	s = strings.Replace(s, "&nbsp;", " ", -1)
	s = html.UnescapeString(s)
	s = Trim(s)
	return s, nil
}
