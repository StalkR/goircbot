// Urban Dictionary client.
package main

import (
	"fmt"
	"os"

	"github.com/StalkR/goircbot/lib/urbandictionary"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %v <query>\n", os.Args[0])
		os.Exit(1)
	}
	r, err := urbandictionary.Define(os.Args[1])
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println(r.String())
}
