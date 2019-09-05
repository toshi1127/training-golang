package main

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func main() {
	for _, url := range os.Args[1:] {
		res, err := http.Get(url)

		if err != nil {
			fmt.Errorf("error at Outline: %v", err)
		}

		doc, err := html.Parse(res.Body)

		if err != nil {
			fmt.Errorf("error at Outline: %v", err)
		}

		n := ElementByID(doc, "lowframe")

		if n != nil {
			fmt.Printf("<%s", n.Data)
			for _, a := range n.Attr {
				fmt.Printf("%s=%s ", a.Key, a.Val)
			}
			fmt.Println(">")
		}
	}
}

func ElementByID(doc *html.Node, id string) *html.Node {
	var node *html.Node

	forEachNode(doc, func(n *html.Node) bool {
		if n.Type == html.ElementNode {
			for _, a := range n.Attr {
				if a.Key == "id" && a.Val == id {
					node = n
					return false
				}
			}
		}
		return true
	}, nil)
	return node
}

func forEachNode(n *html.Node, pre, post func(*html.Node) bool) bool {
	if pre != nil {
		if !pre(n) {
			return false
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if !forEachNode(c, pre, post) {
			return false
		}
	}

	if post != nil {
		_ = post(n) /*要らないけど問題の指定上残している*/
	}
	return true
}