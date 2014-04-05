// MLDonkey client for stats and add link by URL.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/StalkR/goircbot/lib/mldonkey"
)

var (
	url   = flag.String("url", "", "MLDonkey URL.")
	stats = flag.Bool("stats", false, "Show stats.")
	add   = flag.String("add", "", "Add link by URL.")
)

func main() {
	flag.Parse()
	if *url == "" {
		flag.Usage()
		os.Exit(1)
	}
	c, err := mldonkey.New(*url)
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
		if err := c.Add(*add); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v", err)
			os.Exit(1)
		}
		fmt.Println("Link added.")
	} else {
		flag.Usage()
		os.Exit(1)
	}
}
