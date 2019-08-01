package main

import (
	"fmt"
	"unicode"
)

func main() {
	b := []byte("ab   de  f g")
	fmt.Printf("%v\n", string(uniqueSpace(b)))
}

func uniqueSpace(b []byte) []byte {
	res := []byte{}
	for i, c := range b {
		if unicode.IsSpace(rune(c)) {
			if i > 0 && unicode.IsSpace(rune(b[i-1])) {
				continue
			} else {
				res = append(res, ' ')
			}
		} else {
			res = append(res, c)
		}
	}
	return res
}