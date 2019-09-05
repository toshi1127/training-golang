package ex05

import (
	"bytes"
	"fmt"
)

type IntSet struct {
	words []uint
}

func (s *IntSet) AddAll(numbers ...int) {
	for _, v := range numbers {
		s.Add(v)
	}
}

func (s *IntSet) Len() int {
	return len(s.words)
}

const bitSize = 32 << (^uint(0) >> 63)

func (s *IntSet) Remove(x int) {
	word, bit := x/bitSize, uint(x%bitSize)
	s.words[word] = s.words[word] & ^(1 << bit)
}

func (s *IntSet) Clear() {
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		s.words[i] = word & 0
	}
}

func (s *IntSet) Elems() []int {
	elems := make([]int, 0, s.Len())

	for i, sword := range s.words {
		for n := 0; n < bitSize; n++ {
			if sword&(1<<uint(n)) != 0 {
				elems = append(elems, i*bitSize+n)
			}
		}
	}
	return elems
}

func (s IntSet) Copy() IntSet {
	return s
}

func (s *IntSet) Has(x int) bool {
	word, bit := x/bitSize, uint(x%bitSize)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

func (s *IntSet) Add(x int) {
	word, bit := x/bitSize, uint(x%bitSize)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

func (s *IntSet) IntersectWith(t *IntSet) {
	for i, _ := range s.words {
		if i >= len(t.words) {
			s.words[i] = 0
			continue
		}
		s.words[i] &= t.words[i]
	}
}

func (s *IntSet) DifferenceWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] &^= tword
		}
	}
}

func (s *IntSet) SymmetricDifferenceWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] ^= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < bitSize; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", bitSize*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}
