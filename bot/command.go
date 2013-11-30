package bot

import (
	"fmt"
	"strings"
)

// Command represents a command that plugins can add.
type Command struct {
	Help    string  // Help string for help command.
	Handler Handler // Handler to call.
	Pub     bool    // Command can be accessed publicly on a channel.
	Priv    bool    // Command can be accessed privately in query.
	Hidden  bool    // Hide command from list of all available commands.
}

// String formats a command with its attributes for display in help.
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
