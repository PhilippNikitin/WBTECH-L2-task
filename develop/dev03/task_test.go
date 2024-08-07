package main

import (
	"testing"
)

// Тестируем функцию less для различных случаев
func TestLess(t *testing.T) {
	tests := []struct {
		name     string
		a, b     string
		sf       SortFlags
		expected bool
	}{
		{"Simple strings", "apple", "banana", SortFlags{}, true},
		{"Numeric strings", "10", "2", SortFlags{numeric: true}, false},
		{"Non-numeric strings", "apple", "2", SortFlags{numeric: true}, false},
		{"Month sorting", "Jan", "Feb", SortFlags{month: true}, true},
		{"Ignoring blanks", "   apple  ", "banana", SortFlags{ignoreBlanks: true}, true},
		{"Reverse sorting", "banana", "apple", SortFlags{reverse: true}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := less(tt.a, tt.b, tt.sf)
			if got != tt.expected {
				t.Errorf("less(%q, %q, %+v) = %v; want %v", tt.a, tt.b, tt.sf, got, tt.expected)
			}
		})
	}
}

// Тестируем функцию unique
func TestUnique(t *testing.T) {
	lines := []string{"apple", "banana", "apple", "cherry"}
	expected := []string{"apple", "banana", "cherry"}

	got := unique(lines)
	if len(got) != len(expected) {
		t.Fatalf("unique() = %v; want %v", got, expected)
	}

	for i, line := range expected {
		if got[i] != line {
			t.Errorf("unique()[%d] = %q; want %q", i, got[i], line)
		}
	}
}

// Тестируем функцию sortLines
func TestSortLines(t *testing.T) {
	lines := []string{"banana", "apple", "cherry"}
	sf := SortFlags{
		reverse: true,
	}

	expected := []string{"cherry", "banana", "apple"}
	got := sortLines(lines, sf)
	if len(got) != len(expected) {
		t.Fatalf("sortLines() = %v; want %v", got, expected)
	}

	for i, line := range expected {
		if got[i] != line {
			t.Errorf("sortLines()[%d] = %q; want %q", i, got[i], line)
		}
	}
}
