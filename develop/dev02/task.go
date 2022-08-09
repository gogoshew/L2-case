package main

import (
	"fmt"
	"log"
	"strings"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	str := "a4bc2d5e"
	unpacked, err := unpackStr(str)
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Printf("Unpacked string: '%s'", unpacked)
}

func unpackStr(str string) (string, error) {
	// Будем собирать строку записывая руны с помощью билдера из пакета strings
	var unpackedStr strings.Builder
	var lastRune rune
	var isEscape bool

	// Запустим цикл по исходной строке
	for _, char := range str {
		// С помощью switch будем проверять условия относительно каждой руны
		switch {

		case isEscape:
			isEscape = false
			lastRune = char

		case char >= '0' && char <= '9':
			if lastRune != 0 {
				// Для перевода из числовой руны в int нужно вычесть из кода этой руны код руны '0'
				iterationCount := int(char - '0')
				for i := 0; i < iterationCount; i++ {
					unpackedStr.WriteRune(lastRune)
				}
				lastRune = 0
			} else {
				return "", fmt.Errorf("incorrect string")
			}

		case char == '\\':
			isEscape = true
			if lastRune != 0 {
				unpackedStr.WriteRune(lastRune)
			}

		default:
			if lastRune != 0 {
				unpackedStr.WriteRune(lastRune)
			}
			lastRune = char
		}
	}

	if lastRune != 0 {
		unpackedStr.WriteRune(lastRune)
	}
	return unpackedStr.String(), nil
}
