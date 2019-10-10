package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

type client struct {
	who string
	ch  chan<- string
}

const limitTimeout = 5 * time.Minute

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli.ch <- msg
			}

		case cli := <-entering:
			if len(clients) > 0 {
				cli.ch <- "active users: "
				for other := range clients {
					cli.ch <- "\t" + other.who
				}
			}
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli.ch)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)

	input := bufio.NewScanner(conn)

	who := conn.RemoteAddr().String()

	ch <- "Input your name"
	if input.Scan() {
		name := input.Text()
		if len(name) > 0 {
			who = name
		}
	}

	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- client{who, ch}

	timer := time.AfterFunc(limitTimeout, func() {
		conn.Close()
	})

	for input.Scan() {
		messages <- who + ": " + input.Text()
		timer.Reset(limitTimeout)
	}

	leaving <- client{who, ch}
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}
