package old

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Old represents URLs seen.
type Old struct {
	sync.Mutex
	URLs map[string]Info
}

// Info represents information about an observed URL.
type Info struct {
	Channel, Nick string
	Time          time.Time
}

// NewOld returns a new initialized Old.
func NewOld() *Old {
	return &Old{
		URLs: make(map[string]Info),
	}
}

// Old returns info of an URL if it is old, error if it does not exist.
func (o *Old) Old(url string) (Info, error) {
	o.Lock()
	defer o.Unlock()
	if i, ok := o.URLs[url]; ok {
		return i, nil
	}
	return Info{}, errors.New("old: does not exist")
}

// Add adds a new URL with its info, error if it already exists.
func (o *Old) Add(url, channel, nick string) error {
	o.Lock()
	defer o.Unlock()
	if _, ok := o.URLs[url]; ok {
		return errors.New("old: already exists")
	}
	o.URLs[url] = Info{
		Channel: channel,
		Nick:    nick,
		Time:    time.Now(),
	}
	return nil
}

// Clean removes URLs older than a given duration.
func (o *Old) Clean(expiry time.Duration) {
	o.Lock()
	defer o.Unlock()
	for url, i := range o.URLs {
		if time.Since(i.Time) > expiry {
			delete(o.URLs, url)
		}
	}
}

// String returns formatted information about a URL.
func (i Info) String() string {
	duration := time.Since(i.Time) / time.Second * time.Second
	return fmt.Sprintf("old! first shared by %v %v ago", i.Nick, duration)
}
