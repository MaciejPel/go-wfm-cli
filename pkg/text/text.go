package text

import (
	"strings"
	"unicode"

	"github.com/MaciejPel/go-wfm-cli/pkg/market"
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

func BestMatch(input string, candidates []market.ItemValue) market.ItemValue {
	result := market.ItemValue{}
	score := -1.0

	for _, item := range candidates {
		s := similarity(input, item.Name)
		if s > score {
			score = s
			result = item
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

func ItemNameToSlug(in string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(in), " ", "_"), "&", "and")
}
