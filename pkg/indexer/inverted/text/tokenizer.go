package textinvertedindex

import (
	"strings"
	"unicode"
)

// tokenize returns a slice of tokens for the given text.
func (i *TextInvertedIndex[T]) tokenize(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		// Split on any character that is not a letter or a number.
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

// analyze analyzes the text and returns a slice of tokens.
func (i *TextInvertedIndex[T]) analyze(text string) []string {
	tokens := i.tokenize(text)
	tokens = i.lowercaseFilter(tokens)
	tokens = i.stopwordFilter(tokens)
	tokens = i.stemmerFilter(tokens)
	return tokens
}
