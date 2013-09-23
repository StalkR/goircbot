package duration

import (
	"fmt"
	"time"
)

// Format represents a duration in the most appropriate unit by approximation.
// It knows about years, months, days, minutes, hours, seconds.
func Format(d time.Duration) string {
	n, unit := 0, ""
	if d > time.Hour*24*365 {
		n, unit = int(d/time.Hour/24/365), "year"
	} else if d > time.Hour*24*31 {
		n, unit = int(d/time.Hour/24/31), "month"
	} else if d > time.Hour*24 {
		n, unit = int(d/time.Hour/24), "day"
	} else {
		return fmt.Sprintf("%v", d/time.Second*time.Second)
	}
	if n > 1 {
		unit += "s"
	}
	return fmt.Sprintf("%v %v", n, unit)
}
