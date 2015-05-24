// Package metal is a library to represent metal bands.
package metal

import "fmt"

// A Band represents a band search result.
type Band struct {
	Name    string
	Genre   string
	Country string
}

// String formats a band information.
func (b Band) String() string {
	return fmt.Sprintf("%s (%s - %s)", b.Name, b.Country, b.Genre)
}
