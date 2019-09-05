package main

import (
	"fmt"
	"os"
	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	for _, list := range Visit(nil, doc) {
		if list != "" {
			fmt.Fprintf(os.Stderr, "error: %v", err)
		}
	}
}

var targetAttributes = map[string]string{
	"a":      "href",
	"script": "src",
	"img":    "src",
	"link":   "href",
}

func Visit(links []string, n *html.Node) []string {
	if n == nil {
		return links
	}

	if n.Type == html.ElementNode {
		a := getAttribute(n, targetAttributes[n.Data])
		if a != "" {
			links = append(links, a)
		}
	}

	links = Visit(links, n.FirstChild)
	return Visit(links, n.NextSibling)
}

func getAttribute(n *html.Node, attr string) string {
	for _, a := range n.Attr {
		if a.Key == attr {
			return a.Val
		}
	}
	return ""
}