package main

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	// Определение флагов командной строки
	timeoutFlag := flag.String("timeout", "10s", "Timeout for connecting to the server")
	host := flag.String("host", "", "Host to connect to")
	port := flag.String("port", "", "Port to connect to")
	flag.Parse()

	if *host == "" || *port == "" {
		fmt.Println("Host and port must be specified")
		os.Exit(1)
	}

	// Парсинг таймаута
	timeout, err := time.ParseDuration(*timeoutFlag)
	if err != nil {
		fmt.Printf("Invalid timeout value: %v\n", err)
		os.Exit(1)
	}

	// Формирование адреса для подключения
	address := fmt.Sprintf("%s:%s", *host, *port)

	// Создание TCP-соединения
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		fmt.Printf("Error connecting to %s: %v\n", address, err)
		os.Exit(1)
	}
	defer conn.Close()

	// Создание каналов для ввода и вывода данных
	done := make(chan struct{})

	go func() {
		_, err := io.Copy(conn, os.Stdin)
		if err != nil && err != io.EOF {
			fmt.Printf("Error while copying from stdin to socket: %v\n", err)
		}
		done <- struct{}{}
	}()

	go func() {
		_, err := io.Copy(os.Stdout, conn)
		if err != nil && err != io.EOF {
			fmt.Printf("Error while copying from socket to stdout: %v\n", err)
		}
		done <- struct{}{}
	}()

	// Ожидание завершения работы
	<-done
}
