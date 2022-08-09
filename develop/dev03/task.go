package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type SortConfig struct {
	setColumn    int
	sortByDigit  bool
	reversedSort bool
	uniqueRows   bool
	filename     string
}

func createSortConfig() *SortConfig {
	s := SortConfig{}
	flag.IntVar(&s.setColumn, "k", 0, "указываем колонку для сортировки")
	flagN := flag.Bool("n", false, "сортируем по числовому значению")
	flagR := flag.Bool("r", false, "сортируем в обратном порядке")
	flagU := flag.Bool("u", false, "выводим только уникальные значения")

	flag.Parse()

	args := flag.Args()
	s.sortByDigit = *flagN
	s.reversedSort = *flagR
	s.uniqueRows = *flagU

	if len(args) == 1 {
		s.filename = args[0]
	} else {
		log.Fatalf("Единовременно может сортироваться только один файл!")
	}
	return &s
}

func readFile(filename string) ([]string, error) {
	var rows []string
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return rows, err
	}

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		rows = append(rows, sc.Text())
	}
	return rows, nil
}

func getUnique(rows []string) []string {
	uniqueMap := make(map[string]bool)
	for _, row := range rows {
		uniqueMap[row] = true
	}
	// Создаем пустой слайс пол уникальные строки на основе начального
	uniqueRows := rows[:0]
	for key := range uniqueMap {
		uniqueRows = append(uniqueRows, key)
	}
	return uniqueRows
}

func start(s *SortConfig) (string, error) {
	rows, err := readFile(s.filename)
	if err != nil {
		return "", fmt.Errorf("не могу прочитать файл %s\n, %s", s.filename, err.Error())
	}
	return sortRows(rows, s)
}

func getColumn(row string, s *SortConfig) (string, error) {
	re := regexp.MustCompile(`\s+`)
	// Будем отрезать хвостовые пробелы по умолчанию, без вызова флага -b
	columnsSlice := re.Split(strings.TrimSpace(row), -1)
	if len(columnsSlice) >= s.setColumn {
		return columnsSlice[s.setColumn-1], nil
	}
	return "", fmt.Errorf("не могу найти такой столбец")
}

func sortRows(rows []string, s *SortConfig) (string, error) {
	switch {
	case s.setColumn > 0:
		// делаем строки уникальными если true
		if s.uniqueRows {
			rows = getUnique(rows)
		}
		sort.SliceStable(rows, func(i, j int) bool {
			ix, err := getColumn(rows[i], s)
			if err != nil {
				return false
			}
			jx, err := getColumn(rows[j], s)
			if err != nil {
				return false
			}
			// Сортируем по числам
			if s.sortByDigit {
				if s.reversedSort {
					// ! - изменит порядок в обратную сторону
					return !iLessJ(ix, jx)
				}
				return iLessJ(ix, jx)
			}
			if s.reversedSort {
				return ix < jx
			}
			return ix > jx
		})
	case s.uniqueRows:
		rows = getUnique(rows)
		sort.SliceStable(rows, func(i, j int) bool {
			if s.reversedSort {
				return rows[i] > rows[j]
			}
			return rows[i] < rows[j]
		})
	case s.sortByDigit:
		if s.uniqueRows {
			rows = getUnique(rows)
		}
		sort.SliceStable(rows, func(i, j int) bool {
			if s.reversedSort {
				return !iLessJ(rows[i], rows[j])
			}
			return iLessJ(rows[i], rows[j])
		})

	default:
		sort.SliceStable(rows, func(i, j int) bool {
			if s.reversedSort {
				return rows[i] > rows[j]
			}
			return rows[i] < rows[j]
		})
	}

	var res strings.Builder
	rowLength := len(rows)
	for i, row := range rows {
		// Пишем строку, если она не последняя, то добавляем перенос строки
		if i < rowLength-1 {
			_, _ = res.WriteString(row + "\n")
		} else {
			_, _ = res.WriteString(row)
		}
	}

	return res.String(), nil
}

// Создадим функцию, как в примере документации к пакету sort для SliceStable
func iLessJ(strI, strJ string) bool {
	i, _ := strconv.Atoi(strI)
	j, _ := strconv.Atoi(strJ)
	return i < j
}
func main() {
	s := createSortConfig()
	res, err := start(s)
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(res)
}
