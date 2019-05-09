package url

import (
	"errors"
	"html"
	"net/http"
	"regexp"
	"strings"

	"github.com/StalkR/goircbot/lib/transport"
)

var (
	twitterRE = regexp.MustCompile(`^https?://twitter\.com/.*?/status/\d+(/photo/\d+)?(#|$)`)
	tweetRE   = regexp.MustCompile(`(?s)<div class="permalink-inner[^"]*">.*?<p class="[^"]*?js-tweet-text[^"]*"[^>]*>(.*?)</p>`)
	twitpicRE = regexp.MustCompile(`([^ ])<a[^>]*>(pic.twitter.com/[^<]*)</a>`)
)

func handleTwitter(url string) (string, error) {
	if !twitterRE.MatchString(url) {
		return "", errSkip
	}
	client, err := transport.Client(url)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	// Twitter now serves a dummy page with JavaScript location.replace
	// to the same path if not coming with the right referer.
	req.Header.Add("Referer", url)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	body, err := readBody(resp.Body)
	if err != nil {
		return "", err
	}
	text := tweetRE.FindStringSubmatch(body)
	if text == nil {
		return "", errors.New("url: twitter: cannot parse tweet")
	}
	s := text[1]
	// insert a space before pics and add scheme
	s = twitpicRE.ReplaceAllString(s, "$1 https://$2")
	s = stripTags(s)
	// unescape would replace &nbsp; by \u00a0 but we prefer normal space \u0020
	s = strings.Replace(s, "&nbsp;", " ", -1)
	s = html.UnescapeString(s)
	s = trim(s)
	return s, nil
}
