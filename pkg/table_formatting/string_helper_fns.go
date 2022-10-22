package table_formatting

import (
	"strings"
	"unicode"

	"github.com/zeus-fyi/tables-to-go/pkg/settings"
	"github.com/zeus-fyi/tables-to-go/pkg/tagger"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	Taggers tagger.Tagger
	Caser   = cases.Title(language.English, cases.NoLower)
	// Initialisms some strings for idiomatic go in column names
	// see https://github.com/golang/go/wiki/CodeReviewComments#initialisms
	Initialisms = []string{"ID", "JSON", "XML", "HTTP", "URL"}
)

func CamelCaseString(s string) string {
	if s == "" {
		return s
	}

	splitted := strings.Split(s, "_")

	if len(splitted) == 1 {
		return Caser.String(s)
	}

	var cc string
	for _, part := range splitted {
		cc += Caser.String(strings.ToLower(part))
	}
	return cc
}

func GetNullType(settings *settings.Settings, primitive string, sql string) string {
	if settings.IsNullTypeSQL() {
		return sql
	}
	return primitive
}

func ToInitialisms(s string) string {
	for _, substr := range Initialisms {
		idx := IndexCaseInsensitive(s, substr)
		if idx == -1 {
			continue
		}
		toReplace := s[idx : idx+len(substr)]
		s = strings.ReplaceAll(s, toReplace, substr)
	}
	return s
}

func IndexCaseInsensitive(s, substr string) int {
	s, substr = strings.ToLower(s), strings.ToLower(substr)
	return strings.Index(s, substr)
}

// ValidVariableName checks for the existence of any characters
// outside of Unicode letters, numbers and underscore.
func ValidVariableName(s string) bool {
	for _, r := range s {
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_') {
			return false
		}
	}
	return true
}

// ReplaceSpace swaps any Unicode space characters for underscores
// to create valid Go identifiers
func ReplaceSpace(r rune) rune {
	if unicode.IsSpace(r) || r == '\u200B' {
		return '_'
	}
	return r
}
