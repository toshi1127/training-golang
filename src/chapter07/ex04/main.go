package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args[1:]) == 0 {
		return
	}
	in := os.Args[1]
	r := NewReader(in)
	out := make([]byte, 16)
	size, err := r.Read(out)
	for err == nil {
		fmt.Printf("readed: %s\n", string(out[:size]))
		size, err = r.Read(out)
	}
	if err != io.EOF {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func NewReader(s string) io.Reader {
	return &reader{s: s}
}

type reader struct {
	s string
	i int
}

func (r *reader) Read(p []byte) (int, error) {
	if r.i >= len(r.s) {
		return 0, io.EOF
	}
	size := copy(p, r.s[r.i:])
	r.i += size
	return size, nil
}
