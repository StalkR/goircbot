package bot

import (
	"fmt"
	"sort"
	"strings"
)

// Help displays all commands available or help on a single command.
func (s *Commands) Help(e *Event) {
	s.Lock()
	defer s.Unlock()
	var reply string
	pub := strings.HasPrefix(e.Target, "#")
	priv := !pub
	if e.Args == "" {
		cmds := make([]string, 0, len(s.cmds))
		for name, cmd := range s.cmds {
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
		if cmd, present := s.cmds[e.Args]; present {
			reply = fmt.Sprintf("%s: %s", e.Args, cmd.String())
		} else {
			reply = "There is no such command."
		}
	}
	e.Bot.Privmsg(e.Target, reply)
}
