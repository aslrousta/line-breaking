# Line-breaking

This package provides a Go implementation of two line-breaking algorithms with support for bi-directional texts. However, it abstracts text elements (normally, words) with a `Box` interface that has a `Width` and `Direction` instead of plain `string`s, for more flexibility.

An example usage would be like the following:

```go
package main

import (
    "strings"
    lb "github.com/aslrousta/line-breaking"
)

var paragraph = "Alice was beginning ..."

type Word string
func (w Word) Direction() lb.Direction { return lb.LeftToRight }
func (w Word) Width() float32 {
    // Compute the extent of the word in a desired font-face.
}

func main() {

    // Convert words in the paragraph into boxes.
    words := strings.Split(paragraph, " ")
    boxes := make([]lb.Box, 0, len(words))
    for _, w := range words {
        boxes = append(boxes, Word(w))
    }

    spaceWidth := /* Compute the extent of a `space` character */

    // Do the line-breaking using Knuth-Plass algorithm.
    lines := lb.KnuthPlass(boxes, &lb.Options{
        TextWidth:     60,
        TextDirection: lb.LeftToRight,
        GlueWidth:     spaceWidth,
        GlueShrink:    spaceWidth / 5, /* 20% shrink */
        GlueExpand:    spaceWidth / 3, /* 33% expand */
    })

    renderLines(lines)
    ...
}
```

See the `example` folder for a more practical sample code.

## Algorithms

Two `Greedy` and `Knuth-Plass` line-breaking algorithms are provided. The greedy approach is a fast algorithm that tries to fit as boxes as possible within a line. On the other hand, the Knuth-Plass algorithm is relatively slow but gives better results

## Copyright

This package is distributed under the MIT license. See the `LICENSE` file for more information.
