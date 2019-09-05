package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	input := os.Args[1]
	url, e := ParseAndCount(input)
	if e != nil {
		fmt.Errorf("failed counting from %s: %v", url, e)
	}
}

func ParseAndCount(url string) (string, error) {
	r, e := http.Get(url)
	if e != nil {
		return url, e
	}

	doc, e := html.Parse(r.Body)
	if e != nil {
		return url, e
	}

	counts := CountWordsAndImages(map[string]int{"words": 0, "images": 0}, doc)

	for element, count := range counts {
		fmt.Fprintf(writer, "%s: %d\n", element, count)
	}
	return url, nil
}

func CountWordsAndImages(counts map[string]int, n *html.Node) map[string]int {
	if n == nil {
		return counts
	}

	if n.Type == html.TextNode {
		s := bufio.NewScanner(strings.NewReader(n.Data))
		s.Split(bufio.ScanWords)

		for s.Scan() {
			counts["words"]++
		}
	}

	if n.Type == html.ElementNode && n.Data == "img" {
		counts["images"]++
	}

	if !(n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style")) {
		counts = CountWordsAndImages(counts, n.FirstChild)
	}

	return CountWordsAndImages(counts, n.NextSibling)
}