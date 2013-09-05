// Wunderground client.
package main

import (
	"fmt"
	"os"

	"github.com/StalkR/goircbot/lib/wunderground"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %v <API key> <location>\n", os.Args[0])
		os.Exit(1)
	}
	w, err := wunderground.Conditions(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println(w.String())
}
