package ex09

import (
	"regexp"
)

func Expand(s string, f func(string) string) string {
	re, _ := regexp.Compile(`\$\w+`)
	return re.ReplaceAllStringFunc(s, f)
}