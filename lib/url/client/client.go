// Binary client shows the title of an url.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/StalkR/goircbot/lib/url"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %v <url>\n", os.Args[0])
		os.Exit(1)
	}
	title, err := url.Title(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(title)
}
