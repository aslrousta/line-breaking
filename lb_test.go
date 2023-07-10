package lb_test

import (
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"

	lb "github.com/aslrousta/line-breaking"
)

var (
	alice = "Alice was beginning to get very tired of sitting by her sister" +
		" on the bank, and of having nothing to do: once or twice she had" +
		" peeped into the book her sister was reading, but it had no pictures" +
		" or conversations in it, 'and what is the use of a book,' thought" +
		" Alice 'without pictures or conversation?'"
)

type Word string

func (w Word) Direction() lb.Direction {
	r, _ := utf8.DecodeRuneInString(string(w))
	if unicode.Is(unicode.Arabic, r) {
		return lb.RightToLeft
	}
	return lb.LeftToRight
}

func (w Word) Width() float32 {
	return float32(utf8.RuneCountInString(string(w)))
}

func BenchmarkGreedy(b *testing.B) {
	para := splitWords(alice)
	for i := 0; i < b.N; i++ {
		lb.Greedy(para, &lb.Options{
			TextWidth:     50,
			TextDirection: lb.LeftToRight,
			GlueWidth:     1,
			GlueExpand:    1,
		})
	}
}

func BenchmarkKnuthPlass(b *testing.B) {
	para := splitWords(alice)
	for i := 0; i < b.N; i++ {
		lb.KnuthPlass(para, &lb.Options{
			TextWidth:     50,
			TextDirection: lb.LeftToRight,
			GlueWidth:     1,
			GlueExpand:    1,
		})
	}
}

func splitWords(text string) []lb.Box {
	words := strings.Split(text, " ")
	result := make([]lb.Box, 0, len(words))
	for _, w := range words {
		result = append(result, Word(w))
	}
	return result
}
