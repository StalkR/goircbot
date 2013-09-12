// Package tor implements a library to query a TOR node via its control port.
package tor

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/StalkR/goircbot/lib/size"
)

// An Information represents various TOR info on the node.
type Information struct {
	Version     string
	In          int
	Out         int
	Fingerprint string
	Flags       []string
	Circuits    int
}

// String formats Information on one line.
func (i *Information) String() string {
	in, out := size.Byte(i.In).String(), size.Byte(i.Out).String()
	url := fmt.Sprintf("https://atlas.torproject.org/#details/%v", i.Fingerprint)
	return fmt.Sprintf("Running %v, %v in, %v out, %v circuits - Flags: %v - %v",
		i.Version, in, out, i.Circuits, strings.Join(i.Flags, ", "), url)
}

// Info obtains information on a TOR node by connecting to its control protocol.
func Info(hostPort, pwd string) (*Information, error) {
	c, err := Connect(hostPort, pwd)
	if err != nil {
		return nil, err
	}
	i := &Information{}

	i.Version, err = getinfo(c, "version")
	if err != nil {
		return nil, err
	}

	i.Fingerprint, err = getinfo(c, "fingerprint")
	if err != nil {
		return nil, err
	}

	i.In, err = getinfoInt(c, "traffic/read")
	if err != nil {
		return nil, err
	}

	i.Out, err = getinfoInt(c, "traffic/written")
	if err != nil {
		return nil, err
	}

	ns, err := getinfoMany(c, fmt.Sprintf("ns/id/%s", i.Fingerprint))
	if err != nil {
		return nil, err
	}
	for _, line := range ns {
		if strings.HasPrefix(line, "s ") {
			i.Flags = strings.Split(line[2:], " ")
			break
		}
	}

	circuitStatus, err := getinfoMany(c, "circuit-status")
	if err != nil {
		return nil, err
	}
	i.Circuits = len(circuitStatus)

	return i, nil
}

// Connect creates a new authenticated connection to TOR control port.
func Connect(hostPort, pwd string) (c net.Conn, err error) {
	c, err = net.DialTimeout("tcp", hostPort, 2*time.Second)
	if err != nil {
		return nil, errors.New("tor: down")
	}
	if err := send(c, fmt.Sprintf(`authenticate "%s"`, pwd)); err != nil {
		return nil, err
	}
	line, err := read(c)
	if err != nil {
		return nil, err
	}
	if line != "250 OK" {
		return nil, errors.New("tor: auth failed")
	}
	return c, nil
}

// send sends a command.
func send(c net.Conn, cmd string) error {
	_, err := fmt.Fprintf(c, cmd+"\r\n")
	return err
}

// read reads a line.
func read(c net.Conn) (string, error) {
	return readUntil(c, '\n')
}

// readUntil reads one byte at a time until stop character.
func readUntil(c net.Conn, until byte) (string, error) {
	var bytes []byte
	char := make([]byte, 1)
	for {
		n, err := c.Read(char)
		if n != 1 || err != nil {
			return "", fmt.Errorf("tor: error reading %#v: %v", string(bytes), err)
		}
		bytes = append(bytes, char[0])
		if char[0] == until {
			break
		}
	}
	return strings.TrimSpace(string(bytes)), nil
}

// getinfo sends a getinfo command for a key and returns its value.
func getinfo(c net.Conn, key string) (string, error) {
	if err := send(c, "getinfo "+key); err != nil {
		return "", err
	}
	line, err := read(c)
	if err != nil {
		return "", err
	}
	if line[:4] != "250-" {
		return "", fmt.Errorf("tor: code error: %#v", line)
	}
	line = line[4:]
	if line[:len(key)+1] != key+"=" {
		return "", fmt.Errorf("tor: key error: %#v", line)
	}
	value := line[len(key)+1:]
	line, err = read(c)
	if err != nil {
		return "", err
	}
	if line != "250 OK" {
		return "", fmt.Errorf("tor: status error: %#v", line)
	}
	return value, nil
}

// getinfo sends a getinfo command for a key and returns its integer value.
func getinfoInt(c net.Conn, key string) (int, error) {
	value, err := getinfo(c, key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(value)
}

// getinfoMany sends a getinfo command for a key with many output lines.
func getinfoMany(c net.Conn, key string) ([]string, error) {
	if err := send(c, "getinfo "+key); err != nil {
		return nil, err
	}
	line, err := read(c)
	if err != nil {
		return nil, err
	}
	if line[:4] != "250+" {
		return nil, fmt.Errorf("tor: code error: %#v", line)
	}
	line = line[4:]
	if line[:len(key)+1] != key+"=" {
		return nil, fmt.Errorf("tor: key error: %#v", line)
	}
	var lines []string
	for {
		line, err = read(c)
		if err != nil {
			return nil, err
		}
		if line == "." {
			break
		}
		lines = append(lines, line)
	}
	line, err = read(c)
	if err != nil {
		return nil, err
	}
	if line != "250 OK" {
		return nil, fmt.Errorf("tor: status error: %#v", line)
	}
	return lines, nil
}
