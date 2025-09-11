package utils

import (
	"strings"
)

func CleanString(s string, remove ...string) string {
	for _, r := range remove {
		s = strings.ReplaceAll(s, r, "")
	}
	return s
}
