package ex15

import "fmt"

func Max(first int, vals ...int) int {
	max := first
	for _, val := range vals {
		if val > max {
			max = val
		}
	}
	return max
}

func Min(first int, vals ...int) int {
	min := first
	for _, val := range vals {
		if val < min {
			min = val
		}
	}
	return min
}

func MaxRequireArgs(vals ...int) (int, error) {
	if len(vals) == 0 {
		return 0, fmt.Errorf("need at least one argument!")
	}
	return Max(vals[0], vals[1:]...), nil
}