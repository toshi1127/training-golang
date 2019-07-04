package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// io.ReadAllの代わりにio.Copyを使う

func main() {
	for _, url := range os.Args[1:] {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}

		// srcをresp.Body、dstをos.Stdoutにしてio.Copyを呼ぶ
		_, err = io.Copy(os.Stdout, resp.Body)

		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
			os.Exit(1)
		}
		resp.Body.Close()
	}
}