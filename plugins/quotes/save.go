package quotes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
)

func load(f string) *Quotes {
	q := NewQuotes()
	if len(f) == 0 {
		return q
	}
	b, err := ioutil.ReadFile(f)
	if err != nil {
		log.Println("quotes: unable to open quotes file")
		return q
	}
	if err := json.Unmarshal(b, &q.quotes); err != nil {
		log.Println("quotes: unable to load quotes")
		return q
	}
	for k, v := range q.quotes {
		id, err := strconv.Atoi(k)
		if err != nil { // ignore non-numeric keys if any
			delete(q.quotes, k)
			continue
		}
		// restore ID from the map key
		v.ID = id
		q.quotes[k] = v
		if id >= q.next {
			q.next = id + 1
		}
	}
	log.Println("quotes: loaded successfully")
	return q
}

func save(f string, q *Quotes) {
	q.Lock()
	defer q.Unlock()
	if !q.dirty {
		return
	}
	b, err := json.MarshalIndent(q.quotes, "", "  ")
	if err != nil {
		log.Println("quotes: unable to encode quotes for saving")
		return
	}
	if err := ioutil.WriteFile(f, b, 0644); err != nil {
		log.Println("quotes: unable to save quotes")
		return
	}
	q.dirty = false
}
