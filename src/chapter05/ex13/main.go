package ex13

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/nwidger/gopl.io/ch5/links"
)

var host string
var paths []string

const depth = 3

func main() {
	host = "https://golang.org"
	breadthFirst(crawl, []string{host})
}

func crawl(url string) []string {

	if strings.HasPrefix(url, host) {
		path := strings.Split(url, host)[1]

		if strings.HasSuffix(url, "/") {
			paths = append(paths, path)
			os.MkdirAll("."+path, 0777)
		} else {
			res, _ := http.Get(url)
			f, _ := os.Create("." + path)

			defer f.Close()
			defer res.Body.Close()

			body, _ := ioutil.ReadAll(res.Body)

			f.Write(body)
		}
	}

	fmt.Println("%s", paths)
	if len(paths) > depth {
		return nil
	}

	list, _ := links.Extract(url)

	return list
}

func breadthFirst(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true //見たサイトはtrue
				worklist = append(worklist, f(item)...)
			}
		}
	}
}
