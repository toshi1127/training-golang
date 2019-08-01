package main

import "fmt"

func main() {
	str := []string{"a", "a", "a", "b", "b"}
	fmt.Printf("%v\n", unique(str))
}

func unique(str []string) []string {
	i := 0
	for _, s := range str {
		if str[i] == s { //初回はここ、2回目 i=0でsがi=1の時の値
			continue
		}
		i++
		str[i] = s
	}
	return str[:i+1]
}