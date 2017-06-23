// Package geo implements geographic functions such as location of an IP address.
package geo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/StalkR/goircbot/lib/transport"
)

// Locate locates an IP address or hostname using MaxMind demo.
// If a hostname it given, it first resolves to an IP.
// It is limited to 25 addresses per day.
func Locate(address string) (*Response, error) {
	if ip := net.ParseIP(address); ip != nil {
		return locationIP(ip.String())
	}
	addresses, err := net.LookupHost(address)
	if err != nil {
		return nil, err
	}
	if len(addresses) == 0 {
		return nil, errors.New("geoip: no IP addresses for this host")
	}
	ip := strings.TrimRight(addresses[0], ".")
	return locationIP(ip)
}

func locationIP(ip string) (*Response, error) {
	url := fmt.Sprintf("http://www.maxmind.com/geoip/v2.0/city_isp_org/%s?demo=1", ip)
	client, err := transport.Client(url)
	if err != nil {
		return nil, err
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var r Response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}
