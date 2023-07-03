// Package url implements a library to get a meaningful title of web URLs.
package url

import "errors"

// errSkip is used by handlers to skip to the next handler.
var errSkip = errors.New("url: skip to next handler")

// handlers is the ordered list of handlers.
var handlers = []func(url string) (string, error){}

// Title gets an URL and returns its title.
func Title(url string) (string, error) {
	for _, handler := range handlers {
		title, err := handler(url)
		if err == errSkip {
			continue
		}
		return title, err
	}
	return handleDefault(url)
}
