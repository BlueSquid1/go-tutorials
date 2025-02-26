package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type client struct {
	name    string
	channel chan<- string
}

var (
	// this is a channel of read only channels
	entering = make(chan client)
	leaving  = make(chan client)
	// Channels are thread safe
	messages = make(chan string)
)

func broadcaster() {
	// Only a map so it's easy to delete from.
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			// Broadcast to all clients.
			for cli := range clients {
				cli.channel <- msg
			}
		case newCli := <-entering:
			// send currently connected clients to the newly connected peer
			newCli.channel <- "currently connected peers:"
			for cli := range clients {
				newCli.channel <- cli.name
			}

			// store a read+write channel as a write only channel here
			clients[newCli] = true

		case cli := <-leaving:
			// Can do comparison with channels
			delete(clients, cli)
			close(cli.channel)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()

	clientDetails := client{
		name:    who,
		channel: ch,
	}

	ch <- "You are " + who
	messages <- who + " has arrived"

	entering <- clientDetails

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + " : " + input.Text()
	}
	leaving <- clientDetails
	messages <- who + " has left"
	// net.Conn is thread safe so it's fine to close the connection in this goroutine
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
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
