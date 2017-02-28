// Client is a binary to test battlefield plugin.
// go run client.go -email x -password y -id 0 -game bf1
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/StalkR/goircbot/plugins/battlefield"
)

var (
	email    = flag.String("email", "", "EA account email.")
	password = flag.String("password", "", "EA account password.")
	id       = flag.Int("id", 0, "Show stats for given persona ID.")
	game     = flag.String("game", "bf1", "Game: bf1 or bf4")
)

func main() {
	flag.Parse()
	s := battlefield.NewSession(*email, *password)
	if err := s.Login(); err != nil {
		log.Fatal(err)
	}
	stats, err := s.Stats(uint64(*id), "")
	if err != nil {
		log.Fatal(err)
	}
	switch *game {
	case "bf1":
		fmt.Printf("BF1 stats: %s", stats.BF1.String())
	case "bf4":
		fmt.Printf("BF4 stats: %s", stats.BF4.String())
	default:
		log.Fatalf("unknown game: %s", *game)
	}
}
