package main

import (
	"fmt"
	"os"
	"strings"
	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	for _, text := range ShowTexts(nil, doc) {
		if text != "" {
			fmt.Fprintf(os.Stderr, "error: %v", err)
		}
	}
}

func ShowTexts(list []string, n *html.Node) []string {
	if n == nil {
		return list
	}

	if n.Type == html.TextNode {
		s := strings.TrimSpace(n.Data)
		list = append(list, s)
	}

	if !(n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style")) {
		list = ShowTexts(list, n.FirstChild)
	}

	return ShowTexts(list, n.NextSibling)
}