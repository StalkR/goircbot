package errors

import (
	"bufio"
	"os"
	"strings"
)

func parse(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines [][]string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		r := strings.Split(scanner.Text(), "\t")
		if len(r) != 3 {
			continue
		}
		lines = append(lines, r)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
