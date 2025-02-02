package misc

import (
	"strings"
	"unicode"
)

func FormaterName(s string) string {
	sl := strings.ToLower(s)
	sf := strings.Fields(sl)

	var sb strings.Builder

	for _, v := range sf {
		if len(v) > 0 {
			firstLetter := unicode.ToUpper(rune(v[0]))
			rest := v[1:]
			sb.WriteRune(firstLetter)
			sb.WriteString(rest)
			sb.WriteRune(' ')
		}
	}

	return strings.TrimSpace(sb.String())
}
