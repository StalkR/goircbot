// Binary ieeeoui generates a Go file with mapping of OUI
// and organization from IEEE public OUI.
package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const ouiURL = "http://standards.ieee.org/develop/regauth/oui/oui.txt"

var (
	pkg     = flag.String("p", "mac", "Package name.")
	varname = flag.String("v", "ieeeoui", "Variable name.")
	out     = flag.String("o", "ieeeoui.go", "Output file.")
)

func main() {
	flag.Parse()
	ouis, err := Get()
	if err != nil {
		log.Fatal(err)
	}
	var b bytes.Buffer
	fmt.Fprintf(&b, "package %s\n", *pkg)
	b.WriteString("\n")
	fmt.Fprintf(&b, "var %s = map[uint64]string{\n", *varname)
	for _, e := range ouis {
		fmt.Fprintf(&b, "\t0x%06x: %#v,\n", e.ID, e.Org)
	}
	b.WriteString("}\n")
	if err := ioutil.WriteFile(*out, b.Bytes(), 0660); err != nil {
		log.Fatal(err)
	}
}

var (
	ouiRE       = regexp.MustCompile(`^([\da-fA-F-]+)\s+\(hex\)\s+(.*)$`)
	ignoreChars = strings.NewReplacer("-", "", ":", "")
)

// Get obtains the XML reference, parses it and returns a slice of mnemonics
// sorted by name and with a brief description.
func Get() ([]OUI, error) {
	resp, err := http.Get(ouiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	set := make(map[uint64]string)
	for scanner.Scan() {
		m := ouiRE.FindStringSubmatch(strings.TrimSpace(scanner.Text()))
		if m == nil {
			continue
		}
		oui, err := strconv.ParseUint(ignoreChars.Replace(m[1]), 16, 0)
		if err != nil {
			continue
		}
		if v, ok := set[oui]; ok {
			set[oui] = v + ", " + m[2]
		} else {
			set[oui] = m[2]
		}

	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(set) == 0 {
		return nil, errors.New("OUI data empty")
	}
	var ouis []OUI
	for id, org := range set {
		ouis = append(ouis, OUI{ID: id, Org: org})
	}
	sort.Sort(byID(ouis))
	return ouis, nil
}

// An OUI represents an ID assigned by IEEE for a given organization.
type OUI struct {
	ID  uint64
	Org string
}

type byID []OUI

func (s byID) Len() int           { return len(s) }
func (s byID) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byID) Less(i, j int) bool { return s[i].ID < s[j].ID }
