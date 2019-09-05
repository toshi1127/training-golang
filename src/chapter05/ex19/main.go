package main

import (
	"fmt"
)

func main() {
	panicTest()
}

// なんでdefer func呼び出されないんだろう・・・
func panicTest() (res string, err error) {
	defer func() {
		fmt.Errorf("panic occured, and deferred func was called\n")
		if p := recover(); p != nil {
			err = fmt.Errorf("internal error: %v", p)
		}
	}()
	panic("panic!")
	return "hoge", err
}
