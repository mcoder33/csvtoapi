package main

import (
	"strings"
	"unicode"
)

func CleanString(s string, excl rune) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == excl {
			return r
		}
		return -1
	}, s)
}

func CleanStringFrom(s string, remove ...string) string {
	for _, r := range remove {
		s = strings.ReplaceAll(s, r, "")
	}
	return s
}
