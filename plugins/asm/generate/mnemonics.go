// Binary mnemonics generates a Go file with mapping of mnemonics
// and their description from the x86/x86-64 reference.
package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
)

const xmlURL = "http://ref.x86asm.net/x86reference.xml"

var (
	pkg     = flag.String("p", "asm", "Package name.")
	varname = flag.String("v", "mnemonics", "Variable name.")
	out     = flag.String("o", "mnemonics.go", "Output file.")
)

func main() {
	flag.Parse()
	mnemonics, err := Get()
	if err != nil {
		log.Fatal(err)
	}
	var b bytes.Buffer
	fmt.Fprintf(&b, "package %s\n", *pkg)
	b.WriteString("\n")
	fmt.Fprintf(&b, "var %s = map[string]string{\n", *varname)
	for _, m := range mnemonics {
		fmt.Fprintf(&b, "\t%#v: %#v,\n", m.Name, m.Brief)
	}
	b.WriteString("}\n")
	if err := ioutil.WriteFile(*out, b.Bytes(), 0660); err != nil {
		log.Fatal(err)
	}
}

// Get obtains the XML reference, parses it and returns a slice of mnemonics
// sorted by name and with a brief description.
func Get() ([]Mnemonic, error) {
	resp, err := http.Get(xmlURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	page, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var x xmlRef
	if err := xml.Unmarshal(page, &x); err != nil {
		return nil, err
	}
	var mnemonics []Mnemonic
	for name, brief := range x.Mnemonics() {
		mnemonics = append(mnemonics, Mnemonic{Name: name, Brief: brief})
	}
	sort.Sort(byName(mnemonics))
	return mnemonics, nil
}

// An xmlRef represents the partially-parsed XML reference.
type xmlRef struct {
	OneByte byteEl `xml:"one-byte"`
	TwoByte byteEl `xml:"two-byte"`
}

type byteEl struct {
	OpCode []opCode `xml:"pri_opcd"`
}

type opCode struct {
	Entry []entry `xml:"entry"`
}

type entry struct {
	Syntax []syntax `xml:"syntax"`
	Brief  string   `xml:"note>brief"`
}

type syntax struct {
	Mnem string `xml:"mnem"`
}

func (r *xmlRef) Mnemonics() map[string]string {
	set := make(map[string]string)
	for _, a := range []byteEl{r.OneByte, r.TwoByte} {
		for _, b := range a.OpCode {
			for _, c := range b.Entry {
				for _, d := range c.Syntax {
					if d.Mnem == "" {
						continue
					}
					set[d.Mnem] = c.Brief
				}
			}
		}
	}
	return set
}

// A Mnemonic represents an x86/x86-64 mnemonic with its name and description brief.
type Mnemonic struct {
	Name  string
	Brief string
}

type byName []Mnemonic

func (s byName) Len() int           { return len(s) }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byName) Less(i, j int) bool { return s[i].Name < s[j].Name }
