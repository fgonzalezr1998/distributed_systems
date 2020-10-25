package main

import (
	"io"
	"log"
	"net"
	"time"
	"os"
	"os/signal"
	"syscall"
	"flag"
	"strconv"
)

const DELAY int = 1

func sigtermHandler() {
	c := make(chan os.Signal)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		io.WriteString(os.Stdout, "[INFO] Program exited cleanly\n")
		os.Exit(0)
	}()
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return
		}
		time.Sleep(time.Duration(DELAY) * time.Second)
	}
}

func main() {
	var port int

	sigtermHandler()  // ctrl+C handler

	flag.IntVar(&port, "port", 8000, "Binding port")
	flag.Parse()

	listener, err := net.Listen("tcp", "localhost:" + strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}
