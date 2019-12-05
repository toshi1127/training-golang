package display

import "testing"

func TestDisplay(t *testing.T) {

	// struct
	key := struct {
		key string
	}{
		key: "hoge",
	}
	type MS map[struct{ key string }]string
	nma := make(MS)
	nma[key] = "poge"
	Display("struct", nma)

	// map
	m := map[struct{ x int }]int{{1}: 2, {2}: 3}
	Display("map", m)

	// slice
	s := []string{"slice1", "slice2"}
	Display("slice", s)

}
