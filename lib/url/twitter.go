package url

import (
	"errors"
	"html"
	"regexp"
)

var (
	twitterRE = regexp.MustCompile(`^https?://twitter\.com/.*?/status/\d+(#|$)`)
	tweetRE   = regexp.MustCompile(`<p class="js-tweet-text">(.*?)</p>`)
)

type twitter struct{}

func (d *twitter) Match(url string) bool {
	return twitterRE.MatchString(url)
}

func (d *twitter) Parse(body string) (string, error) {
	text := tweetRE.FindStringSubmatch(body)
	if text == nil {
		return "", errors.New("url: twitter: cannot parse tweet")
	}
	return Trim(html.UnescapeString(StripTags(text[1]))), nil
}
