package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestCut(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		input    string
		fileData string
		want     string
	}{
		{
			name:  "simple field",
			args:  []string{"cut", "-f", "1", "-d", ","},
			input: "a,b,c\nd,e,f",
			want:  "a\nd\n",
		},
		{
			name:  "multiple fields",
			args:  []string{"cut", "-f", "1,3", "-d", ","},
			input: "a,b,c\nd,e,f",
			want:  "a,c\nd,f\n",
		},
		{
			name:  "range fields",
			args:  []string{"cut", "-f", "2-4", "-d", ":"},
			input: "a:b:c:d:e\n1:2:3:4:5",
			want:  "b:c:d\n2:3:4\n",
		},
		{
			name:  "mixed fields",
			args:  []string{"cut", "-f", "1,3-4,6", "-d", " "},
			input: "a b c d e f\none two three four five six",
			want:  "a c d f\none three four six\n",
		},
		{
			name:  "default delimiter (tab)",
			args:  []string{"cut", "-f", "2"},
			input: "a\tb\tc\nd\te\tf",
			want:  "b\ne\n",
		},
		{
			name:  "separated only (-s)",
			args:  []string{"cut", "-f", "1", "-d", ",", "-s"},
			input: "a,b\nc\nd,e",
			want:  "a\nd\n",
		},
		{
			name:  "field out of range",
			args:  []string{"cut", "-f", "10", "-d", ","},
			input: "a,b,c",
			want:  "", // поле 10 не существует
		},
		{
			name:  "empty field selection",
			args:  []string{"cut", "-f", "3-1", "-d", ","}, // некорректный диапазон
			input: "a,b,c",
			want:  "", // парсинг вернёт пустой список
		},
		{
			name:     "file input",
			args:     []string{"cut", "-f", "1", "-d", ",", "test.csv"},
			fileData: "a,b,c\nd,e,f",
			want:     "a\nd\n",
		},
		{
			name:  "no delimiter no -s",
			args:  []string{"cut", "-f", "1", "-d", ","},
			input: "abc\ndef",
			want:  "abc\ndef\n",
		},
		{
			name:  "empty input",
			args:  []string{"cut", "-f", "1"},
			input: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tempFile string
			if tt.fileData != "" {
				file, err := os.CreateTemp("", "cut_test_*.csv")
				if err != nil {
					t.Fatal(err)
				}
				defer os.Remove(file.Name())
				defer file.Close()

				if _, err := file.WriteString(tt.fileData); err != nil {
					t.Fatal(err)
				}
				if err := file.Close(); err != nil {
					t.Fatal(err)
				}
				for i, arg := range tt.args {
					if arg == "test.csv" {
						tt.args[i] = file.Name()
					}
				}
				tempFile = file.Name()
			}
			cmd := exec.Command("go", append([]string{"run", "cut.go"}, tt.args[1:]...)...)
			if tt.input != "" {
				cmd.Stdin = strings.NewReader(tt.input)
			}

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err := cmd.Run()
			if err != nil && !strings.Contains(err.Error(), "exit status") {
				t.Fatalf("failed to run command: %v, stderr: %s", err, stderr.String())
			}

			if stdout.String() != tt.want {
				t.Errorf("got %q, want %q", stdout.String(), tt.want)
			}
			if tempFile != "" {
				os.Remove(tempFile)
			}
		})
	}
}
