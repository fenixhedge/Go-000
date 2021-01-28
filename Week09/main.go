package main

import (
	"bufio"
	"log"
	"net"
	"fmt"
	"sync"
)

var (
	ID    int
	lock  sync.Mutex
)

func GenerateId() int {
	lock.Lock()
	defer lock.Unlock()

	ID++
	return ID
}

func sendMessage(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	id := GenerateId()
	messageCh := make(chan string, 8)

	go sendMessage(conn, messageCh)

	input := bufio.NewScanner(conn)
	for input.Scan() {
		fmt.Println(id, input.Text())
		messageCh <- input.Text()
	}

	if err := input.Err(); err != nil {
		log.Println("read failure:", err)
	}
}

func main() {
	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}

	log.Println("start tcp listen")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("listener.Accept error{%v}")
			continue
		}
		go handleConn(conn)
	}
}