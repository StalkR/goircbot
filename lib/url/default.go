package url

import (
	"errors"
	"html"
	"regexp"
)

var titleRE = regexp.MustCompile(`(?i)<title[^>]*>([^<]+)<`)

type Default struct{}

func (p *Default) Match(url string) bool { return true }

func (p *Default) Parse(body string) (string, error) {
	text := titleRE.FindStringSubmatch(body)
	if text == nil {
		return "", errors.New("url: cannot parse title")
	}
	return Trim(html.UnescapeString(text[1])), nil
}
