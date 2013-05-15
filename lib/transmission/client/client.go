// Transmission client for stats and add torrent by URL.
package main

import (
	"fmt"
	"os"

	"github.com/StalkR/goircbot/lib/transmission"
)

func usage() {
	fmt.Printf("Usage: %s <Transmission URL> (stats|add <URL>)\n", os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) < 3 {
		usage()
	}

	c, err := transmission.New(os.Args[1])
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
		name, err := c.Add(os.Args[3])
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Torrent added: ", name)

	default:
		usage()
	}
}
