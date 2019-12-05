package intset

import "testing"

func (s *IntSet) equals(m map[int]bool) bool {
	if s.Len() != len(m) {
		return false
	}
	for v, _ := range m {
		if !s.Has(v) {
			return false
		}
	}
	return true
}

func TestAdd(t *testing.T) {
	s := &IntSet{}
	s.Add(1)
	s.Add(10)
	s.Add(100)
	s.Add(1000)
	want := map[int]bool{
		1: true, 10: true, 100: true, 1000: true,
	}
	if !s.equals(want) {
		t.Errorf("=> %v != %v", s, want)
	}
}
