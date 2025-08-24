package main

import (
	"bytes"
	"os"
	"testing"
)

func TestSort(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		input    string
		fileData string
		want     string
		wantErr  bool
	}{

		{
			name:  "simple sort",
			args:  []string{"sort"},
			input: "banana\napple\ncherry\n",
			want:  "apple\nbanana\ncherry\n",
		},
		{
			name:  "reverse sort",
			args:  []string{"sort", "-r"},
			input: "apple\nbanana\ncherry\n",
			want:  "cherry\nbanana\napple\n",
		},
		{
			name:  "numeric sort",
			args:  []string{"sort", "-n"},
			input: "10\n2\n5\n",
			want:  "2\n5\n10\n",
		},
		{
			name:  "numeric reverse",
			args:  []string{"sort", "-nr"},
			input: "10\n2\n5\n",
			want:  "10\n5\n2\n",
		},
		{
			name:  "unique",
			args:  []string{"sort", "-u"},
			input: "apple\napple\nbanana\napple\n",
			want:  "apple\nbanana\n",
		},
		{
			name:    "check sorted",
			args:    []string{"sort", "-c"},
			input:   "apple\nbanana\ncherry\n",
			wantErr: false,
		},
		{
			name:    "check not sorted",
			args:    []string{"sort", "-c"},
			input:   "banana\napple\n",
			wantErr: true,
		},
		{
			name:  "sort by column numeric",
			args:  []string{"sort", "-k2", "-n"},
			input: "Alice\t30\nBob\t25\nCharlie\t35\n",
			want:  "Bob\t25\nAlice\t30\nCharlie\t35\n",
		},
		{
			name:  "month sort",
			args:  []string{"sort", "-M"},
			input: "Mar\nJan\nFeb\n",
			want:  "Jan\nFeb\nMar\n",
		},
		{
			name:  "human sort",
			args:  []string{"sort", "-h"},
			input: "1M\n512\n10K\n",
			want:  "512\n10K\n1M\n",
		},
		{
			name:  "ignore trailing blanks reverse",
			args:  []string{"sort", "-b", "-r"},
			input: "apple   \napple\nbanana\n",
			want:  "banana\napple\napple   \n",
		},
		{
			name:     "file input",
			args:     []string{"sort", "test.txt"},
			fileData: "3\n1\n2\n",
			want:     "1\n2\n3\n",
		},
		{
			name:  "empty input",
			args:  []string{"sort"},
			input: "",
			want:  "",
		},
		{
			name:  "one line",
			args:  []string{"sort"},
			input: "hello\n",
			want:  "hello\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdin bytes.Buffer
			stdin.WriteString(tt.input)

			var stdout, stderr bytes.Buffer
			var tempFile *os.File
			if tt.fileData != "" {
				var err error
				tempFile, err = os.CreateTemp("", "sort_test_*.txt")
				if err != nil {
					t.Fatal(err)
				}
				defer os.Remove(tempFile.Name())
				defer tempFile.Close()

				if _, err := tempFile.WriteString(tt.fileData); err != nil {
					t.Fatal(err)
				}
				if err := tempFile.Close(); err != nil {
					t.Fatal(err)
				}

				args := make([]string, len(tt.args))
				copy(args, tt.args)
				for i, arg := range args {
					if arg == "test.txt" {
						args[i] = tempFile.Name()
					}
				}
				tt.args = args
			}

			err := Run(tt.args, &stdin, &stdout, &stderr)
			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.want != "" {
				if stdout.String() != tt.want {
					t.Errorf("got %q, want %q", stdout.String(), tt.want)
				}
			}
		})
	}
}
