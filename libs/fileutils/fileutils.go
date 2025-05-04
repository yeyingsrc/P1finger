package fileutils

import (
	"bufio"
	"github.com/P001water/P1finger/libs/p1httputils"
	"os"
)

// LoadFile file content to slice
func ReadLinesFromFile(filename string) (lines []string, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close() //nolint
	s := bufio.NewScanner(f)
	for s.Scan() {
		l := p1httputils.ReplaceAll(s.Text(), "", "\r\n", "\n", "\t", "\v", "\f", "\r", " ")
		if l != "" {
			lines = append(lines, l)
		}

	}
	return
}
