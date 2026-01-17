package canonization

import (
	"strings"
	"unicode"
)

type Canonizer interface {
	Canonize(input []string) string
	CanonizeString(input string) string
}

type DefaultCanonizer struct{}

func (c *DefaultCanonizer) Canonize(input []string) string {
	canonized := strings.Builder{}
	for _, word := range input {
		canonized.WriteString(c.CanonizeString(word))
	}
	return canonized.String()
}

func (c *DefaultCanonizer) CanonizeString(input string) string {
	canonized := strings.Builder{}
	for _, char := range input {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			canonized.WriteRune(unicode.ToLower(rune(char)))
		}
	}
	return canonized.String()
}
