package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type Config struct {
	after       int
	before      int
	contextRows int
	count       bool
	ignoreCase  bool
	invert      bool
	fixed       bool
	strNum      bool
	reg         string
	filename    string
}

// NewConfig Конструктор конфига флагов
func NewConfig() *Config {
	config := Config{}
	flag.IntVar(&config.after, "A", 0, "печатать +N строк после совпадения")
	flag.IntVar(&config.before, "B", 0, "печатать +N строк до совпадения")
	flag.IntVar(&config.contextRows, "C", 0, "(A+B) печатать ±N строк вокруг совпадения")
	flagC := flag.Bool("c", false, "количество строк")
	flagI := flag.Bool("i", false, "игнорировать регистр")
	flagV := flag.Bool("v", false, "вместо совпадения, исключать")
	flagF := flag.Bool("f", false, "точное совпадение со строкой, не паттерн")
	flagN := flag.Bool("n", false, "печатать номер строки")

	flag.Parse()

	args := flag.Args()
	config.count = *flagC
	config.ignoreCase = *flagI
	config.invert = *flagV
	config.fixed = *flagF
	config.strNum = *flagN

	if len(args) == 2 {
		config.reg = args[0]
		config.filename = args[1]
	} else {
		log.Fatalf("Одновременно может проверяться только один файл")
	}
	return &config
}

// Функция чтения файла
func readFile(filename string) ([]string, error) {
	var rows []string
	file, err := os.Open(filename)
	if err != nil {
		return rows, err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		rows = append(rows, sc.Text())
	}
	return rows, nil
}

func grep(rows []string, config *Config) (interface{}, error) {

	var prefix, postfix string

	if config.ignoreCase {
		// Одно или ноль вхождений + игнорирование регистра
		prefix = "(?i)"
	}
	if config.fixed {
		// ^ Начало строки
		prefix = "^"
		// $ Конец строки
		postfix = "$"
	}
	// Формируем регулярку
	re, err := regexp.Compile(prefix + config.reg + postfix)
	if err != nil {
		return "Error", fmt.Errorf("неверное регулярное выражение ")
	}

	switch {
	case config.after != 0:
		for i, row := range rows {
			// Проверим, содержит ли строка соответствие регулярному выражению
			if re.MatchString(row) {
				if config.after <= len(rows)-i {
					return rows[i : i+config.after+1], nil
				}
				return rows[i:], nil
			}
		}
		return "не найдено", nil

	case config.before != 0:
		for i, row := range rows {
			// Проверим, содержит ли строка соответствие регулярному выражению
			if re.MatchString(row) {
				if config.before-1 <= i {
					return rows[i-config.before : i+1], nil
				}
				return rows[:i+1], nil
			}
		}
		return "не найдено", nil

	case config.contextRows != 0:
		for i, row := range rows {
			if re.MatchString(row) {
				firstI, lastI := 0, len(rows)

				if config.contextRows-1 <= i {
					firstI = i - config.contextRows
				}

				if config.contextRows <= len(rows)-i {
					lastI = i + config.contextRows + 1
				}

				return rows[firstI:lastI], nil
			}
		}
		return "не найдено", nil

	case config.count:
		total := 0
		for _, row := range rows {
			total += len(re.FindAllString(row, -1))
		}
		if config.invert {
			return len(rows) - total, nil
		}
		return total, nil

	case config.strNum:
		var numOfRows []int
		for i, row := range rows {
			if re.MatchString(row) {
				numOfRows = append(numOfRows, i)
			}
		}
		return numOfRows, nil

	default:
		var res []string
		for _, row := range rows {
			if re.MatchString(row) {
				res = append(res, row)
			}
		}
		return res, nil
	}
}

func initGrep(config *Config) (interface{}, error) {
	rows, err := readFile(config.filename)
	if err != nil {
		return "", fmt.Errorf("не могу прочитать файл %s: %s", config.filename, err.Error())
	}
	return grep(rows, config)
}

func main() {
	config := NewConfig()
	res, err := initGrep(config)
	if err != nil {
		log.Fatalf(err.Error())
	}

	switch result := res.(type) {
	case []string:
		for _, row := range result {
			fmt.Println(row)
		}
	case []int:
		for _, row := range result {
			fmt.Println(row)
		}
	case string:
		fmt.Println(result)
	case int:
		fmt.Println(result)
	default:
		fmt.Printf("Неизвестный тип %T\n", result)
	}
}
