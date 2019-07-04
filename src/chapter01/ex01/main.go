package main

import (
	"fmt"
	"os"
	"strings"
)

// strings.Joinにos.Argsを丸っと渡して、から文字でつなげる

func main() {
	fmt.Println(strings.Join(os.Args, " "))
}
