package main

import (
	"io"
	"os"
	"os/signal"
	"syscall"
	"net"
	"sync"
	"bufio"
	"strings"
	"fmt"
)

func sigtermHandler() {
	c := make(chan os.Signal)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		io.WriteString(os.Stdout, "[INFO] Program exited cleanly\n")
		os.Exit(0)
	}()
}

func argOk(arg string) bool {
	/*
	 * Returns true if 'arg' follow the format: foo=foo:foo
	 */

	var s string

	if (strings.Index(arg, "=") == -1) {
		return false
	}

	s = arg[strings.Index(arg, "=") + 1:]
	
	if (strings.Index(s, ":") == -1) {
		return false
	}

	return true
}

func argsOk(args []string) bool {
	for _, arg := range args {

		if (!argOk(arg)) {
			return false
		}
	}

	return true
}

func listen(wg *sync.WaitGroup, timezone string, conn net.Conn) {
	reader := bufio.NewReader(conn)

	defer wg.Done()
	for {
		str, err := reader.ReadString('\n')
		if (err != nil) {
			break
		}
		io.WriteString(os.Stdout, timezone + ": " + str)
	}
	
}

func listenForAll(args []string) {
	var address string
	var wg sync.WaitGroup

	for _, arg := range args {
		address = arg[strings.Index(arg, "=") + 1:]

		conn, err := net.Dial("tcp", address)

		fmt.Println("Connecting to " + address)

		if err != nil {
			io.WriteString(os.Stderr, "[ERROR] Connection refused!\n")
			continue
		}
		wg.Add(1)
		go listen(&wg, arg[:strings.Index(arg, "=")], conn)
	}

	wg.Wait()
}

func main() {
	sigtermHandler()  // ctrl+C handler

	if (!argsOk(os.Args[1:])) {
		io.WriteString(os.Stderr, "[ERROR] Invalid argument!\n")
		os.Exit(1)
	}

	listenForAll(os.Args[1:])

	os.Exit(0)
}