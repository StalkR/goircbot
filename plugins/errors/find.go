package errors

import (
	"fmt"
	"strconv"
	"strings"
)

type info struct {
	id   uint32
	name string
	doc  string
}

func (i info) String() string {
	return fmt.Sprintf("%d/0x%08x (%s): %s", i.id, i.id, i.name, i.doc)
}

func find(table []info, arg string) (info, error) {
	if r, err := findName(table, arg); err == nil {
		return r, nil
	}
	code, err := atoi(arg)
	if err != nil {
		return info{}, fmt.Errorf("%s: not found", arg)
	}
	r, err := findCode(table, code)
	if err != nil {
		return info{}, err
	}
	return r, nil
}

func findName(table []info, name string) (info, error) {
	for _, r := range table {
		if name == r.name {
			return r, nil
		}
	}
	return info{}, fmt.Errorf("%s: not found", name)
}

func findCode(table []info, code uint32) (info, error) {
	for _, r := range table {
		if code == r.id {
			return r, nil
		}
	}
	return info{}, fmt.Errorf("%s: not found", code)
}

func atoi(s string) (uint32, error) {
	base := 10
	if strings.HasPrefix(s, "0x") {
		s = s[2:]
		base = 16
	} else if strings.HasPrefix(s, "h") {
		s = s[2:]
		base = 16
	}
	r, err := strconv.ParseUint(s, base, 32)
	if err != nil {
		return 0, err
	}
	return uint32(r), nil
}
