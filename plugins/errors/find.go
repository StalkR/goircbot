package errors

import (
	"fmt"
	"strconv"
	"strings"
)

func find(table [][]string, arg string) string {
	if r := findName(table, arg); r != nil {
		return fmt.Sprintf("%s (%s): %s", r[0], r[1], r[2])
	}
	code, err := atoi(arg)
	if err != nil {
		return fmt.Sprintf("%s: not found", arg)
	}
	r := findCode(table, code)
	if r == nil {
		return fmt.Sprintf("%s: not found", arg)
	}
	return fmt.Sprintf("%s (%s): %s", r[0], r[1], r[2])
}

func findName(table [][]string, name string) []string {
	for _, r := range table {
		if name == r[1] {
			return r
		}
	}
	return nil
}

func findCode(table [][]string, code uint64) []string {
	for _, r := range table {
		if i, err := atoi(r[0]); err == nil && code == i {
			return r
		}
	}
	return nil
}

func atoi(s string) (uint64, error) {
	base := 10
	if strings.HasPrefix(s, "0x") {
		s = s[2:]
		base = 16
	} else if strings.HasPrefix(s, "h") {
		s = s[2:]
		base = 16
	}
	return strconv.ParseUint(s, base, 64)
}
