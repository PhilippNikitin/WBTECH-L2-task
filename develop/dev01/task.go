package main

/*
=== Базовая задача ===

Создать программу печатающую точное время с использованием NTP библиотеки.Инициализировать как go module.
Использовать библиотеку https://github.com/beevik/ntp.
Написать программу печатающую текущее время / точное время с использованием этой библиотеки.

Программа должна быть оформлена с использованием как go module.
Программа должна корректно обрабатывать ошибки библиотеки: распечатывать их в STDERR и возвращать ненулевой код выхода в OS.
Программа должна проходить проверки go vet и golint.
*/

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/beevik/ntp"
)

// NTPClient представляет интерфейс для получения времени с NTP сервера
type NTPClient interface {
	Time(host string) (time.Time, error)
}

// RealNTPClient представляет реальный NTP клиент
type RealNTPClient struct{}

// Time возвращает текущее время с указанного NTP сервера
func (c RealNTPClient) Time(host string) (time.Time, error) {
	return ntp.Time(host)
}

// getCurrentTime получает текущее время с NTP сервера через предоставленный клиент
func getCurrentTime(client NTPClient) (time.Time, error) {
	return client.Time("0.beevik-ntp.pool.ntp.org")
}

// printCurrentTime печатает текущее время, используя предоставленный клиент NTP
func printCurrentTime(client NTPClient) {
	currentTime, err := getCurrentTime(client)
	if err != nil {
		// Логируем ошибку в STDERR в случае ее возникновения
		log.Printf("Ошибка получения времени с NTP сервера: %v", err)
		// Возвращаем ненулевой код выхода
		os.Exit(1)
	}

	// Печатаем текущее время
	fmt.Printf("Текущее время: %s\n", currentTime.Format(time.RFC1123))
}

func main() {
	client := RealNTPClient{}
	printCurrentTime(client)
}
