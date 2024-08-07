package main

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	fieldsFlag := flag.String("f", "", "fields to select (e.g., 1,2,3)")
	delimiterFlag := flag.String("d", "\t", "delimiter to use")
	separatedFlag := flag.Bool("s", false, "only lines with delimiter")

	flag.Parse()

	if *fieldsFlag == "" {
		fmt.Fprintln(os.Stderr, "fields (-f) flag is required")
		os.Exit(1)
	}

	fields, err := parseFields(*fieldsFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid fields: %v\n", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		if *separatedFlag && !strings.Contains(line, *delimiterFlag) {
			continue
		}

		columns := strings.Split(line, *delimiterFlag)
		output := selectFields(columns, fields)
		fmt.Println(strings.Join(output, *delimiterFlag))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
		os.Exit(1)
	}
}

func parseFields(fieldsStr string) ([]int, error) {
	var fields []int
	for _, field := range strings.Split(fieldsStr, ",") {
		var f int
		_, err := fmt.Sscanf(field, "%d", &f)
		if err != nil {
			return nil, err
		}
		fields = append(fields, f-1) // Convert to 0-based index
	}
	return fields, nil
}

func selectFields(columns []string, fields []int) []string {
	selected := []string{}
	for _, f := range fields {
		if f >= 0 && f < len(columns) {
			selected = append(selected, columns[f])
		}
	}
	return selected
}
