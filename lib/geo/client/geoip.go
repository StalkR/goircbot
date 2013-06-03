// GeoIP client.
package main

import (
	"fmt"
	"os"

	"github.com/StalkR/goircbot/lib/geo"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %v <ip|host>\n", os.Args[0])
		os.Exit(1)
	}
	g, err := geo.Location(os.Args[1])
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println(g.String())
}
