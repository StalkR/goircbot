// Package size implements human representation of bytes size.
package size

import "fmt"

// A Byte represents a number of bytes.
type Byte int64

const (
	_       = iota // ignore first value by assigning to blank identifier
	KB Byte = 1 << (10 * iota)
	MB
	GB
	TB
	PB
)

func (b Byte) String() string {
	switch {
	case b >= PB:
		return fmt.Sprintf("%dPB", b/PB)
	case b >= TB:
		return fmt.Sprintf("%dTB", b/TB)
	case b >= GB:
		return fmt.Sprintf("%dGB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%dMB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%dKB", b/KB)
	}
	return fmt.Sprintf("%dB", b)
}
