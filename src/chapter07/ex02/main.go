package ex02

import (
	"io"
)

type wrapper struct {
	w     io.Writer
	count int64
}

func CountingWriter(w io.Writer) (io.Writer, *int64) {
	var wrap = wrapper{w, 0}
	return &wrap, &wrap.count
}

func (w *wrapper) Write(p []byte) (int, error) {
	n, err := w.w.Write(p)
	w.count += int64(n)
	return n, err
}