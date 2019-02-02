// Package size implements human representation of bytes size.
package size

import "fmt"

// A Byte represents a number of bytes.
type Byte uint64

const (
	_       = iota // ignore first value by assigning to blank identifier
	KB Byte = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
)

// String represents a Byte with the most appropriate unit.
func (b Byte) String() string {
	switch {
	case b >= EB:
		return fmt.Sprintf("%vEB", b.fmt(EB))
	case b >= PB:
		return fmt.Sprintf("%vPB", b.fmt(PB))
	case b >= TB:
		return fmt.Sprintf("%vTB", b.fmt(TB))
	case b >= GB:
		return fmt.Sprintf("%vGB", b.fmt(GB))
	case b >= MB:
		return fmt.Sprintf("%vMB", b.fmt(MB))
	case b >= KB:
		return fmt.Sprintf("%vKB", b.fmt(KB))
	}
	return fmt.Sprintf("%dB", b)
}

// fmt represents a Byte in a unit with 3 significant digits minimum.
func (b Byte) fmt(unit Byte) string {
	f := "%.0f"
	if b < 10*unit {
		f = "%.2f"
	} else if b < 100*unit {
		f = "%.1f"
	}
	return fmt.Sprintf(f, float64(b)/float64(unit))
}
