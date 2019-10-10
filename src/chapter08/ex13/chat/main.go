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
	ch  chan<- string // 送信用メッセージチャネル
}

const limitTimeout = 5 * time.Minute

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // クライアントからのメッセージ
)

func broadcaster() {
	clients := make(map[client]bool) // 接続中のクライアント
	for {
		select {
		case msg := <-messages:
			// 受信したメッセージを全てのクライアントの送信用チャネルへブロードキャストする
			for cli := range clients {
				cli.ch <- msg
			}

		case cli := <-entering:
			//.現在のクライアントの集まりを通知
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

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- client{who, ch}

	timer := time.AfterFunc(limitTimeout, func() {
		conn.Close()
	})

	input := bufio.NewScanner(conn)
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