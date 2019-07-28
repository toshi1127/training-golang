// 二つの文字列がお互いにアナグラムになっているか、すなわち同じ文字を異なる順番で含んでいるかを報告する関数を書きなさい
package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	p := "subject: %s\ncandidate%s\nisAnagram:%v\n"

	sub := "あとうかい"
	can := "かとうあい"

	fmt.Printf(p, sub, can, isAnagram(sub, can))

	sub = "ないようがいい"
	can = "いいようがない"

	fmt.Printf(p, sub, can, isAnagram(sub, can))

	sub = "もうねむい"
	can = "もうあさだ"

	fmt.Printf(p, sub, can, isAnagram(sub, can))
}

func isAnagram(sub, can string) bool {
	if sub == can {
		return false
	}

	return alphagram(sub) == alphagram(can)
}

func alphagram(s string) string { // 配列の中身をソート(アナグラムなら等しくなる)して、繋げる
	chars := strings.Split(strings.ToLower(s), "") // ToLower: 文字列sをUnicodeの小文字にマッピングしたコピーを返す。
	sort.Strings(chars)
	return strings.Join(chars, "")
}