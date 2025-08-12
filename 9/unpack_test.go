package main

import (
	"testing"
)

func TestSimple(t *testing.T) {
	res, err := unpack("a4")
	expected := "aaaa"
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if res != expected {
		t.Errorf("got %q, want %q", res, expected)
	}
}
func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"a4bc2d5e", "aaaabccddddde", false},
		{"abcd", "abcd", false},
		{"", "", false},
		{"45", "", true},
		{"qwe\\4\\5", "qwe45", false},
		{"qwe\\45", "qwe44444", false},
		{"a\\", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			res, err := unpack(tt.input)
			if tt.hasError {
				if err == nil {
					t.Errorf("expected error: %v", err)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if res != tt.expected {
				t.Errorf("got %q, want %q", res, tt.expected)
			}
		})
	}

}
