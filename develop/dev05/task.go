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

// Config содержит параметры фильтрации.
type Config struct {
	AfterLines   int
	BeforeLines  int
	ContextLines int
	Count        bool
	IgnoreCase   bool
	InvertMatch  bool
	FixedString  bool
	LineNumber   bool
	Pattern      string
	FilePath     string
}

func parseFlags() Config {
	config := Config{}

	// флаги, которые будем затем использовать для сортировки
	flag.IntVar(&config.AfterLines, "A", 0, "печатать +N строк после совпадения")
	flag.IntVar(&config.BeforeLines, "B", 0, "печатать +N строк до совпадения")
	flag.IntVar(&config.ContextLines, "C", 0, "печатать +N строк вокруг совпадения")
	flag.BoolVar(&config.Count, "c", false, "печатать количество совпадающих строк")
	flag.BoolVar(&config.IgnoreCase, "i", false, "игнорировать регистр")
	flag.BoolVar(&config.InvertMatch, "v", false, "печатать все строки, кроме совпадающей")
	flag.BoolVar(&config.FixedString, "F", false, "точное совпадение со строкой, не паттерн")
	flag.BoolVar(&config.LineNumber, "n", false, "напечтать номер совпадающей строки строки")

	flag.Parse()

	config.Pattern = flag.Arg(0)
	config.FilePath = flag.Arg(1)
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

func filterLines(lines []string, config Config) []string {
	if config.AfterLines > 0 || config.BeforeLines > 0 || config.ContextLines > 0 {
		return filterLinesContext(lines, config)
	}
	return filterLinesDefault(lines, config)
}

func filterLinesDefault(lines []string, config Config) []string {
	var resultLines []string
	var count int

	for i, line := range lines {
		match := matchLine(line, config)

		if match && config.Count {
			count++
			continue
		}

		if config.LineNumber {
			line = fmt.Sprintf("%d:%s", i+1, line)
		}

		if match {
			resultLines = append(resultLines, line)
		}
	}
	if count != 0 {
		resultLines = append(resultLines, strconv.Itoa(count))
	}

	return resultLines
}

func filterLinesContext(lines []string, config Config) []string {
	var resultLines []string
	var buffer []string
	var count int
	var countAfter int

	if config.ContextLines != 0 {
		config.AfterLines = config.ContextLines
		config.BeforeLines = config.ContextLines
	}

	for i, line := range lines {
		match := matchLine(line, config)

		if match && config.Count {
			count++
			continue
		}

		if config.LineNumber {
			line = fmt.Sprintf("%d:%s", i+1, line)
		}

		buffer = append(buffer, line)
		if match {
			if config.AfterLines > 0 {
				countAfter = config.AfterLines
			}

			if config.BeforeLines > 0 {
				if len(buffer) > 1 {
					buffer = buffer[len(buffer)-(config.BeforeLines+1):]
					resultLines = append(resultLines, buffer...)
					buffer = nil
					continue
				}
				buffer = nil
			}

			resultLines = append(resultLines, line)
			continue
		}

		if countAfter != 0 {
			resultLines = append(resultLines, line)
			countAfter--
		}
	}

	if count != 0 {
		resultLines = append(resultLines, strconv.Itoa(count))
	}

	return resultLines
}

func matchLine(line string, config Config) bool {
	//если есть точное совпадение
	if config.FixedString {
		return strings.Contains(line, config.Pattern)
	}

	//если не важен регистр
	if config.IgnoreCase {
		line = strings.ToLower(line)
		config.Pattern = strings.ToLower(config.Pattern)
	}

	match := strings.Contains(line, config.Pattern)

	if config.InvertMatch {
		return !match
	}

	return match
}

func printLines(lines []string) {
	for _, line := range lines {
		fmt.Println(line)
	}
}

func main() {
	config := parseFlags()

	lines, err := readLines(config.FilePath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	resultLines := filterLines(lines, config)
	printLines(resultLines)
}
