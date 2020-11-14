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
	"io"
	"os"
	"os/signal"
	"syscall"
	"./broadcaster_lib"
)

func sigtermHandler(listener net.Listener) {
	c := make(chan os.Signal)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func(c chan os.Signal, listener net.Listener) {
		<-c
		io.WriteString(os.Stdout, "[INFO] Server exited cleanly\n")

		listener.Close()
		os.Exit(0)
	}(c, listener)
}

func getUsername(conn net.Conn) string {
	var input *bufio.Scanner
	var username string

	fmt.Fprint(conn, "Introduce a username: ")

	input = bufio.NewScanner(conn)
	for input.Scan() {
		username = input.Text()
		if (username != "") {
			break
		}
	}

	return username
}

//!+handleConn
func handleConn(conn net.Conn, broadcast * broadcaster_lib.BroadcastType) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()

	username := getUsername(conn)

	ch <- "You are " + username
	
	// Add the new client to the clients list:

	broadcast.AddClient(who, username, ch)

	broadcast.SendBroadcast(who, "[" + username + "]" + " has arrived")

	// Announce the new list:

	go broadcast.AnnounceClients()

	input := bufio.NewScanner(conn)
	for input.Scan() {
		broadcast.SendBroadcast("", username + ": " + input.Text())
	}
	// NOTE: ignoring potential errors from input.Err()

	conn.Close()

	// Remove the client from teh clients list:

	broadcast.DeleteClient(who)

	go broadcast.SendBroadcast(who, "[" + username + "]" + " has left")
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
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	sigtermHandler(listener)

	// go broadcaster(broadcast)
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
