package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
)

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

type grepConfig struct {
	After      int
	Before     int
	Context    int
	Count      bool
	IgnoreCase bool
	Inverse    bool
	Fixed      bool
	LineNum    bool
	Pattern    string
}

func main() {
	err := Run(os.Args, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		fmt.Println("Error: %w", err)
		os.Exit(1)
	}
}

// Run executes the grep command with given arguments, reading from
// stdin or file, writing to stdout, and reporting errors to stderr.
//
// args must include the command name (e.g., "grep"), followed by
// flags and an optional file path.
//
// Returns nil on success, or an error if execution fails.
func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	flagSet := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flagSet.SetOutput(stderr)
	var (
		AFlag = flagSet.Int("A", 0, "Print N lines after match")
		BFlag = flagSet.Int("B", 0, "Print N lines before match")
		CFlag = flagSet.Int("C", 0, "Print number lines around match")
		cFlag = flagSet.Bool("c", false, "Only count matching lines")
		iFlag = flagSet.Bool("i", false, "Ignore case")
		vFlag = flagSet.Bool("v", false, "Invert match")
		FFlag = flagSet.Bool("F", false, "Equal to template")
		nFlag = flagSet.Bool("n", false, "Print line numbers")
	)

	if err := flagSet.Parse(args[1:]); err != nil {
		return err
	}
	template := flagSet.Args()[0]

	config := &grepConfig{
		After:      *AFlag,
		Before:     *BFlag,
		Context:    *CFlag,
		Count:      *cFlag,
		IgnoreCase: *iFlag,
		Inverse:    *vFlag,
		Fixed:      *FFlag,
		LineNum:    *nFlag,
		Pattern:    template,
	}
	var source io.Reader = stdin
	if len(flagSet.Args()) > 1 {
		var err error
		file, err := os.Open(flagSet.Arg(1))
		if err != nil {
			return fmt.Errorf("error opening file: %w", err)
		}
		defer file.Close()
		source = file
	}

	scanner := bufio.NewScanner(source)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if *cFlag {
		count := 0
		for _, line := range lines {
			if match(line, config) {
				count++
			}
		}
		fmt.Fprintln(stdout, count)
		return nil
	}
	after := max(*AFlag, *CFlag)
	before := max(*BFlag, *CFlag)

	res := make(map[int]struct{})
	for i, e := range lines {
		if match(e, config) {
			res[i] = struct{}{}

			for j := 1; j <= before; j++ {
				idx := i - j
				if idx >= 0 {
					res[idx] = struct{}{}
				}
			}

			for j := 1; j <= after; j++ {
				idx := i + j
				if idx < len(lines) {
					res[idx] = struct{}{}
				}
			}

		}

	}
	var keys []int
	for i := range res {
		keys = append(keys, i)
	}
	sort.Ints(keys)

	if *nFlag {
		for i, idx := range keys {
			if i > 0 && idx > keys[i-1]+1 && (before > 0 || after > 0) {
				fmt.Fprintln(stdout, "--")
			}
			fmt.Fprintf(stdout, "%d:%s\n", idx+1, lines[idx])
		}
		return nil
	}
	for i, idx := range keys {
		if i > 0 && idx > keys[i-1]+1 && (before > 0 || after > 0) {
			fmt.Fprintln(stdout, "--")
		}
		fmt.Fprintf(stdout, "%s\n", lines[idx])
	}

	return nil
}

func match(line string, config *grepConfig) bool {
	var found bool
	if config.Fixed {
		a, b := line, config.Pattern
		if config.IgnoreCase {
			a, b = strings.ToLower(a), strings.ToLower(b)
		}
		found = strings.Contains(a, b)
	} else {
		pat := config.Pattern
		if config.IgnoreCase {
			pat = "(?i)" + pat
		}
		found, _ = regexp.MatchString(pat, line)
	}
	if config.Inverse {
		found = !found
	}
	return found
}
