// Package darkstat provides a library to read information from darkstat web server.
// Darkstat: Captures network traffic, calculates statistics about usage, and serves reports over HTTP.
// http://unix4lyfe.org/darkstat/
package darkstat

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"

	"github.com/StalkR/goircbot/lib/transport"
)

// A Conn represents a connection to Darkstat.
type Conn struct {
	url    string
	client *http.Client
}

// New prepares a Darkstat connection by returning a *Conn.
func New(serverURL string) (*Conn, error) {
	client, err := transport.Client(serverURL)
	if err != nil {
		return nil, err
	}
	return &Conn{url: serverURL, client: client}, nil
}

// Graphs gets and parses darkstat's graphs.xml page.
func (c *Conn) Graphs() (*GraphsXML, error) {
	resp, err := c.client.Get(c.url + "/graphs.xml")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	graphs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var g GraphsXML
	if err := xml.Unmarshal(graphs, &g); err != nil {
		return nil, err
	}
	return &g, nil
}

// A GraphsXML represents darkstat's graphs.xml.
type GraphsXML struct {
	Packets    int         `xml:"tp,attr"`
	Bytes      int         `xml:"tb,attr"`
	Captured   int         `xml:"pc,attr"`
	Dropped    int         `xml:"pd,attr"`
	RunningFor string      `xml:"rf,attr"`
	Seconds    []DataPoint `xml:"seconds>e"`
	Minutes    []DataPoint `xml:"minutes>e"`
	Hours      []DataPoint `xml:"hours>e"`
	Days       []DataPoint `xml:"days>e"`
}

// An DataPoint represents a data point of bytes in/out at a given instant.
type DataPoint struct {
	Pos int `xml:"p,attr"`
	In  int `xml:"i,attr"`
	Out int `xml:"o,attr"`
}

// Bandwidth returns the average bandwidth in bytes observed in the last minute.
func (g *GraphsXML) Bandwidth() (in, out int) {
	for _, d := range g.Seconds {
		in += d.In
		out += d.Out
	}
	return in / 60, out / 60
}

// Bandwidth returns the average bandwidth in bytes observed in the last minute.
func (c *Conn) Bandwidth() (in, out int, err error) {
	g, err := c.Graphs()
	if err != nil {
		return
	}
	in, out = g.Bandwidth()
	return
}
