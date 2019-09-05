package main

import (
	"bytes"
	"fmt"
	"io"
)

func main() {
	in := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	r := bytes.NewReader(in)
	limited := LimitReader(r, 5)
	out := make([]byte, 10)
	limited.Read(out)
	fmt.Printf("Input=%v\n", in)
	fmt.Printf("Output(limited 5)=%v\n", out)
}

func LimitReader(r io.Reader, n int) io.Reader {
	return &limitReader{r: r, n: n}
}

type limitReader struct {
	r io.Reader
	i int
	n int
}

func (l *limitReader) Read(b []byte) (int, error) {
	rem := l.n - l.i
	if rem <= 0 {
		return 0, io.EOF
	}
	if rem > len(b) {
		size, err := l.r.Read(b)
		l.i += size
		return size, err
	}
	buf := make([]byte, rem)
	size, err := l.r.Read(buf)
	copy(b, buf[:size])
	l.i += size
	return size, err
}
