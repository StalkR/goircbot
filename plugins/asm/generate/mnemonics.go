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
)

const xmlURL = "http://ref.x86asm.net/x86reference.xml"

var (
	pkg     = flag.String("p", "", "Package name.")
	varname = flag.String("v", "mnemonics", "Variable name.")
	out     = flag.String("o", "", "Output file.")
)

func main() {
	flag.Parse()
	resp, err := http.Get(xmlURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	page, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var x xmlRef
	if err := xml.Unmarshal(page, &x); err != nil {
		log.Fatal(err)
	}
	var b bytes.Buffer
	fmt.Fprintf(&b, "package %s\n", *pkg)
	b.WriteString("\n")
	fmt.Fprintf(&b, "var %s = map[string]string{\n", *varname)
	for mnem, brief := range x.Mnemonics() {
		fmt.Fprintf(&b, "\t%#v: %#v,\n", mnem, brief)
	}
	b.WriteString("}\n")
	if err := ioutil.WriteFile(*out, b.Bytes(), 0660); err != nil {
		log.Fatal(err)
	}
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
	m := make(map[string]string)
	for _, a := range []byteEl{r.OneByte, r.TwoByte} {
		for _, b := range a.OpCode {
			for _, c := range b.Entry {
				for _, d := range c.Syntax {
					if d.Mnem == "" {
						continue
					}
					m[d.Mnem] = c.Brief
				}
			}
		}
	}
	return m
}
