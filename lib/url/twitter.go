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
	twitpicRE = regexp.MustCompile(`([^ ])<a[^>]*>(pic.twitter.com/[^<]*)</a>`)
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
	s := text[1]
	// insert a space before pics, and add http://
	s = twitpicRE.ReplaceAllString(s, "$1 http://$2")
	s = StripTags(s)
	// unescape would replace &nbsp; by \u00a0 but we prefer normal space \u0020
	s = strings.Replace(s, "&nbsp;", " ", -1)
	s = html.UnescapeString(s)
	s = Trim(s)
	return s, nil
}
