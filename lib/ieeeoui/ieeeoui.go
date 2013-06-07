// Package ieeeoui implements a library to know manufacturer of a MAC (IEEE public OUI).
package ieeeoui

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const ouiURL = "http://standards.ieee.org/develop/regauth/oui/oui.txt"

var (
	ouiRE = regexp.MustCompile(`^([\da-fA-F-]+)\s+\(hex\)\s+(.*)$`)
	rp    = strings.NewReplacer("-", "", ":", "")
)

func New() *Resolver {
	r := Resolver{manufacturer: make(map[int64]string)}
	r.Add(1)
	go r.init()
	return &r
}

type Resolver struct {
	sync.Mutex     // to protect WaitGroup for Add
	sync.WaitGroup // to have Find wait for init to complete
	manufacturer   map[int64]string
	err            error
}

func (r *Resolver) init() {
	defer r.Done()
	resp, err := http.Get(ouiURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	rd := bufio.NewReader(resp.Body)
	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println("ieeeoui: error reading OUI data: ", err)
			r.err = errors.New("could not read OUI data")
			return
		}
		m := ouiRE.FindStringSubmatch(strings.TrimSpace(line))
		if m == nil {
			continue
		}
		oui, err := strconv.ParseInt(rp.Replace(m[1]), 16, 0)
		if err != nil {
			continue
		}
		if v, ok := r.manufacturer[oui]; ok {
			r.manufacturer[oui] = v + " / " + m[2]
		} else {
			r.manufacturer[oui] = m[2]
		}
	}
	if len(r.manufacturer) == 0 {
		log.Println("ieeeoui: OUI data empty, parsing error?")
		r.err = errors.New("could not parse OUI data")
		return
	}
	r.err = nil
}

func (r *Resolver) Find(address string) (string, error) {
	hex := rp.Replace(address)
	if len(hex) < 6 {
		return "", errors.New("need at least 3 address bytes")
	}
	oui, err := strconv.ParseInt(hex[:6], 16, 0)
	if err != nil {
		return "", errors.New("invalid address")
	}
	r.Lock()
	defer r.Unlock()
	r.Wait()
	if r.err != nil {
		r.Add(1)
		go r.init() // retry
		return "", r.err
	}
	manufacturer, found := r.manufacturer[oui]
	if !found {
		return "", errors.New("not found")
	}
	return manufacturer, nil
}
