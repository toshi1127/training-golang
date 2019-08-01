package main

import "fmt"

func main() {
	a := [...]int{0, 1, 2, 3, 4, 5, 6, 7}
	reverse(&a)
	fmt.Println(a)
}

func reverse(nums *[8]int) {
	for i, j := 0, len(nums)-1; i < j; i, j = i+1, j-1 {
		nums[i], nums[j] = nums[j], nums[i]
	}
}