package main

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

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
)

func main() {
	// Define flags
	after := flag.Int("A", 0, "Print N lines after each match")
	before := flag.Int("B", 0, "Print N lines before each match")
	context := flag.Int("C", 0, "Print N lines before and after each match (context)")
	count := flag.Bool("c", false, "Count matching lines only")
	ignoreCase := flag.Bool("i", false, "Ignore case distinctions")
	invert := flag.Bool("v", false, "Invert match (select non-matching lines)")
	fixed := flag.Bool("F", false, "Fixed string match (literal match)")
	lineNumber := flag.Bool("n", false, "Print line numbers")
	flag.Parse()

	// Adjust context flags if necessary
	if *context > 0 {
		*after = *context
		*before = *context
	}

	// Fetch the pattern to search for
	pattern := flag.Arg(0)
	if pattern == "" {
		fmt.Println("Usage: grep [OPTIONS] pattern [file ...]")
		flag.PrintDefaults()
		return
	}

	// Compile the pattern for matching
	var regex *regexp.Regexp
	if *fixed {
		pattern = regexp.QuoteMeta(pattern) // Treat pattern as literal string
	}
	if *ignoreCase {
		regex = regexp.MustCompile("(?i)" + pattern)
	} else {
		regex = regexp.MustCompile(pattern)
	}

	// Process files or stdin if no files are provided
	files := flag.Args()[1:]
	if len(files) == 0 {
		// Read from stdin
		processInput(os.Stdin, regex, *after, *before, *count, *invert, *lineNumber)
	} else {
		// Process each file
		for _, filename := range files {
			file, err := os.Open(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening file %s: %v\n", filename, err)
				continue
			}
			processInput(file, regex, *after, *before, *count, *invert, *lineNumber)
			file.Close()
		}
	}
}

func processInput(input *os.File, regex *regexp.Regexp, after, before int, count, invert, lineNumber bool) {
	scanner := bufio.NewScanner(input)
	var (
		lineNum       int
		matchCount    int
		lines         []string
		matchLineNums []int
	)

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		lineNum++
		if regex.MatchString(line) != invert {
			matchCount++
			matchLineNums = append(matchLineNums, lineNum)
		}
	}

	if count {
		fmt.Println(matchCount)
		return
	}

	lineMap := make(map[int]bool)
	for _, matchNum := range matchLineNums {
		start := matchNum - before - 1
		if start < 0 {
			start = 0
		}
		end := matchNum + after
		if end > len(lines) {
			end = len(lines)
		}
		for i := start; i < end; i++ {
			lineMap[i] = true
		}
	}

	for i := 0; i < len(lines); i++ {
		if lineMap[i] {
			if lineNumber {
				fmt.Printf("%d:%s\n", i+1, lines[i])
			} else {
				fmt.Println(lines[i])
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}
