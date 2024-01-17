package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
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

// функция для чтения файла
func readLines(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		os.Exit(1)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func sortLines(lines []string, column int, numericSort bool) {
	less := func(i, j int) bool {
		str1 := extractColumn(lines[i], column)
		str2 := extractColumn(lines[j], column)

		// Преобразование строк в числа, если требуется числовая сортировка
		if numericSort {
			num1, err1 := strconv.Atoi(str1)
			num2, err2 := strconv.Atoi(str2)

			if err1 == nil && err2 == nil {
				return num1 < num2
			}
		}
		// Сравнение строк
		return str1 < str2
	}

	sort.SliceStable(lines, less)
}

func sortFile(lines []string, column int, numericSort, reverseSort, uniqueSort bool) []string {
	if uniqueSort {
		//используем map для уникальной сортировки
		seen := make(map[string]struct{})
		var uniqueLines []string

		for _, line := range lines {
			if _, ok := seen[line]; !ok {
				seen[line] = struct{}{}
				uniqueLines = append(uniqueLines, line)
			}
		}

		sortLines(uniqueLines, column, numericSort)
		lines = uniqueLines
	} else {
		//обычная сортировка
		sortLines(lines, column, numericSort)
	}

	//обратная сортировка
	if reverseSort {
		reverseSlice(lines)
	}

	return lines
}

// выбираем строку по колонке
func extractColumn(line string, column int) string {
	fields := strings.Fields(line)
	if column > 0 && column <= len(fields) {
		return fields[column-1]
	}
	return line
}

// переворачиваем результат сортировки
func reverseSlice(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// записываем в файл результаты сортировки
func writeLines(filePath string, lines []string) error {
	file, err := os.Create(fmt.Sprintf("%v_sorted.txt", filePath))
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := fmt.Fprintln(writer, line)
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

func main() {
	// флаги, которые будем затем использовать для сортировки
	column := flag.Int("k", 0, "номер колонки для сортировки (по умолчанию 0 - вся строка)")
	numericSort := flag.Bool("n", false, "сортировать по числовому значению")
	reverseSort := flag.Bool("r", false, "сортировать в обратном порядке")
	uniqueSort := flag.Bool("u", false, "не выводить повторяющиеся строки")

	flag.Parse()

	filename := flag.Arg(0)

	lines, err := readLines(filename)
	if err != nil {
		os.Exit(1)
	}

	result := sortFile(lines, *column, *numericSort, *reverseSort, *uniqueSort)

	err = writeLines(filename, result)
	if err != nil {
		fmt.Printf("Ошибка при записи в файл: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Файл успешно отсортирован.")
}
