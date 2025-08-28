package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	err := Run(os.Args, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		fmt.Println("Error: %w", err)
		os.Exit(1)
	}
}

// Run executes the cut command with given arguments, reading from
// stdin or file, writing to stdout, and reporting errors to stderr.
//
// args must include the command name (e.g., "cut"), followed by
// flags and an optional file path.
//
// Returns nil on success, or an error if execution fails.
func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	flagSet := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flagSet.SetOutput(stderr)
	var (
		fFlag = flagSet.String("f", "", "field to output, e.g. 1,3-5")
		dFlag = flagSet.String("d", "\t", "delimiter (default: tab)")
		sFlag = flagSet.Bool("s", false, "only lines with delimiter")
	)

	if err := flagSet.Parse(args[1:]); err != nil {
		return err
	}
	fields, err := parseFields(*fFlag)
	if err != nil {
		return err
	}
	source := stdin
	if flagSet.NArg() > 0 {
		file, err := os.Open(flagSet.Arg(0))
		if err != nil {
			return fmt.Errorf("error opening file: %w", err)
		}
		defer file.Close()
		source = file
	}

	scanner := bufio.NewScanner(source)
	for scanner.Scan() {
		line := scanner.Text()

		if *sFlag && !strings.Contains(line, *dFlag) {
			continue
		}
		parts := strings.Split(line, *dFlag)
		var res []string
		for _, f := range fields {
			if f > 0 && f <= len(parts) {
				res = append(res, parts[f-1])
			}
		}
		if len(res) > 0 {
			fmt.Fprintln(stdout, strings.Join(res, *dFlag))
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func parseFields(s string) ([]int, error) {
	var fields []int
	parts := strings.Split(s, ",")
	for _, part := range parts {
		if strings.Contains(part, "-") {
			partsOfPart := strings.Split(part, "-")
			l, err1 := strconv.Atoi(partsOfPart[0])
			r, err2 := strconv.Atoi(partsOfPart[1])
			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("invalid range: %s", part)
			}
			for i := l; i <= r; i++ {
				fields = append(fields, i)
			}
		} else {
			num, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid field: %s", part)
			}
			fields = append(fields, num)
		}
	}
	return fields, nil
}
