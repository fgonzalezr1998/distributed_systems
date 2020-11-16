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
	"strings"
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

func getUsername(conn net.Conn, b * broadcaster_lib.BroadcastType) string {
	var input *bufio.Scanner
	var username string

	fmt.Fprint(conn, "[SERVER] Introduce a username: ")

	input = bufio.NewScanner(conn)
	for input.Scan() {
		username = input.Text()
		if (username != "" && !b.Exists(username)) {
			break
		} else {
			fmt.Fprint(conn, "[SERVER] Username invalid or already in use. Try other: ")
		}
	}

	return username
}

func privateChanRequested(tags []string) bool {
	return len(tags) == 2 && tags[0] == "priv" 
}

func exitPrivChanRequested(m string) bool {
	return m == "end priv"
}

func sendMsg(c chan string, b * broadcaster_lib.BroadcastType, id, username, msg string) {
	var tags []string
	var output_user string

	tags = strings.Split(msg, " ")

	if (privateChanRequested(tags)) {
		b.SetPrivateChan(id, tags[1])
		c <- "[SERVER] To exit the private channel, type: 'end priv'"
		b.SendBroadcast(id, "[SERVER] " + username + " is in a private channel")
		return
	}

	if (exitPrivChanRequested(msg)) {
		b.SetPrivateChan(id, "")
		b.SendBroadcast(id, "[SERVER] " + username + " is in the public channel")
		return
	}

	if (b.IsInPrivate(id, &output_user)) {
		b.SendTo(output_user, username + ": " + msg)
	} else {
		b.SendBroadcast(id, username + ": " + msg)
	}
}

//!+handleConn
func handleConn(conn net.Conn, broadcast * broadcaster_lib.BroadcastType) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()

	username := getUsername(conn, broadcast)

	ch <- "[SERVER] You are " + username
	ch <- "[SERVER] If you want to open a private channel with a username, type: 'priv [username]'"
	
	// Add the new client to the clients list:

	broadcast.AddClient(who, username, ch)

	broadcast.SendBroadcast(who, "[SERVER] " + "[" + username + "]" + " has arrived")

	// Announce the new list:

	go broadcast.AnnounceClients()

	input := bufio.NewScanner(conn)
	for input.Scan() {
		str := input.Text()
		sendMsg(ch, broadcast, who, username, str)
	}
	// NOTE: ignoring potential errors from input.Err()

	err := conn.Close()
	if (err != nil) {
		log.Print(err)
	}

	// Remove the client from teh clients list:

	broadcast.DeleteClient(who)

	go broadcast.SendBroadcast(who, "[SERVER] " + "[" + username + "]" + " has left")
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
