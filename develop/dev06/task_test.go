package main

import (
	"fmt"
	"reflect"
	"testing"
)

// Тест для функции parseFields
func TestParseFields(t *testing.T) {
	cases := []struct {
		input    string
		expected []int
		err      bool
	}{
		{"1,2,3", []int{0, 1, 2}, false},
		{"4,5,6", []int{3, 4, 5}, false},
		{"1", []int{0}, false},
		{"", nil, true},
		{"1,a,3", nil, true},
	}

	for _, c := range cases {
		result, err := parseFields(c.input)
		if (err != nil) != c.err {
			t.Errorf("parseFields(%q) error = %v, expected error = %v", c.input, err, c.err)
			continue
		}
		if !reflect.DeepEqual(result, c.expected) {
			t.Errorf("parseFields(%q) = %v, expected %v", c.input, result, c.expected)
		}
	}
}

// Тест для функции selectFields
func TestSelectFields(t *testing.T) {
	cases := []struct {
		columns  []string
		fields   []int
		expected []string
	}{
		{[]string{"a", "b", "c"}, []int{0, 2}, []string{"a", "c"}},
		{[]string{"a", "b", "c"}, []int{1}, []string{"b"}},
		{[]string{"a", "b", "c"}, []int{0, 1, 2}, []string{"a", "b", "c"}},
		{[]string{"a", "b"}, []int{1, 2}, []string{"b"}},
		{[]string{"a"}, []int{1}, []string{}},
	}

	for _, c := range cases {
		result := selectFields(c.columns, c.fields)
		fmt.Printf("Testing selectFields(%v, %v): got %v, expected %v\n", c.columns, c.fields, result, c.expected)
		if !reflect.DeepEqual(result, c.expected) {
			t.Errorf("selectFields(%v, %v) = %v, expected %v", c.columns, c.fields, result, c.expected)
		}
	}
}
