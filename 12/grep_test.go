// grep_test.go
package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestGrep(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		input    string
		fileData string
		want     string
		wantCode int
	}{
		{
			name:  "simple match stdin",
			args:  []string{"grep", "match"},
			input: "1\n2\nmatch\n4\n5",
			want:  "match\n",
		},
		{
			name:  "case insensitive",
			args:  []string{"grep", "-i", "MATCH"},
			input: "match\nMATCH\nMatch",
			want:  "match\nMATCH\nMatch\n",
		},
		{
			name:  "invert match",
			args:  []string{"grep", "-v", "match"},
			input: "1\n2\nmatch\n4\n5",
			want:  "1\n2\n4\n5\n",
		},
		{
			name:  "fixed string",
			args:  []string{"grep", "-F", "a.b"},
			input: "a.b\nab",
			want:  "a.b\n",
		},
		{
			name:  "line numbers",
			args:  []string{"grep", "-n", "match"},
			input: "1\nmatch\n3",
			want:  "2:match\n",
		},
		{
			name:  "count lines",
			args:  []string{"grep", "-c", "match"},
			input: "match\nmatch\nno",
			want:  "2\n",
		},
		{
			name:  "after context",
			args:  []string{"grep", "-A", "1", "match"},
			input: "1\nmatch\n3\n4",
			want:  "match\n3\n",
		},
		{
			name:  "before context",
			args:  []string{"grep", "-B", "1", "match"},
			input: "1\n2\nmatch\n4",
			want:  "2\nmatch\n",
		},
		{
			name:  "context",
			args:  []string{"grep", "-C", "1", "match"},
			input: "1\n2\nmatch\n4\n5",
			want:  "2\nmatch\n4\n",
		},
		{
			name:  "context with separator",
			args:  []string{"grep", "-A", "1", "match"},
			input: "1\nmatch\n3\n\n5\nmatch\n7",
			want:  "match\n3\n--\nmatch\n7\n",
		},
		{
			name:     "file input",
			args:     []string{"grep", "match", "test.txt"},
			fileData: "1\nmatch\n3",
			want:     "match\n",
		},
		{
			name:  "empty input",
			args:  []string{"grep", "match"},
			input: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var tempFile string
			if tt.fileData != "" {
				file, err := os.CreateTemp("", "grep_test_*.txt")
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
					if arg == "test.txt" {
						tt.args[i] = file.Name()
					}
				}
				tempFile = file.Name()
			}
			cmd := exec.Command("go", append([]string{"run", "grep.go"}, tt.args[1:]...)...)
			if tt.input != "" {
				cmd.Stdin = strings.NewReader(tt.input)
			}
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			exitCode := 0
			if err != nil {
				if exitError, ok := err.(*exec.ExitError); ok {
					exitCode = exitError.ExitCode()
				} else {
					t.Fatalf("failed to run command: %v", err)
				}
			}
			if exitCode != tt.wantCode {
				t.Errorf("expected exit code %d, got %d, stderr: %s", tt.wantCode, exitCode, stderr.String())
			}
			if tt.want != "" {
				if stdout.String() != tt.want {
					t.Errorf("got %q, want %q", stdout.String(), tt.want)
				}
			}
			if tempFile != "" {
				os.Remove(tempFile)
			}
		})
	}
}
