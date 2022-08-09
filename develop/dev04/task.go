package main

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func addToDict(words []string) map[int]string {
	charCountDict := make(map[int]string)
	// Определяем количество букв в слове, преобразуем слова в нижний регистр и отрезаем пробелы
	// добавляем полученные слова в мапу по ключу, равному количеству букв в слове
	for _, word := range words {
		charCountDict[utf8.RuneCountInString(word)] += strings.ToLower(strings.TrimSpace(word)) + " "
	}
	return charCountDict
}

func checkAnagrams(word string, dict map[int]string) []string {
	okWord := strings.ToLower(strings.TrimSpace(word))
	okWordLen := utf8.RuneCountInString(okWord)

	// Создадим регулярное выражение с помощью билдера в виде [w,o,r,d]{len(word)}
	var reg strings.Builder
	reg.WriteString("[")

	for _, char := range okWord {
		reg.WriteString(string(char) + ",")
	}
	reg.WriteString("]{" + strconv.Itoa(okWordLen) + "}")
	// Найдем анаграмму с помощью регулярного выражение
	re := regexp.MustCompile(reg.String())
	return re.FindAllString(dict[okWordLen], -1)
}

// Функция принимает слайс слов и мапу со словами, слова из слайса будут являться ключами,
// а найденные анаграммы - значениями, на выходе функция возвращает словарь анаграмм по введенным словам.
func getAnagrams(words []string, dict map[int]string) map[string][]string {
	res := make(map[string][]string)
	for _, word := range words {
		anagrams := checkAnagrams(word, dict)
		// Проверяем, чтобы в слайсе было больше 1 анаграммы
		if len(anagrams) > 1 {
			sort.Strings(anagrams)
			res[strings.ToLower(strings.TrimSpace(word))] = anagrams
		}
	}
	return res
}

func main() {
	myDict := addToDict([]string{"ЛАСКОВ", "СЛОВАК", "СЛАВОК", "СКОВАЛ", "ВИДАЛ", "ВДАЛИ", "ВЛАДИ",
		"ПЯТАК", "ПЯТКА", "ТЯПКА", "ЛИСТОК", "СЛИТОК", "СТОЛИК", "АВТОР", "ВАРТО", "ВТОРА", "ОТВАР",
		"РВОТА", "ТАВРО", "ТОВАР", "КАЧУР", "КРАУЧ", "КРУЧА", "КУРЧА", "РУЧКА", "ЧУРКА", "АБНЯ",
		"БАНЯ", "БАЯН", "КОРТ", "КРОТ", "ТРОК", "КОТ", "КТО", "ОТК", "ТОК",
	})
	fmt.Println(getAnagrams([]string{}, myDict))
}
