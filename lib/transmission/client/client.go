// Transmission client to print stats.
package main

import (
	"fmt"
	"os"

	"github.com/StalkR/goircbot/lib/transmission"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <transmission url>\n", os.Args[0])
		os.Exit(1)
	}
	c, err := transmission.New(os.Args[1])
	if err != nil {
		fmt.Println("err", err)
	}
	s, err := c.Stats()
	if err != nil {
		fmt.Println("err", err)
	}
	fmt.Printf("%v KB/s DL, %v KB/s UL, %v torrents (%v active, %v paused)\n",
		s.DownloadSpeed/1024, s.UploadSpeed/1024, s.TorrentCount,
		s.ActiveTorrentCount, s.PausedTorrentcount)
}
