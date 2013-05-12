// MLDonkey client to print stats.
package main

import (
	"fmt"
	"os"

	"github.com/StalkR/goircbot/lib/mldonkey"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <mldonkey url>\n", os.Args[0])
		os.Exit(1)
	}
	s, err := mldonkey.Stats(os.Args[1])
	if err != nil {
		fmt.Println("err", err)
	}
	fmt.Println(s.String())
}
