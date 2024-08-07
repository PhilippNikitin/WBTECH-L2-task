package main

import (
	"reflect"
	"testing"
)

func TestFindAnagramSets(t *testing.T) {
	tests := []struct {
		input    []string
		expected map[string][]string
	}{
		{
			input: []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "кот", "ток", "окт"},
			expected: map[string][]string{
				"пятак":  {"пятак", "пятка", "тяпка"},
				"листок": {"листок", "слиток", "столик"},
				"кот":    {"кот", "окт", "ток"},
			},
		},
		{
			input: []string{"автор", "товар", "отвар", "втор"},
			expected: map[string][]string{
				"автор": {"автор", "отвар", "товар"},
			},
		},
		{
			input: []string{"дом", "мод", "мода", "домик"},
			expected: map[string][]string{
				"дом": {"дом", "мод"},
			},
		},
		{
			input:    []string{"слово", "антивирус", "компьютер"},
			expected: map[string][]string{},
		},
		{
			input: []string{"мама", "мам", "амам", "маам", "амма"},
			expected: map[string][]string{
				"мама": {"амам", "амма", "маам", "мама"},
			},
		},
	}

	for _, test := range tests {
		result := findAnagramSets(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("findAnagramSets(%v) = %v; expected %v", test.input, result, test.expected)
		}
	}
}

func TestSortString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"пятак", "акптя"},
		{"слово", "влоос"},
		{"мама", "аамм"},
		{"кот", "кот"},
	}

	for _, test := range tests {
		result := sortString(test.input)
		if result != test.expected {
			t.Errorf("sortString(%v) = %v; expected %v", test.input, result, test.expected)
		}
	}
}
