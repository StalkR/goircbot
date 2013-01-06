package goircbot

import (
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"log"
	"sort"
	"strings"
)

// Event is given to plugins.
type Event struct {
	Line   *irc.Line
	Target string
	Args   string
}

// CmdHandler is a callback used by plugins to receive events.
type CmdHandler func(b *Bot, e *Event)

// Command is an IRC command that plugins can add.
type Command struct {
	Help    string
	Handler CmdHandler
	Pub     bool
	Priv    bool
	Hidden  bool
}

func (c *Command) String() string {
	var opts []string
	if c.Pub {
		opts = append(opts, "pub")
	}
	if c.Priv {
		opts = append(opts, "priv")
	}
	if c.Hidden {
		opts = append(opts, "hidden")
	}
	return fmt.Sprintf("%s (%s)", c.Help, strings.Join(opts, "+"))
}

// AddCommand is a convenient way for plugins to add commands with help.
// Use hidden bool to have command hidden from list of all available commands.
func (b *Bot) AddCommand(name string, c Command) {
	if _, present := b.commands[name]; present {
		log.Println("AddCommand: already defined", name)
		return
	}
	b.commands[name] = c
}

// DelCommand is to remove commands added via AddCommand.
func (b *Bot) DelCommand(name string) {
	if _, present := b.commands[name]; !present {
		log.Println("DelCommand: not defined", name)
		return
	}
	delete(b.commands, name)
}

// handleCommand parses a line, matches with registered commands and dispatches.
// A command can be triggered in public with !cmd or bot: cmd and in private with cmd.
func handleCommand(b *Bot, line *irc.Line) {
	words := strings.Split(line.Args[1], " ")
	var pub bool
	var target string
	if strings.HasPrefix(line.Args[0], "#") {
		pub = true
		target = line.Args[0]
		switch {
		case strings.HasPrefix(words[0], "!"):
			words[0] = words[0][1:]
		case words[0] == b.Conn.Me.Nick+":":
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
	for name, cmd := range b.commands {
		if words[0] == name && (pub && cmd.Pub || priv && cmd.Priv) {
			go cmd.Handler(b, &Event{Line: line, Target: target, Args: args})
		}
	}
}

// handleHelp displays all help or help on a single command.
func handleHelp(b *Bot, e *Event) {
	var reply string
	pub := strings.HasPrefix(e.Target, "#")
	priv := !pub
	if e.Args == "" {
		cmds := make([]string, 0, len(b.commands))
		for name, cmd := range b.commands {
			if cmd.Hidden {
				continue
			}
			if pub && cmd.Pub {
				cmds = append(cmds, "!"+name)
			} else if priv && cmd.Priv {
				cmds = append(cmds, name)
			}
		}
		sort.Strings(cmds)
		reply = fmt.Sprintf("Commands: %s", strings.Join(cmds, ", "))
	} else {
		if cmd, present := b.commands[e.Args]; present {
			reply = fmt.Sprintf("%s: %s", e.Args, cmd.String())
		} else {
			reply = "There is no such command."
		}
	}
	b.Conn.Privmsg(e.Target, reply)
}
