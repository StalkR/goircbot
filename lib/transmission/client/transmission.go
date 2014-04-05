// Transmission client for stats and add torrent by URL.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/StalkR/goircbot/lib/transmission"
)

var (
	url   = flag.String("url", "", "Transmission URL.")
	stats = flag.Bool("stats", false, "Show stats.")
	add   = flag.String("add", "", "Add link by URL.")
)

func main() {
	flag.Parse()
	if *url == "" {
		flag.Usage()
		os.Exit(1)
	}
	c, err := transmission.New(*url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
	if *stats {
		stats, err := c.Stats()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v", err)
			os.Exit(1)
		}
		fmt.Println(stats.String())
	} else if *add != "" {
		name, err := c.Add(*add)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v", err)
			os.Exit(1)
		}
		fmt.Println("Torrent added: ", name)
	} else {
		flag.Usage()
		os.Exit(1)
	}
}
