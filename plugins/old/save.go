package old

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func load(oldfile string) *Old {
	o := NewOld()
	if len(oldfile) == 0 {
		return o
	}
	b, err := ioutil.ReadFile(oldfile)
	if err != nil {
		log.Println("old: unable to open old file")
		return o
	}
	if err := json.Unmarshal(b, &o.urls); err != nil {
		log.Println("old: unable to load old")
		return o
	}
	log.Println("old: loaded successfully")
	return o
}

func save(oldfile string, o *Old) {
	o.Lock()
	defer o.Unlock()
	if !o.dirty {
		return
	}
	b, err := json.Marshal(o.urls)
	if err != nil {
		log.Println("old: unable to encode old for saving")
		return
	}
	if err := ioutil.WriteFile(oldfile, b, 0644); err != nil {
		log.Println("old: unable to save old")
		return
	}
	o.dirty = false
}
