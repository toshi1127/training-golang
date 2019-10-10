package popcount

import (
	"fmt"
	"testing"
)

const x = 2 ^ 64 - 1

func TestPopCount(t *testing.T) {
	var tests = []struct {
		x    uint64
		want int
	}{
		{0, 0},
		{255, 8},
		{256, 1},
		{511, 9},
		{512, 1},
		{18446744073709551615, 64},
	}

	for _, test := range tests {
		descr := fmt.Sprintf("PopCountByLoop(%v)", test.x)
		got := PopCount(test.x)
		if got != test.want {
			t.Errorf("%s = %v, want %v", descr, got, test.want)
		}
	}
}
