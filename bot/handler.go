package bot

import (
	"strings"

	"github.com/fluffle/goirc/client"
)

// Handler is a callback used by plugins to receive command events.
type Handler func(e *Event)

// Event is given to a command's Handler when a command matched by name.
type Event struct {
	Bot    Bot
	Target string // Channel if public, nick if private.
	Args   string // Command args (excluding command name).
	Line   *client.Line
}

// Handle parses a line, matches with commands and dispatches to handler.
// A command can be triggered
//   - in public on a channel with: /msg #chan !cmd [args]
//     or directly at a bot: /msg #chan bot: cmd [args]
//   - in private in a query: /msg bot cmd [args]
func (s *Commands) Handle(b Bot, line *client.Line) {
	words := strings.Split(line.Args[1], " ")
	var pub bool
	var target string
	if strings.HasPrefix(line.Args[0], "#") {
		pub = true
		target = line.Args[0]
		switch {
		case strings.HasPrefix(words[0], "!"):
			words[0] = words[0][1:]
		case words[0] == b.Me().Nick+":":
			words = words[1:]
			if len(words) <= 0 {
				return
			}
		default:
			return // Not a command.
		}
	} else {
		pub = false
		target = line.Nick
	}
	priv := !pub
	args := strings.Join(words[1:], " ")
	s.Lock()
	defer s.Unlock()
	for name, c := range s.cmds {
		if words[0] == name && (pub && c.Pub || priv && c.Priv) {
			go c.Handler(&Event{Bot: b, Line: line, Target: target, Args: args})
		}
	}
}
