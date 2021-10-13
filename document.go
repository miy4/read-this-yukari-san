package yonde

import (
	"bufio"
	"io"
)

type document struct {
	lines []string
}

func open(r io.Reader) (*document, error) {
	s := bufio.NewScanner(r)
	ln := make([]string, 0)
	for s.Scan() {
		ln = append(ln, s.Text())
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return &document{lines: ln}, nil
}

func (doc document) len() int {
	return len(doc.lines)
}
