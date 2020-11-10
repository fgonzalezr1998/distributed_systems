// Copyright © 2020
// License: APACHE
// Author: Fernando González <fergonzaramos@yahoo.es>

/*
 * Chat is a server that lets clients chat with each other
 */

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"./broadcaster_lib"
)

func broadcaster(broadcast * broadcaster_lib.BroadcastType) {
	clients := make(map[broadcaster_lib.ClientChannelType]bool) // all connected clients
	for {
		select {
		case msg := <-broadcast.Messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli <- msg
			}

		case cli := <-broadcast.Entering:
			clients[cli] = true

		case cli := <-broadcast.Leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

//!-broadcaster

//!+handleConn
func handleConn(conn net.Conn, broadcast * broadcaster_lib.BroadcastType) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	broadcast.Messages <- who + " has arrived"
	broadcast.Entering <- ch

	broadcast.AddClient(who)
	broadcast.PrintConnectedClients()

	input := bufio.NewScanner(conn)
	for input.Scan() {
		broadcast.Messages <- who + ": " + input.Text()
	}
	// NOTE: ignoring potential errors from input.Err()

	broadcast.Leaving <- ch
	broadcast.Messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

//!-handleConn

//!+main
func main() {
	var broadcast * broadcaster_lib.BroadcastType = new(broadcaster_lib.BroadcastType)
	broadcast.Init()
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster(broadcast)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn, broadcast)
	}
}

//!-main
