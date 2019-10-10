package main

import (
	"bufio"
	"fmt"
)

type WordCounter int
type LineCounter int

func main() {
	p := []byte("世界\n世界\n世界")
	var l LineCounter
	l.Write(p)
	fmt.Printf("%v", l)
}

func (c *WordCounter) Write(p []byte) (int, error) {
	ploc := p
	for i := 0; i < len(ploc); {
		adv, token, _ := bufio.ScanWords(ploc, true)

		if token != nil {
			*c += 1
		}

		ploc = ploc[adv:]
	}
	return 0, nil
}

func (l *LineCounter) Write(p []byte) (int, error) {
	for _, b := range p {
		if b == '\n' {
			*l += 1
		}
	}
	return 0, nil
}