package main

import (
	"fmt"
	"strings"
	"unicode"
)

func unpack(s string) (string, error) {
	var builder strings.Builder
	var escape bool
	var prev rune
	for _, e := range s {
		if e == '\\' {
			if escape {
				return "", fmt.Errorf("invalid escape character after escape")
			}
			escape = true
		} else if unicode.IsDigit(e) {
			if escape {
				builder.WriteRune(e)
				escape = false
				prev = e
			} else {
				if prev != 0 {
					for i := 0; i < int(e-'0'-1); i++ {
						builder.WriteRune(prev)
					}
				} else {
					return "", fmt.Errorf("digit without preceding character")
				}
			}
		} else {
			prev = e
			builder.WriteRune(e)
			escape = false
		}

	}
	if escape {
		return "", fmt.Errorf("invalid escape at end of string")
	}
	return builder.String(), nil
}
