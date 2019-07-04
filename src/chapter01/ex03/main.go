package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)


// 非効率な実装
func echo1() {
	s, sep := "", ""
	for _, arg := range os.Args[1:] {
		s += sep + arg
		sep = " "
	}
	fmt.Println(s)
}

// strings.Joinを使った実装
func echo2() {
	fmt.Println(strings.Join(os.Args[1:], " "))
}

func main() {
	// echo1の計測
	start1 := time.Now()
	echo1()
	elapsed1 := time.Since(start1).Seconds()

	// echo2の計測
	start2 := time.Now()
	echo2()
	elapsed2 := time.Since(start2).Seconds()

	// 結果の出力
	fmt.Println("echo1", elapsed1)
	fmt.Println("echo2", elapsed2)
}
