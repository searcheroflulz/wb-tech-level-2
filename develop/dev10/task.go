package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port
go-telnet mysite.ru 8080
go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

func telnet(host, port string, timeout time.Duration) {
	//канал для отслеживания системных сигналов
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

	address := fmt.Sprintf("%s:%s", host, port)

	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		fmt.Printf("Error connecting to %s: %v\n", address, err)
		os.Exit(1)
	}

	var dnsErr *net.DNSError
	switch {
	case errors.As(err, &dnsErr):
		time.Sleep(timeout)
		log.Println("DNS error: ", err)
		os.Exit(0)
	case err != nil:
		log.Fatal("Cannot open connection: ", err)
	}

	//отслеживаем сигналы и закрываем сокет
	go func(conn net.Conn) {
		<-signalChannel
		err := conn.Close()
		if err != nil {
			log.Println("Cannot close socket")
			os.Exit(1)
		}
		log.Println("Closing connection...")
		os.Exit(0)
	}(conn)

	var wg sync.WaitGroup

	wg.Add(2)

	go readFromConn(conn, &wg, signalChannel)
	go writeToConn(conn, &wg, signalChannel)

	wg.Wait()
}

func writeToConn(conn net.Conn, wg *sync.WaitGroup, signalChannel chan os.Signal) {
	defer wg.Done()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		message := scanner.Text()

		//message = strings.TrimRight(message, "\r\n")

		_, err := conn.Write([]byte(message + "\r"))
		if err != nil {
			signalChannel <- syscall.SIGQUIT
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error reading from stdin:", err)
		signalChannel <- syscall.SIGQUIT
	}
}

func readFromConn(conn net.Conn, wg *sync.WaitGroup, signalChannel chan os.Signal) {
	defer wg.Done()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		fmt.Println(message)
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error reading from connection:", err)
		signalChannel <- syscall.SIGQUIT
	}
}

func parseFlags() (string, string, time.Duration) {
	var host string
	var port string
	var timeout string
	flag.Parse()

	if len(flag.Args()) < 3 {
		fmt.Println("Usage: go-telnet --timeout=10s host port")
		os.Exit(1)
	}

	if flag.Arg(0) != "go-telnet" {
		fmt.Println("Usage: go-telnet --timeout=10s host port")
		os.Exit(1)
	}
	if len(flag.Args()) > 3 {
		timeout = flag.Arg(1)
		timeout, _ = strings.CutPrefix(timeout, "--timeout=")
		host = flag.Arg(2)
		port = flag.Arg(3)
	} else {
		timeout = "10s"
		host = flag.Arg(1)
		port = flag.Arg(2)
	}
	duration, err := time.ParseDuration(timeout)
	if err != nil {
		fmt.Println("Invalid timeout value:", err)
		os.Exit(1)
	}
	return host, port, duration
}

func main() {
	host, port, timeout := parseFlags()

	telnet(host, port, timeout)
}
