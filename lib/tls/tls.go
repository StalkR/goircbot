// Package tls configures TLS with root CAs given by flag.
package tls

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"strings"
	"sync"
)

var (
	cacerts = flag.String("cacerts", "",
		"Root CA PEM files (separated by comma), if empty use system default")
	once                     = sync.Once{}
	rootCerts *x509.CertPool = nil
)

// Config takes an optional server name (for SNI) and returns a TLS config with Root CAs set.
func Config(serverName string) *tls.Config {
	once.Do(func() {
		if *cacerts == "" {
			return
		}
		rootCerts = x509.NewCertPool()
		for _, file := range strings.Split(*cacerts, ",") {
			b, err := ioutil.ReadFile(file)
			if err != nil {
				panic("tls: could not read CA file " + file + ": " + err.Error())
			}
			rootCerts.AppendCertsFromPEM(b)
		}
	})
	if rootCerts == nil {
		return &tls.Config{ServerName: serverName}
	}
	return &tls.Config{ServerName: serverName, RootCAs: rootCerts}
}
