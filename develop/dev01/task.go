package main

import (
	"fmt"
	"github.com/beevik/ntp"
	"os"
	"time"
)

/*
=== Базовая задача ===

Создать программу печатающую точное время с использованием NTP библиотеки.Инициализировать как go module.
Использовать библиотеку https://github.com/beevik/ntp.
Написать программу печатающую текущее время / точное время с использованием этой библиотеки.

Программа должна быть оформлена с использованием как go module.
Программа должна корректно обрабатывать ошибки библиотеки: распечатывать их в STDERR и возвращать ненулевой код выхода в OS.
Программа должна проходить проверки go vet и golint.
*/

func main() {
	response, err := ntp.Query("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error at geeting time: %s\n", err)
		os.Exit(1)
	}
	// Получение точного времени
	ntpTime := time.Now().Add(response.ClockOffset)
	// Печать текущего и точного времени
	fmt.Printf("Текущее время: %v\n", time.Now())
	fmt.Printf("Точное время (по NTP): %v\n", ntpTime)
}
