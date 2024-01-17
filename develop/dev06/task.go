package main

import (
	"bufio"
	"flag"
	"fmt"
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
	Fields    string
	Delimiter string
	Separated bool
	FilePath  string
}

func parseFlags() Config {
	config := Config{}
	// флаги, которые будем затем использовать для cut
	flag.StringVar(&config.Fields, "f", "", "выбрать поля (колонки)")
	flag.StringVar(&config.Delimiter, "d", "\t", "использовать другой разделитель")
	flag.BoolVar(&config.Separated, "s", false, "только строки с разделителем")

	flag.Parse()

	//путь до файла
	config.FilePath = flag.Arg(0)

	return config
}

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

// разбиваем пользовательский ввод тех полей, которые нам нужны
func parseFields(fieldsStr string) map[int]bool {
	fields := make(map[int]bool)

	if fieldsStr == "" {
		return fields
	}

	fieldList := strings.Split(fieldsStr, ",")
	for _, fieldStr := range fieldList {
		field := strings.Trim(fieldStr, " ")
		fieldsNum, err := strconv.Atoi(field)
		if err != nil {
			_, err := fmt.Fprintln(os.Stderr, "Error parsing fields:", err)
			if err != nil {
				return nil
			}
			os.Exit(1)
		}
		fields[fieldsNum] = true
	}

	return fields
}

// проходим по всему файлу
func processInput(lines []string, selectedFields map[int]bool, config Config) {
	for _, line := range lines {
		//если есть флаг s, то мы можем сразу пропускать строки без delimiter
		if config.Separated && !strings.Contains(line, config.Delimiter) {
			continue
		}
		//разбиваем на поля по delimiter
		fields := strings.Split(line, config.Delimiter)
		selectedFields := selectFields(fields, selectedFields)

		fmt.Println(strings.Join(selectedFields, config.Delimiter))
	}
}

// если данное поле нам нужно, то мы добавляем его
func selectFields(fields []string, selectedFields map[int]bool) []string {
	result := make([]string, 0)

	for i, field := range fields {
		if selectedFields[i+1] {
			result = append(result, field)
		}
	}

	return result
}

func main() {
	config := parseFlags()

	fields := parseFields(config.Fields)

	lines, err := readLines(config.FilePath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	processInput(lines, fields, config)
}
