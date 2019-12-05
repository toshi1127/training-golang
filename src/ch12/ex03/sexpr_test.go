package sexpr

import (
	"testing"
)

func TestMarshal(t *testing.T) {
	ts := []struct {
		i        interface{}
		expected string
	}{
		{true, "t"},
		{false, "nil"},
		{float64(10.3), "10.3"},
		{complex(10, 3), "#C(10 3)"},
		{complex(-300, -300), "#C(-300 -300)"},
	}

	for _, k := range ts {
		actual, err := Marshal(k.i)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}
		if k.expected != string(actual) {
			t.Fatalf("exptected %v, but actual %v", k.expected, actual)
		}
	}
}
