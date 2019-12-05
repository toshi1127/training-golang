package display

import "testing"

func TestCyclic(t *testing.T) {
	type Cycle struct {
		Value int
		Tail  *Cycle
	}
	var c Cycle
	c = Cycle{42, &c}
	Display("c", c)
}
