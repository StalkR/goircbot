package quotes

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// A Quote represents a quote text, added by a nick at a given time.
type Quote struct {
	ID              int    `json:"-"` // ignore ID in JSON as it's the map key
	Added, By, Text string `json:",omitempty"`
}

// String formats a single quote.
func (q Quote) String() string {
	if len(q.By) > 0 {
		return fmt.Sprintf("%s [#%d %s %s]", q.Text, q.ID, q.By, q.Added)
	}
	return fmt.Sprintf("%s [#%d]", q.Text, q.ID)
}

// ByID implements sort.Interface for []Quote based on the ID field.
type ByID []Quote

func (s ByID) Len() int           { return len(s) }
func (s ByID) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByID) Less(i, j int) bool { return s[i].ID < s[j].ID }

// Quotes is the main structure to hold quotes in the plugin.
type Quotes struct {
	sync.Mutex
	quotes map[string]Quote // indexed by ID, as string for JSON
	next   int              // next unused ID
	dirty  bool
}

// NewQuotes creates a new initialized Quotes.
func NewQuotes() *Quotes {
	q := &Quotes{}
	q.quotes = make(map[string]Quote)
	return q
}

// Add adds a new quote and return it.
func (q *Quotes) Add(by, text string) Quote {
	q.Lock()
	defer q.Unlock()
	e := Quote{ID: q.next, Added: time.Now().Format("2006-01-02"), By: by, Text: text}
	q.quotes[strconv.Itoa(q.next)] = e
	q.next++
	q.dirty = true
	return e
}

// Delete removes a quote by its ID and returns whether it was there.
func (q *Quotes) Delete(id int) bool {
	q.Lock()
	defer q.Unlock()
	ids := strconv.Itoa(id)
	_, present := q.quotes[ids]
	if !present {
		return false
	}
	delete(q.quotes, ids)
	q.dirty = true
	return true
}

// Search finds quotes with simple contains or regexp.
func (q *Quotes) Search(term string) []Quote {
	q.Lock()
	defer q.Unlock()
	var results []Quote
	for _, quote := range q.quotes {
		if !strings.Contains(quote.Text, term) {
			re, err := regexp.Compile(term)
			if err != nil || !re.MatchString(quote.Text) {
				continue
			}
		}
		results = append(results, quote)
	}
	// Fisher–Yates shuffle with inside-out algorithm
	for i := range results {
		j := rand.Intn(i + 1)
		results[i], results[j] = results[j], results[i]
	}
	return results
}

// Empty returns whether there are no quotes yet.
func (q *Quotes) Empty() bool {
	q.Lock()
	defer q.Unlock()
	return len(q.quotes) == 0
}

// Random picks a random quote (empty if no quotes).
func (q *Quotes) Random() Quote {
	q.Lock()
	defer q.Unlock()
	if len(q.quotes) == 0 {
		return Quote{}
	}
	var ids []string
	for id := range q.quotes {
		ids = append(ids, id)
	}
	randID := ids[rand.Intn(len(ids))]
	return q.quotes[randID]
}
