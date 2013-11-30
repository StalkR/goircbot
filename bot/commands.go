package bot

import (
	"fmt"
	"sync"
)

// Commands holds commands by their name.
type Commands struct {
	sync.Mutex
	cmds map[string]Command
}

// NewStore initializes a new Commands.
func NewCommands() *Commands {
	return &Commands{cmds: make(map[string]Command)}
}

// Add allows plugins to add commands with help.
func (s *Commands) Add(name string, c Command) error {
	s.Lock()
	defer s.Unlock()
	if _, present := s.cmds[name]; present {
		return fmt.Errorf("commands: %s already defined", name)
	}
	s.cmds[name] = c
	return nil
}

// Del allows plugins to remove commands added with Add (no error if not present).
func (s *Commands) Del(name string) {
	s.Lock()
	defer s.Unlock()
	delete(s.cmds, name)
}
