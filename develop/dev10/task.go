package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему серверу, программа должна завершаться через timeout.
*/

type Config struct {
	TimeOut     time.Duration
	Host        string
	Port        string
	ServerStart bool
	ServerHost  string
	ServerPort  string
}

func NewConfig() *Config {
	config := Config{}

	flag.Usage = func() {
		fmt.Println("Usage flags: --timeout host port")
		flag.PrintDefaults()
	}

	timeoutFlag := flag.Duration("timeout", 10*time.Second, "timeout")

	flag.Parse()
	args := flag.Args()

	if len(args) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	config.Host = args[0]
	config.Port = args[1]
	config.TimeOut = *timeoutFlag
	config.ServerStart = true
	config.ServerHost = "127.0.0.1"
	config.ServerPort = "8080"

	return &config
}

//read Функция чтения из соединения и записи в Stdout
func read(conn net.Conn, cf context.CancelFunc) {
	sc := bufio.NewScanner(conn)
	if !sc.Scan() {
		log.Printf("read: connection closed")
		cf()
		return
	}
	text := sc.Text()
	fmt.Printf("%s\n", text)
}

//write Функция чтения из Stdin и записи в соединение
func write(conn net.Conn, cf context.CancelFunc) {
	sc := bufio.NewScanner(os.Stdin)
	if !sc.Scan() {
		log.Printf("write: can't scan from stdin")
		cf()
		return
	}
	text := sc.Text()
	_, err := conn.Write([]byte(fmt.Sprintln(text)))
	if err != nil {
		log.Printf("write: can't write to server connection")
		cf()
		return
	}
}

func startTcpServer(config *Config, wg *sync.WaitGroup) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", config.ServerHost, config.ServerPort))
	if err != nil {
		log.Fatal("can't listen TCP server\n", err)
	}
	wg.Done()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("can't accept connection\n", err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()
	conn.Write([]byte("Handling " + remoteAddr + "\n\r"))
	log.Printf("%+v connected\n", remoteAddr)
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	for sc.Scan() {
		text := sc.Text()
		if text == "exit" {
			conn.Write([]byte(fmt.Sprintf("Exit from connection, %+v\n\r", remoteAddr)))
			log.Printf("%+v disconnected\n", remoteAddr)
			break
		} else if text == "" {
			conn.Write([]byte(fmt.Sprintf("%+v message is '%s'\n\r", remoteAddr, text)))
		}
	}
}

func main() {
	config := NewConfig()
	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT)
	go func() {
		<-sigCh
		cancel()
	}()

	if config.ServerStart {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go startTcpServer(config, wg)
		wg.Wait()
	}

	addr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(config.Host, config.Port))
	if err != nil {
		log.Fatal("can't resolve tcp address", err)
	}

	conn, err := net.DialTimeout(addr.Network(), addr.String(), config.TimeOut)
	if err != nil {
		log.Fatal("timeout to connection", err)
	}
	defer conn.Close()

	go read(conn, cancel)
	go write(conn, cancel)

	<-ctx.Done()
	log.Println("finished telnet client")
}
