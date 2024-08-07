package main

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

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

// Unpack выполняет распаковку строки, содержащую повторяющиеся символы.
func Unpack(s string) (r string, err error) {
	if _, err := strconv.Atoi(s); err == nil {
		return r, errors.New("некорректная строка")
	}

	var prev rune
	var escaped bool
	var b strings.Builder
	for _, char := range s {
		if unicode.IsDigit(char) && !escaped {
			m := int(char - '0')
			r := strings.Repeat(string(prev), m-1)
			b.WriteString(r)
		} else {
			escaped = string(char) == "\\" && string(prev) != "\\"
			if !escaped {
				b.WriteRune(char)
			}
			prev = char
		}
	}

	return b.String(), err
}
