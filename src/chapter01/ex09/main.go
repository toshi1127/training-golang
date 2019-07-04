package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// HTTPのステータスコードを表示するようにする

func main() {
	for _, url := range os.Args[1:] {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}

		// resp.Bosyを出力してHTTPのステータスコードを表示する
		fmt.Println("status:", resp.Status)

		b, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
			os.Exit(1)
		}
		fmt.Printf("%s", b)
	}
}

