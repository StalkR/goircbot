// Binary client searches for metal bands.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/StalkR/goircbot/lib/metal/all"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %v <band>\n", os.Args[0])
		os.Exit(1)
	}
	bands, err := all.Search(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	for _, band := range bands {
		fmt.Println(band.String())
	}
}
