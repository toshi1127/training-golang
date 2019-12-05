package main

import (
	"net"
)

func main() {
    conn, _ := net.Dial("tcp", "localhost:8888")
  // connをread/write
}