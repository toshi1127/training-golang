package main

import "fmt"

const (
	B  = 1
	KB = B * 1000
	MB = KB * 1000
	GB = MB * 1000
	TB = GB * 1000
	PG = TB * 1000
	EB = PG * 1000
	ZB = EB * 1000
	YB = ZB * 1000
)

func main() {
	fmt.Println(B, KB, MB)
}