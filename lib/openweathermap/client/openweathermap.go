// Binary openweathermap gets weather information for a location.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/StalkR/goircbot/lib/openweathermap"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %v <API key> <location>\n", os.Args[0])
		os.Exit(1)
	}
	w, err := openweathermap.Find(os.Args[1], os.Args[2])
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	fmt.Println(w.String())
}
