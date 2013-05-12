// IMDb client.
package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/StalkR/goircbot/lib/imdb"
)

var ttRE = regexp.MustCompile(`^tt\w+$`)

// findTitle searches for a title and presents up to 10 results.
func findTitle(q string) {
	r, err := imdb.FindTitle(q)
	if err != nil {
		fmt.Println("FindTitle error", err)
		os.Exit(1)
	}
	if len(r) == 0 {
		fmt.Println("No results found.")
		return
	}
	max := len(r)
	if max > 10 {
		max = 10
	}
	for i, tt := range r[:max] {
		t, err := imdb.NewTitle(tt.Id)
		if err != nil {
			fmt.Println("NewTitle error", err)
			os.Exit(1)
		}
		fmt.Printf("%2d. %s\n", i+1, t.String())
	}
}

// title obtains information on a title id and presents it.
func title(id string) {
	t, err := imdb.NewTitle(id)
	if err != nil {
		fmt.Println("NewTitle error", err)
		os.Exit(1)
	}
	fmt.Println(t.String())
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <query|ttID>\n", os.Args[0])
		os.Exit(1)
	}
	if ttRE.MatchString(os.Args[1]) {
		title(os.Args[1])
		return
	}
	findTitle(os.Args[1])
}
