package old

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/StalkR/goircbot/lib/duration"
)

// Old represents URLs seen.
type Old struct {
	sync.Mutex
	urls  map[string]Info
	dirty bool
}

// Info represents information about an observed URL.
type Info struct {
	Channel, Nick string
	Time          time.Time
}

// NewOld returns a new initialized Old.
func NewOld() *Old {
	return &Old{urls: make(map[string]Info)}
}

// Old returns info of an URL if it is old, error if it does not exist.
func (o *Old) Old(url string) (Info, error) {
	o.Lock()
	defer o.Unlock()
	if i, ok := o.urls[url]; ok {
		return i, nil
	}
	return Info{}, errors.New("old: does not exist")
}

// Add adds a new URL with its info, error if it already exists.
func (o *Old) Add(url, channel, nick string) error {
	o.Lock()
	defer o.Unlock()
	if _, ok := o.urls[url]; ok {
		return errors.New("old: already exists")
	}
	o.urls[url] = Info{
		Channel: channel,
		Nick:    nick,
		Time:    time.Now(),
	}
	o.dirty = true
	return nil
}

// Clean removes URLs older than a given duration.
func (o *Old) Clean(expiry time.Duration) {
	o.Lock()
	defer o.Unlock()
	for url, i := range o.urls {
		if time.Since(i.Time) > expiry {
			delete(o.urls, url)
			o.dirty = true
		}
	}
}

// String returns formatted information about a URL.
func (i Info) String() string {
	return fmt.Sprintf("old! first shared by %v %v ago",
		i.Nick, duration.Format(time.Since(i.Time)))
}
