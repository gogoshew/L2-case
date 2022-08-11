package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type Config struct {
	fields    string
	delimiter string
	separated bool
}

func NewConfig() *Config {
	config := Config{}

	flag.StringVar(&config.fields, "f", "", "List of fields to cut")
	flag.StringVar(&config.delimiter, "d", "\t", "Set custom delimeter")
	flagS := flag.Bool("s", false, "Get only separated strings")
	flag.Parse()

	config.separated = *flagS
	return &config
}

func cut(row string, config *Config) (string, error) {
	var res strings.Builder
	fields := make(map[int]bool)

	var delimeter string
	if config.delimiter != "\t" {
		if len(config.delimiter) == 1 {
			delimeter = config.delimiter
		} else {
			return "", fmt.Errorf("you could set only one character for delimeter")
		}
	}

	// Ограничиваем поля для cut
	if config.fields != "" {
		diapason := strings.Split(config.fields, ",")

		for _, dPart := range diapason {
			dRange := strings.Split(strings.TrimSpace(dPart), "-")
			// Проверим если диапазон задан из двух значений
			if len(dRange) == 2 {
				// Левая граница диапазона
				dLeft, err := strconv.Atoi(dRange[0])
				if err != nil {
					return "", fmt.Errorf("invalid left value %s", dLeft)
				}
				// Правая граница диапазона
				dRight, err := strconv.Atoi(dRange[1])
				if err != nil {
					return "", fmt.Errorf("invalid right value %s", dRight)
				}

				// Условия задания диапазона
				if dLeft < 1 || dLeft > dRight {
					return "", fmt.Errorf("your range has started from 0 or left border more than right border")
				}
				for i := dLeft; i <= dRight; i++ {
					fields[i] = true
				}
			} else {
				numOfField, err := strconv.Atoi(strings.TrimSpace(dPart))
				if err != nil {
					return "", fmt.Errorf("invalid field value %s", dPart)
				}
				fields[numOfField] = true
			}
		}
	}
	sliceOfRows := strings.Split(row, delimeter)
	if len(sliceOfRows) == 1 && config.separated {
		return "", nil
	}

	isDelim := false
	for i, val := range sliceOfRows {
		_, ok := fields[i+1]
		if ok {
			if isDelim {
				res.WriteString(delimeter + val)
			} else {
				res.WriteString(val)
				isDelim = true
			}
		}
	}
	return res.String(), nil
}

func startCut(config *Config) {
	var str strings.Builder
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		str.WriteString(sc.Text())
	}

	res, err := cut(str.String(), config)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println(res)
}

func main() {
	config := NewConfig()
	startCut(config)
}
