package yonde

import (
	"bufio"
	"fmt"
	"io"
)

type document struct {
	lines  []string
	chunks []*chunk
}

type chunk struct {
	s int
	e int
}

func (chunk chunk) String() string {
	return fmt.Sprintf("{s: %d, e: %d}", chunk.s, chunk.e)
}

func (chunk chunk) containsAt(at int) bool {
	return chunk.s <= at && at <= chunk.e
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

	chunks := make([]*chunk, 0)
	var c *chunk
	for i, l := range ln {
		if len(l) <= 0 {
			if c != nil {
				chunks = append(chunks, c)
				c = nil
			}
		} else {
			if c == nil {
				c = &chunk{s: i, e: i}
			} else {
				c.e = i
			}
		}
	}
	if c != nil {
		chunks = append(chunks, c)
	}

	return &document{lines: ln, chunks: chunks}, nil
}

func (doc document) len() int {
	return len(doc.lines)
}
