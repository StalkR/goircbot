package url

import (
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/StalkR/goircbot/lib/transport"
)

var whitespaceRE = regexp.MustCompile(`\s+`)

// trim removes all white spaces including duplicates in a string.
func trim(s string) string {
	return strings.TrimSpace(whitespaceRE.ReplaceAllString(s, " "))
}

var htmlTagsRE = regexp.MustCompile(`<[^>]*>`)

// stripTags removes all HTML tags in a string.
func stripTags(s string) string {
	return htmlTagsRE.ReplaceAllString(s, "")
}

// get fetches an URL with a standard GET.
func get(url string) (string, error) {
	client, err := transport.Client(url)
	if err != nil {
		return "", err
	}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	return readBody(resp.Body)
}

const maxBody = 1 << 20 // 1MB

// readBody reads, closes and returns the first MB of an http response body.
func readBody(body io.ReadCloser) (string, error) {
	defer body.Close()
	// Limit reader to avoid denial of service fetching large files.
	lr := &io.LimitedReader{R: body, N: maxBody}
	contents, err := ioutil.ReadAll(lr)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}
