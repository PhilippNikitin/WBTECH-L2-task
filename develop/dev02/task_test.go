package main

import "testing"

func TestUnpack(t *testing.T) {
	data := map[string]string{
		"a4bc2d5e": "aaaabccddddde",
		"abcd":     "abcd",
	}

	for s, e := range data {
		r, err := Unpack(s)
		if err != nil {
			t.Fatalf("bad unpack for %s: got error %v", s, err)
		}
		if r != e {
			t.Fatalf("bad unpack for %s: got %v expected %v", s, r, e)
		}
	}
}

func TestUnpackError(t *testing.T) {
	s := "45"
	r, err := Unpack(s)
	if r != "" {
		t.Fatalf("bad unpack for %s: expected empty string", s)
	}
	if err == nil {
		t.Fatalf("bad unpack for %s: expected error", s)
	}
}

func TestUnpackEscape(t *testing.T) {
	data := map[string]string{
		"qwe\\4\\5": "qwe45",
		"qwe\\45":   "qwe44444",
		"qwe\\\\5":  "qwe\\\\\\\\\\",
	}

	for s, e := range data {
		r, err := Unpack(s)
		if err != nil {
			t.Fatalf("bad unpack for %s: got error %v", s, err)
		}
		if r != e {
			t.Fatalf("bad unpack for %s: got %v expected %v", s, r, e)
		}
	}
}
