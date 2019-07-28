package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var width = flag.Int("w", 256, "hash width (256, 384 or 512)")

func main() {
	flag.Parse()

	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	switch *width {
	case 256:
		h := sha256.Sum256(b)
		print(h[:])
	case 384:
		h := sha512.Sum384(b)
		print(h[:])
	case 512:
		h := sha512.Sum512(b)
		print(h[:])
	default:
		log.Fatal("invalid hash widh specified.")
	}
}

func print(hash []byte) {
	for _, v := range hash {
		fmt.Printf("%X", v)
	}
	fmt.Println()
}