package textinvertedindex

import (
	"strings"

	snowballeng "github.com/kljensen/snowball/english"
)

// lowercaseFilter returns a slice of tokens normalized to lower case.
func (i *TextInvertedIndex[T]) lowercaseFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = strings.ToLower(token)
	}
	return r
}

// TODO: adopt more languages
var stopWordsEN = []string{"a", "an", "and", "are", "as", "at", "be", "but", "by", "for", "if", "in", "into", "is", "it",
	"no", "not", "of", "on", "or", "such", "that", "the", "their", "then", "than", "there", "these", "they",
	"this", "to", "was", "will", "with", "those"}

// stopwordFilter returns a slice of tokens with stop words removed.
func (i *TextInvertedIndex[T]) stopwordFilter(tokens []string) []string {
	r := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if !i.stopWords.Has(token) {
			r = append(r, token)
		}
	}
	return r
}

// stemmerFilter returns a slice of stemmed tokens.
func (i *TextInvertedIndex[T]) stemmerFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = snowballeng.Stem(token, false)
	}
	return r
}
