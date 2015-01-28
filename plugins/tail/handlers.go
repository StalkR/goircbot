// Package tail implements a plugin to tail files and notify new lines on all channels.
package tail

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/StalkR/goircbot/bot"
)

// Tail watches a file and calls notify for every new line.
// File can be truncated but it does not detect if file is renamed.
func Tail(path string, notify func(string)) {
	var f *os.File
	for ; ; time.Sleep(time.Minute) {
		fp, err := os.Open(path)
		if err == nil {
			f = fp
			break
		}
		log.Println("tail: error open", path, err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		if _, err := r.ReadString('\n'); err != nil {
			break
		}
	}
	var line string
	for ; ; time.Sleep(time.Second) {
		for {
			buf, err := r.ReadString('\n')
			line += buf
			if err == io.EOF {
				break
			} else if err == nil {
				notify(strings.TrimSpace(line))
				line = ""
			}
		}
		f.Seek(0, os.SEEK_END) // catch when file is truncated
	}
}

func notify(b bot.Bot, line string) {
	if !b.Connected() {
		return
	}
	for _, channel := range b.Channels() {
		b.Privmsg(channel, line)
	}
}

// Register registers the plugin with a bot.
func Register(b bot.Bot, paths []string) {
	for _, path := range paths {
		go Tail(path, func(line string) { notify(b, line) })
	}
}
