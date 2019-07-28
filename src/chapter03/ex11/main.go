package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

func main() {
	for i := 1; i < len(os.Args); i++ {
		fmt.Printf("comma: %s\n", comma(os.Args[i]))
	}
}

func comma(s string) string {
	p := strings.Split(s, ".")

	n := len(p[0])
	if n <= 3 {
		return s
	}

	init := n % 3
	if init == 0 {
		init = 3
	}

	var b bytes.Buffer
	b.WriteString(p[0][:init])
	for i := init; (i + 3) <= len(p[0]); i += 3 {
		b.WriteString(",")
		b.WriteString(p[0][i : i+3])
	}

	if len(p[1]) != 0 {
		b.WriteString(".")
		b.WriteString(p[1])
	}
	return b.String()
}