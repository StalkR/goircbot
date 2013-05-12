// Transmission client to print stats.
package main

import (
	"fmt"
	"os"

	"github.com/StalkR/goircbot/lib/transmission"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <transmission url>\n", os.Args[0])
		os.Exit(1)
	}
	s, err := transmission.Stats(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(s.String())
}
