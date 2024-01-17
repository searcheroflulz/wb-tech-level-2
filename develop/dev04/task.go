package main

import (
	"fmt"
	"slices"
	"strings"
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

func searchAnagrams(s []string) map[string][]string {
	anagrams := make(map[string][]string)
	result := make(map[string][]string)

	for _, word := range s {
		//приводим слово к нижнему регистру
		sortedWord := sortString(strings.ToLower(word))

		//проверяем, есть ли уже множество для данного слова
		if s, ok := anagrams[sortedWord]; ok {
			// Если да, добавляем слово в множество
			anagrams[sortedWord] = append(s, word)
		} else {
			// Если нет, создаем новое множество
			anagrams[sortedWord] = []string{word}
		}
	}

	//убираем множества, состоящие из одного элемента
	for key, set := range anagrams {
		set = removeDuplicates(set)
		if len(set) <= 1 {
			delete(anagrams, key)
		} else {
			//сортируем по возрастанию
			result[set[0]] = set[1:]
			slices.Sort(result[set[0]])
		}
	}

	return result
}

// удаляем дупликаты с помощью map
func removeDuplicates(slice []string) []string {
	var result []string
	for i := 0; i < len(slice); i++ {
		duplicate := false
		for j := 0; j < len(result); j++ {
			if slice[i] == result[j] {
				duplicate = true
				break
			}
		}
		if !duplicate {
			result = append(result, slice[i])
		}
	}
	return result
}

// сортируем руны в строке по возрастанию
func sortString(s string) string {
	r := []rune(s)
	slices.Sort(r)
	return string(r)
}

func main() {
	words := []string{"пятак", "ааа", "тяпка", "столик", "листок", "слиток", "пятка", "ааа", "ааа", "ббб", "ввв", "ввв"}
	result := searchAnagrams(words)

	for key, set := range result {
		fmt.Printf("%s: %v\n", key, set)
	}
}
