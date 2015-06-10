/*
Package x86 exposes the X86 Opcode and Instruction Reference
available at http://ref.x86asm.net/x86reference.xml
*/
package x86

import (
	"encoding/xml"
	"io/ioutil"

	"github.com/StalkR/goircbot/lib/transport"
)

const xmlURL = "http://ref.x86asm.net/x86reference.xml"

// New fetches the reference XML, parses it and returns a Reference.
func New() (*Reference, error) {
	c, err := transport.Client(xmlURL)
	if err != nil {
		return nil, err
	}
	resp, err := c.Get(xmlURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var x xmlRef
	if err := xml.Unmarshal(b, &x); err != nil {
		return nil, err
	}
	return &Reference{xml: x}, nil
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

// A Reference represents the x86/x86-64 reference and exposes methods to access its data.
type Reference struct {
	xml xmlRef
}

// Mnemonics builds the map of instructions with their description.
func (r *Reference) Mnemonics() map[string]string {
	m := make(map[string]string)
	for _, a := range []byteEl{r.xml.OneByte, r.xml.TwoByte} {
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
