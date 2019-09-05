package ex17

import (
	"golang.org/x/net/html"
)

func ElementsByTagName(doc *html.Node, names ...string) []*html.Node {
	res := []*html.Node{}

	if doc.Type == html.ElementNode {
		for _, name := range names {
			if doc.Data == name {
				res = append(res, doc)
			}
		}
	}

	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		for _, node := range ElementsByTagName(c, names...) {
			res = append(res, node)
		}
	}

	return res
}