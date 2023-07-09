package lb_test

import (
	"fmt"
	"strings"
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

func ExampleGreedy() {
	para := splitWords(alice)
	lines := lb.Greedy(para, &lb.Options{
		TextWidth:     40,
		TextDirection: lb.LeftToRight,
		GlueWidth:     1,
		GlueExpand:    1,
	})
	printLines(lines)

	// Output:
	// Alice was beginning to get very tired of
	// sitting by her sister on the bank, and
	// of having nothing to do: once or twice
	// she had peeped into the book her sister
	// was reading, but it had no pictures or
	// conversations in it, 'and what is the
	// use of a book,' thought Alice 'without
	// pictures or conversation?'
}

func ExampleKnuthPlass() {
	para := splitWords(alice)
	lines := lb.KnuthPlass(para, &lb.Options{
		TextWidth:     40,
		TextDirection: lb.LeftToRight,
		GlueWidth:     1,
		GlueExpand:    1,
	})
	printLines(lines)

	// Output:
	// Alice was beginning to get very tired of
	// sitting  by  her  sister  on  the  bank,
	// and of having nothing to do: once or
	// twice she had peeped into the book her
	// sister was reading, but it had no pictures
	// or conversations in it, 'and what is
	// the  use  of  a  book,'  thought  Alice
	// 'without pictures or conversation?'
}

func splitWords(text string) []lb.Box {
	words := strings.Split(text, " ")
	result := make([]lb.Box, 0, len(words))
	for _, w := range words {
		result = append(result, Word(w))
	}
	return result
}

func printLines(lines []*lb.Line) {
	for _, l := range lines {
		glueWidth := int(l.GlueWidth)
		if glueWidth < 1 {
			glueWidth = 1
		}
		for i, b := range l.Boxes {
			w := b.(Word)
			fmt.Printf("%s", w)
			if i != len(l.Boxes)-1 {
				fmt.Print(strings.Repeat(" ", glueWidth))
			}
		}
		fmt.Println()
	}
}
