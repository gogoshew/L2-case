package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

var (
	listen = flag.Bool("l", false, "Listen")
	host   = flag.String("h", "localhost", "Host")
	port   = flag.Int("p", 0, "Port") // 0 value means that we could choose random port
)

func main() {
	protocol := "tcp"
	flag.Parse()
	fmt.Println(flag.Args())

	if *listen {
		startTcpServer(protocol)
		return
	}
	if len(flag.Args()) < 2 {
		fmt.Println("Hostname and port required")
		return
	}

	serverHost := flag.Arg(0)
	serverPort := flag.Arg(1)
	startTcpClient(fmt.Sprintf("%s:%s", serverHost, serverPort), protocol)
}

func startTcpServer(network string) {
	addr := fmt.Sprintf("%s:%d", *host, *port)
	listener, err := net.Listen(network, addr)

	if err != nil {
		panic(err)
	}

	log.Printf("Listening for connections on %s with %s network",
		listener.Addr().String(), strings.ToTitle(network))

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection from client: %s", err)
		} else {
			go processClient(conn)
		}
	}
}

func processClient(conn net.Conn) {
	_, err := io.Copy(os.Stdout, conn)
	if err != nil {
		fmt.Println(err)
	}
	conn.Close()
}

func startTcpClient(addr string, network string) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		fmt.Printf("Can't connect to server: %s\n", err)
		return
	}
	_, err = io.Copy(conn, os.Stdin)
	if err != nil {
		fmt.Printf("Connection error: %s\n", err)
	}
}
