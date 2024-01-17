package main

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
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

func Unpack(s string) (string, error) {
	var result string
	if len(s) == 0 {
		return "", nil
	}

	if unicode.IsDigit(rune(s[0])) {
		return "", errors.New("некорректная строка")
	}
	for i, char := range s {
		if unicode.IsDigit(char) {
			count, err := strconv.Atoi(string(char))
			if err != nil {
				return "", errors.New("некорректная строка")
			}
			if unicode.IsDigit(rune(s[i-1])) {
				return "", errors.New("некорректная строка")
			}
			for y := 0; y < count-1; y++ {
				result += string(s[i-1])
			}

		} else {
			result += string(char)
		}
	}

	return result, nil
}

func main() {
	// Пример использования
	testCases := []struct {
		input  string
		output string
		err    error
	}{
		{"a4bc2d5e", "aaaabccddddde", nil},
		{"abcd", "abcd", nil},
		{"45", "", errors.New("некорректная строка")},
		{"", "", nil},
		{"a4b6c8910", "", errors.New("некорректная строка")},
	}

	for _, tc := range testCases {
		result, err := Unpack(tc.input)
		if err != nil {
			if err.Error() != tc.err.Error() {
				fmt.Printf("Ошибка: ожидалось %v, получено %v\n", tc.err, err)
			}
		} else if result != tc.output {
			fmt.Printf("Ошибка: ожидалось %v, получено %v\n", tc.output, result)
		} else {
			fmt.Printf("Тест успешен: %v => %v\n", tc.input, result)
		}
	}
}
