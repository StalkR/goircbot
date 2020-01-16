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
//     works also with a comma: /msg #chan bot, cmd [args]
//   - in private in a query: /msg bot cmd [args]
func (s *Commands) Handle(b Bot, line *client.Line) {
	words := strings.Split(line.Args[1], " ")
	var direct, indirect, private bool
	var target string

	if strings.HasPrefix(line.Args[0], "#") {
		target = line.Args[0]
		switch {
		case strings.HasPrefix(words[0], b.Prefix()):
			indirect = true
			words[0] = words[0][1:]

		case words[0] == b.Me().Nick+":" || words[0] == b.Me().Nick+",":
			direct = true
			words = words[1:]
			if len(words) < 1 {
				return
			}

		default:
			return // Not a command.
		}
	} else {
		private = true
		target = line.Nick
	}

	args := strings.Join(words[1:], " ")

	s.Lock()
	defer s.Unlock()

	c, present := s.cmds[words[0]]
	if !present {
		return
	}
	direct = direct && c.Pub
	indirect = indirect && c.Pub && !c.NoExclamation
	private = private && c.Priv
	if direct || indirect || private {
		go c.Handler(&Event{Bot: b, Line: line, Target: target, Args: args})
	}
}
