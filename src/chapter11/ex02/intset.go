package intset

import (
	"bytes"
	"fmt"
)

// IntSet は負ではない小さな整数のセット
type IntSet struct {
	words []uint64
}

// Has は負ではない値xをセットが含んでいるか否かを返す
func (s *IntSet) Has(x int) bool {
	word, bit := x/64, uint(x%64)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Len はセットの要素数を返す
func (s *IntSet) Len() int {
	count := 0
	for _, w := range s.words {
		count += popCount(w)
	}
	return count
}

func popCount(x uint64) int {
	count := 0
	for x != 0 {
		x = x & (x - 1)
		count++
	}
	return count
}

// Add はセットに負ではない値xを追加する
func (s *IntSet) Add(x int) {
	word, bit := x/64, uint(x%64)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// UnionWith は、sとtの和集合をsに設定する
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// String は"{1, 2, 3}"の形式で文字列としてセットを返す
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", 64*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}