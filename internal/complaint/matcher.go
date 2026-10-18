package complaint

import (
	"strings"
	"unicode"
)

// normalize lower-cases text, replaces everything that is not a letter or a
// digit with a single space and returns the words joined by one space.
// "Wi-Fi не работает!" becomes "wi fi не работает".
func normalize(text string) string {
	words := strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	return strings.Join(words, " ")
}

// pad surrounds normalized text with spaces so that every word, including the
// first and the last one, has a space before it and after it.
func pad(normalized string) string {
	return " " + normalized + " "
}

// keyword is a compiled keyword: the needle to look for inside padded text.
type keyword struct {
	needle string
}

// keywordSet is a list of compiled keywords without duplicates.
type keywordSet []keyword

// newKeywordSet compiles raw keywords. A trailing "$" means "whole word only".
func newKeywordSet(raw []string) keywordSet {
	seen := make(map[string]struct{}, len(raw))
	set := make(keywordSet, 0, len(raw))
	for _, kw := range raw {
		exact := strings.HasSuffix(kw, "$")
		kw = normalize(strings.TrimSuffix(kw, "$"))
		if kw == "" {
			continue
		}
		needle := " " + kw // the match must start at the beginning of a word
		if exact {
			needle += " " // ...and must end at the end of a word
		}
		if _, dup := seen[needle]; dup {
			continue
		}
		seen[needle] = struct{}{}
		set = append(set, keyword{needle: needle})
	}
	return set
}

// count returns how many distinct keywords occur in the padded text.
func (ks keywordSet) count(padded string) int {
	n := 0
	for _, kw := range ks {
		if strings.Contains(padded, kw.needle) {
			n++
		}
	}
	return n
}

// any reports whether at least one keyword occurs in the padded text.
func (ks keywordSet) any(padded string) bool {
	for _, kw := range ks {
		if strings.Contains(padded, kw.needle) {
			return true
		}
	}
	return false
}
