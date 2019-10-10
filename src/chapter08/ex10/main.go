package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

var done = make(chan struct{})
var crawlTokens = make(chan struct{}, 20)

func main() {
	// 入力があったらキャンセルする
	go func() {
		os.Stdin.Read(make([]byte, 1))
		close(done)
	}()

	worklist := make(chan []string)
	var depth int = 1

	depth++
	go func() { worklist <- os.Args[1:] }()

	seen := make(map[string]bool)
	for ; depth > 0; depth-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				depth++
				go func(link string) {
					worklist <- crawl(link)
				}(link)
			}
		}
	}
}

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

func crawl(url string) []string {
	fmt.Println(url)
	crawlTokens <- struct{}{}
	list, err := extract(url)
	<-crawlTokens

	if err != nil {
		log.Print(err)
	}
	return list
}

func extract(url string) ([]string, error) {
	if cancelled() {
		return nil, fmt.Errorf("cancelled: get %q", url)
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Cancel = done

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue // ignore bad URLs
				}
				links = append(links, link.String())
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return links, nil
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
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