package complaint

import (
	"strings"
	"unicode"
)

func normalize(text string) string {
	words := strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	return strings.Join(words, " ")
}

func pad(normalized string) string {
	return " " + normalized + " "
}

type keyword struct {
	needle string
}

type keywordSet []keyword

func newKeywordSet(raw []string) keywordSet {
	seen := make(map[string]struct{}, len(raw))
	set := make(keywordSet, 0, len(raw))
	for _, kw := range raw {
		exact := strings.HasSuffix(kw, "$")
		kw = normalize(strings.TrimSuffix(kw, "$"))
		if kw == "" {
			continue
		}
		needle := " " + kw
		if exact {
			needle += " "
		}
		if _, dup := seen[needle]; dup {
			continue
		}
		seen[needle] = struct{}{}
		set = append(set, keyword{needle: needle})
	}
	return set
}

func (ks keywordSet) count(padded string) int {
	n := 0
	for _, kw := range ks {
		if strings.Contains(padded, kw.needle) {
			n++
		}
	}
	return n
}

func (ks keywordSet) any(padded string) bool {
	for _, kw := range ks {
		if strings.Contains(padded, kw.needle) {
			return true
		}
	}
	return false
}
