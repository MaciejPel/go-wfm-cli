package text

import (
	"strings"
	"unicode"

	"github.com/texttheater/golang-levenshtein/levenshtein"
)

func normalize(input string) string {
	return strings.ToLower(strings.ReplaceAll(input, " ", ""))
}

func similarity(a, b string) float64 {
	aNorm := normalize(a)
	bNorm := normalize(b)
	dist := levenshtein.DistanceForStrings([]rune(aNorm), []rune(bNorm), levenshtein.DefaultOptions)
	maxLen := max(len(aNorm), len(bNorm))

	return float64(1) - float64(dist)/float64(maxLen)
}

func BestMatch(input string, candidates []string) string {
	result := ""
	score := -1.0

	for _, c := range candidates {
		s := similarity(input, c)
		if s > score {
			score = s
			result = c
		}
	}

	return result
}

func SeparateByUppercase(in string) string {
	var out strings.Builder
	for i, v := range in {
		if unicode.IsUpper(v) && i > 0 {
			out.WriteString(" " + string(v))
		} else {
			out.WriteString(string(v))
		}
	}
	return out.String()
}

func ItemCamelToSnake(in string) string {
	return strings.ReplaceAll(strings.ToLower(in), " ", "_")
}
