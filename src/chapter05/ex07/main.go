package ex07

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	for _, url := range os.Args[1:] {
		res, err := Outline(url)
		if err != nil {
			fmt.Errorf("error at Outline: %v", err)
		}
		fmt.Printf("%s", res)
	}
}

var writer = new(bytes.Buffer)

func Outline(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	forEachNode(doc, startElement, endElement)

	return writer.String(), nil
}

func forEachNode(n *html.Node, pre, post func(*html.Node)) {
	if pre != nil {
		pre(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}

	if post != nil {
		post(n)
	}
}

var depth int

func startElement(n *html.Node) {
	switch n.Type {
	case html.ElementNode:
		fmt.Fprintf(writer, "%*s<%s", depth*2, "", n.Data)

		for _, a := range n.Attr {
			fmt.Fprintf(writer, "%s", " "+a.Key+"="+"'"+a.Val+"'")
		}

		if n.FirstChild == nil {
			fmt.Fprintf(writer, "/>\n")
		} else {
			fmt.Fprintf(writer, ">\n")
		}
		depth++

	case html.TextNode:
		lines := strings.Split(n.Data, "\n")
		for _, line := range lines {
			if line != "" {
				fmt.Fprintf(writer, "%s\n", strings.Repeat(" ", depth*2)+line)
			}
		}

	case html.CommentNode:
		lines := strings.Split(n.Data, "\n")

		fmt.Fprintf(writer, "%s<!--\n", strings.Repeat(" ", depth*2))

		for _, line := range lines {
			fmt.Fprintf(writer, "%s\n", strings.Repeat(" ", depth*2)+line)
		}

		fmt.Fprintf(writer, "%s-->\n", strings.Repeat(" ", depth*2))
	}
}

func endElement(n *html.Node) {
	if n.Type == html.ElementNode {
		depth--
		if n.FirstChild != nil {
			fmt.Fprintf(writer, "%*s</%s>\n", depth*2, "", n.Data)
		}
	}
}