package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

type lineSorted struct {
	lines   []string
	col     int
	numeric bool
	reverse bool
	month   bool
	ignore  bool
	suffix  bool
}

func (l *lineSorted) Len() int {
	return len(l.lines)
}

func (l *lineSorted) Less(i, j int) bool {
	a := l.lines[i]
	b := l.lines[j]
	if l.col >= 0 {
		partsA := strings.Split(a, "\t")
		partsB := strings.Split(b, "\t")
		if l.col < len(partsA) {
			a = partsA[l.col]
		} else {
			a = ""
		}
		if l.col < len(partsB) {
			b = partsB[l.col]
		} else {
			b = ""
		}
	}
	if l.ignore {
		a = strings.TrimRight(a, " ")
		b = strings.TrimRight(b, " ")
	}
	res := a < b
	if l.numeric {
		aNum, aErr := strconv.Atoi(a)
		bNum, bErr := strconv.Atoi(b)
		if aErr == nil && bErr == nil {
			if aNum == bNum {
				res = a < b
			} else {
				res = aNum < bNum
			}

		}
	}
	if l.month {
		aMonth, aOk := monthMap[a]
		bMonth, bOk := monthMap[b]
		if aOk && bOk {
			res = aMonth < bMonth
		}
	}
	if l.suffix {
		aNum, aErr := parseSuffix(a)
		bNum, bErr := parseSuffix(b)
		if aErr == nil && bErr == nil {
			res = aNum < bNum
		}
	}
	if l.reverse {
		return !res
	}
	return res
}

func (l *lineSorted) Swap(i, j int) {
	l.lines[i], l.lines[j] = l.lines[j], l.lines[i]
}

var monthMap = map[string]int{
	"Jan": 1, "Feb": 2, "Mar": 3,
	"Apr": 4, "May": 5, "Jun": 6,
	"Jul": 7, "Aug": 8, "Sep": 9,
	"Oct": 10, "Nov": 11, "Dec": 12,
}

func main() {
	if err := Run(os.Args, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Run executes the sort command with given arguments, reading from
// stdin, writing to stdout, and reporting errors to stderr.
//
// args must include the command name (e.g., "sort"), followed by
// flags and an optional file path.
//
// Returns nil on success, or an error if execution fails.
func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	flagSet := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flagSet.SetOutput(stderr)

	var (
		rFlag   = flagSet.Bool("r", false, "Sort in reverse order")
		uFlag   = flagSet.Bool("u", false, "Return only unique")
		nFlag   = flagSet.Bool("n", false, "Sort strings like numbers")
		cFlag   = flagSet.Bool("c", false, "Check if the lines are sorted")
		colFlag = flagSet.Int("k", 0, "Sort strings by k column")
		bFlag   = flagSet.Bool("b", false, "Ignore trailing blanks")
		mFlag   = flagSet.Bool("M", false, "Sort by months")
		hFlag   = flagSet.Bool("h", false, "Sorting based on suffixes")
	)

	if err := flagSet.Parse(parseFlags(args)[1:]); err != nil {
		return err
	}
	source := stdin
	if len(flagSet.Args()) > 0 {
		var err error
		file, err := os.Open(flagSet.Args()[len(flagSet.Args())-1])
		if err != nil {
			return fmt.Errorf("error opening file %w", err)
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
		return fmt.Errorf("error scanner %w", err)
	}

	if *uFlag {
		lines = unique(lines)
	}
	if *mFlag && *nFlag || *hFlag && *nFlag || *mFlag && *hFlag {
		return fmt.Errorf("Invalid flags combination")
	}
	sorter := &lineSorted{
		lines:   lines,
		col:     *colFlag - 1,
		numeric: *nFlag,
		reverse: *rFlag,
		month:   *mFlag,
		ignore:  *bFlag,
		suffix:  *hFlag,
	}

	if *cFlag {
		if !sort.IsSorted(sorter) {
			return fmt.Errorf("The lines are not sorted")
		}
		return nil
	}

	sort.Stable(sorter)
	for _, line := range lines {
		fmt.Fprintln(stdout, line)
	}
	return nil
}

func parseSuffix(s string) (int64, error) {
	var multiplier int64 = 1
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return 0, fmt.Errorf("empty string")
	}
	suffix := s[len(s)-1]
	switch suffix {
	case 'K':
		multiplier = 1024
	case 'M':
		multiplier = 1024 * 1024
	case 'G':
		multiplier = 1024 * 1024 * 1024
	case 'T':
		multiplier = 1024 * 1024 * 1024 * 1024
	}
	if multiplier != 1 {
		s = s[:len(s)-1]
	}
	num, err := strconv.ParseInt(s, 10, 64)
	return num * multiplier, err
}

func unique(lines []string) []string {
	unique := make(map[string]struct{})
	res := make([]string, 0, len(lines))
	for _, e := range lines {
		if _, ok := unique[e]; !ok {
			unique[e] = struct{}{}
			res = append(res, e)
		}
	}
	return res
}

func parseFlags(args []string) []string {
	var parsed []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") && len(arg) > 2 {
			if strings.HasPrefix(arg, "-k") {
				parsed = append(parsed, "-k", arg[2:])
			} else {
				for _, f := range arg[1:] {
					parsed = append(parsed, "-"+string(f))
				}
			}
		} else {
			parsed = append(parsed, arg)
		}
	}
	return parsed
}
