// MLDonkey client for stats and add link by URL.
package main

import (
	"fmt"
	"os"

	"github.com/StalkR/goircbot/lib/mldonkey"
)

func usage() {
	fmt.Printf("Usage: %s <MLDonkey URL> (stats|add <URL>)\n", os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) < 3 {
		usage()
	}

	c, err := mldonkey.New(os.Args[1])
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	switch os.Args[2] {

	case "stats":
		stats, err := c.Stats()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println(stats.String())

	case "add":
		if len(os.Args) < 4 {
			usage()
		}
		if err := c.Add(os.Args[3]); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Link added.")

	default:
		usage()
	}
}
