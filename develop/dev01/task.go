package main

import (
	"github.com/beevik/ntp"
	"log"
	"os"
	"time"
)

/*
=== Базовая задача ===

Создать программу, печатающую точное время с использованием NTP библиотеки. Инициализировать как go module.
Использовать библиотеку https://github.com/beevik/ntp.
Написать программу, печатающую текущее время / точное время с использованием этой библиотеки.

Программа должна быть оформлена с использованием как go module.
Программа должна корректно обрабатывать ошибки библиотеки: распечатывать их в STDERR и возвращать ненулевой код выхода в OS.
Программа должна проходить проверки go vet и golint.

SNTP (англ. Simple Network Time Protocol) — протокол синхронизации времени по компьютерной сети.
Является упрощённой реализацией протокола NTP. Используется во встраиваемых системах и устройствах,
не требующих высокой точности, а также в пользовательских программах точного времени.

go vet проверяет исходный код Go и сообщает о подозрительных конструкциях, таких как вызовы Printf,
аргументы которых не совпадают с форматом строки. vet использует эвристику, которая не гарантирует,
что все отчеты являются подлинными проблемами, но он может найти ошибки, не обнаруженные компиляторами.
*/

func main() {
	// STDERR
	logger := log.New(os.Stderr, "", 0)

	// Querying the current time
	curTime, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		logger.Fatal(err.Error())
	}
	log.Printf("The current time is %s", curTime)

	// Querying time metadata
	response, err := ntp.Query("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		logger.Fatal(err.Error())
	}
	metaTime := time.Now().Add(response.ClockOffset)
	log.Printf("The time metadata is %s", metaTime)
}
