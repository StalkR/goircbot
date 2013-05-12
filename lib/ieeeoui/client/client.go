// Client to know manufacturer of a MAC (IEEE public OUI).
package main

import (
	"fmt"
	"os"

	"github.com/StalkR/goircbot/lib/ieeeoui"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <mac address>\n", os.Args[0])
		os.Exit(1)
	}
	r := ieeeoui.New()
	manufacturer, err := r.Find(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(manufacturer)
}
