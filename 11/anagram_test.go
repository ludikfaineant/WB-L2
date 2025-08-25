package main

import (
	"reflect"
	"testing"
)

func TestFindAnagram(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  map[string][]string
	}{
		{
			name:  "basic anagrams",
			input: []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"},
			want: map[string][]string{
				"пятак":  {"пятак", "пятка", "тяпка"},
				"листок": {"листок", "слиток", "столик"},
			},
		},
		{
			name:  "no anagrams",
			input: []string{"стол", "стул", "шкаф"},
			want:  map[string][]string{},
		},
		{
			name:  "mixed case",
			input: []string{"Пятак", "пятка", "Тяпка", "листок", "Слиток"},
			want: map[string][]string{
				"Пятак":  {"Пятак", "Тяпка", "пятка"},
				"листок": {"Слиток", "листок"},
			},
		},
		{
			name:  "empty input",
			input: []string{},
			want:  map[string][]string{},
		},
		{
			name:  "single word",
			input: []string{"пятак"},
			want:  map[string][]string{},
		},
		{
			name:  "two pairs",
			input: []string{"eat", "tea", "tan", "ate", "nat", "bat"},
			want: map[string][]string{
				"eat": {"ate", "eat", "tea"},
				"tan": {"nat", "tan"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findAnagram(tt.input)

			if len(got) != len(tt.want) {
				t.Fatalf("got %d groups, want %d", len(got), len(tt.want))
			}

			for key, wantGroup := range tt.want {
				gotGroup, exists := got[key]
				if !exists {
					t.Errorf("missing key %q in result", key)
					continue
				}
				if !reflect.DeepEqual(gotGroup, wantGroup) {
					t.Errorf("for key %q: got %v, want %v", key, gotGroup, wantGroup)
				}
			}

			for key := range got {
				if _, exists := tt.want[key]; !exists {
					t.Errorf("unexpected key in result: %q", key)
				}
			}
		})
	}
}
