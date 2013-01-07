// Package tail implements a plugin to tail files and notify new lines on all channels.
package tail

import (
	"bufio"
	bot "goircbot"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func tail(b *bot.Bot, path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Println("tail: os.Open", err)
		return
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		if _, err := r.ReadString('\n'); err != nil {
			break
		}
	}
	var line string
	for {
		time.Sleep(time.Duration(1) * time.Second)
		for {
			buf, err := r.ReadString('\n')
			line += buf
			if err == io.EOF {
				break
			} else if err == nil {
				notify(b, strings.TrimSpace(line))
				line = ""
			}
		}
		f.Seek(0, os.SEEK_END) // catch when file is truncated
	}
}

func notify(b *bot.Bot, line string) {
	if !b.Conn.Connected {
		return
	}
	for _, channel := range b.Conn.Me.Channels() {
		b.Conn.Privmsg(channel.Name, line)
	}
}

// Register registers the plugin with a bot.
func Register(b *bot.Bot, paths []string) {
	for _, path := range paths {
		go tail(b, path)
	}
}
