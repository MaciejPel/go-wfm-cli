package text

import (
	"strings"
	"unicode"

	"github.com/texttheater/golang-levenshtein/levenshtein"
)

func normalize(input string) string {
	return strings.ToLower(strings.ReplaceAll(input, " ", ""))
}

func similarity(a, b string) int {
	aNorm := normalize(a)
	bNorm := normalize(b)
	dist := levenshtein.DistanceForStrings([]rune(aNorm), []rune(bNorm), levenshtein.DefaultOptions)
	maxLen := max(len(aNorm), len(bNorm))

	return maxLen - dist
}

func BestMatch(input string, candidates []string) string {
	result := ""
	score := -1

	for _, c := range candidates {
		s := similarity(input, c)
		if s > score {
			score = s
			result = c
		}
	}

	return result
}

func SeparateByUppercase(input string) string {
	var out strings.Builder
	for i, v := range input {
		if unicode.IsUpper(v) && i > 0 {
			out.WriteString(" " + string(v))
		} else {
			out.WriteString(string(v))
		}
	}
	return out.String()
}
