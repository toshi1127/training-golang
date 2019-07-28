package main

import (
	"crypto/sha256"
	"fmt"
)

var pc [8]byte

func init() {
	for i := range pc {
		pc[i] = byte(1 << uint(i))
	}
}

func main() {
	h1 := sha256.Sum256([]byte("x"))
	h2 := sha256.Sum256([]byte("X"))
	fmt.Println(diff(h1, h2))
}

func diff(h1, h2 [32]byte) int {
	count := 0
	for i := 0; i < 32; i++ {
		b1 := h1[i]
		b2 := h2[i]
		for j := 0; j < 8; j++ {
			if b1&pc[j] != b2&pc[j] {
				count++
			}
		}
	}
	return count
}