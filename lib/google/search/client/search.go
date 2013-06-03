// Google Custom Search client.
package main

import (
	"fmt"
	"os"

	"github.com/StalkR/goircbot/lib/google/search"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Printf("Usage: %v <key> <cx> <query>\n", os.Args[0])
		os.Exit(1)
	}
	r, err := search.Search(os.Args[3], os.Args[1], os.Args[2])
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	for i, item := range r.Items {
		fmt.Printf("%2d. %s\n", i+1, item.String())
	}
}
