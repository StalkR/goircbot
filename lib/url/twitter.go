package url

import (
	"errors"
	"html"
	"regexp"
)

var (
	twitterRE = regexp.MustCompile(`^https?://twitter\.com/.*?/status/\d+(/photo/\d+)?(#|$)`)
	tweetRE   = regexp.MustCompile(`(?s)<div class="permalink-inner[^"]*">.*?<p class="js-tweet-text[^"]*">(.*?)</p>`)
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
	return Trim(html.UnescapeString(StripTags(text[1]))), nil
}
