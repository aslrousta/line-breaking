package lb

// Direction is the data-type for word direction.
type Direction int

// List of word direction values.
const (
	LeftToRight Direction = iota
	RightToLeft
)

// Box is a solid text element (typically, a word).
type Box interface {

	// Direction returns the writing direction of the box.
	Direction() Direction

	// Width returns the extent of the box within the text.
	Width() float32
}

// Line is a series of boxes that fit in a line of text.
type Line struct {

	// Boxes is the slice of boxes within the line.
	Boxes []Box

	// GlueWidth is the extent of the glues between the boxes.
	GlueWidth float32
}

// Options contains several options passed to the line-breaking algorithm.
type Options struct {

	// TextWidth is the extent of a line of text.
	TextWidth float32

	// TextDirection is the writing direction of the text.
	TextDirection Direction

	// GlueWidth is the normal extent of the glues between the boxes.
	GlueWidth float32

	// GlueShrink is the extent that a glue can shrink to.
	GlueShrink float32

	// GlueExpand is the extent that a glue can expand to.
	GlueExpand float32
}

// Greedy is a greedy but fast line-breaking algorithm that tries to fit as much
// boxes as possible in each and every line.
func Greedy(para []Box, opt *Options) (lines []*Line) {
	var (
		minGlueWidth = opt.GlueWidth - opt.GlueShrink
		maxGlueWidth = opt.GlueWidth + opt.GlueExpand
		line         = &Line{GlueWidth: opt.GlueWidth}
		lineWidth    = float32(0)
		boxesWidth   = float32(0)
		numGlues     = 0
	)
	for _, box := range para {
		boxWidth := box.Width()
		if len(line.Boxes) == 0 {
			line.Boxes = append(line.Boxes, box)
			lineWidth += boxWidth
			boxesWidth += boxWidth
		} else if lineWidth+minGlueWidth+boxWidth <= opt.TextWidth {
			line.Boxes = append(line.Boxes, box)
			lineWidth += minGlueWidth + boxWidth
			boxesWidth += boxWidth
			numGlues++
		} else {
			if numGlues > 0 && boxesWidth < opt.TextWidth {
				line.GlueWidth = (opt.TextWidth - boxesWidth) / float32(numGlues)
				if line.GlueWidth > maxGlueWidth {
					line.GlueWidth = maxGlueWidth
				}
			}
			bidi(line, opt.TextDirection)
			lines = append(lines, line)
			line = &Line{GlueWidth: opt.GlueWidth, Boxes: []Box{box}}
			lineWidth = boxWidth
			boxesWidth = boxWidth
			numGlues = 0
		}
	}
	if len(line.Boxes) > 0 {
		lines = append(lines, line)
	}
	return lines
}

// KnuthPlass is a relatively slow line-breaking algorithm but gives the best
// result aesthetically.
func KnuthPlass(para []Box, opt *Options) (lines []*Line) {
	ends, _ := knuthPlass(para, opt, 0)
	i := 0
	for _, end := range ends {
		var (
			line       = &Line{GlueWidth: opt.GlueWidth}
			boxesWidth = float32(0)
		)
		for i < end {
			line.Boxes = append(line.Boxes, para[i])
			boxesWidth += para[i].Width()
			i++
		}
		if numGlues := len(line.Boxes) - 1; end < len(para) && numGlues > 0 {
			line.GlueWidth = (opt.TextWidth - boxesWidth) / float32(numGlues)
		}
		bidi(line, opt.TextDirection)
		lines = append(lines, line)
	}
	return lines
}

func knuthPlass(para []Box, opt *Options, start int) (ends []int, badness float32) {
	var end int
	badness = float32(10_000)
	for _, bp := range breakRange(para, opt, start) {
		nextBadness := float32(0)
		if bp.Index < len(para) {
			ends, nextBadness = knuthPlass(para, opt, bp.Index)
		}
		if b := bp.Badness + nextBadness; b < badness {
			badness = bp.Badness + nextBadness
			end = bp.Index
		}
	}
	return append([]int{end}, ends...), badness
}

type breakPoint struct {
	Index   int
	Badness float32
}

func breakRange(para []Box, opt *Options, start int) (bps []breakPoint) {
	var (
		first, last  int
		minGlueWidth = opt.GlueWidth - opt.GlueShrink
		maxGlueWidth = opt.GlueWidth + opt.GlueExpand
		minLineWidth = float32(0)
		maxLineWidth = float32(0)
		boxesWidth   = float32(0)
		numGlues     = 0
	)
	for first = start; first < len(para); first++ {
		boxWidth := para[first].Width()
		if first == start {
			minLineWidth += boxWidth
			maxLineWidth += boxWidth
			boxesWidth += boxWidth
		} else if maxLineWidth+maxGlueWidth+boxWidth <= opt.TextWidth {
			minLineWidth += minGlueWidth + boxWidth
			maxLineWidth += maxGlueWidth + boxWidth
			boxesWidth += boxWidth
			numGlues++
		} else {
			badness := float32(0)
			if numGlues > 0 {
				diff := opt.GlueWidth - (opt.TextWidth-boxesWidth)/float32(numGlues)
				badness = diff * diff
			}
			bps = append(bps, breakPoint{Index: first, Badness: badness})
			break
		}
	}
	if len(bps) == 0 {
		bps = append(bps, breakPoint{Index: first})
		return
	}
	for last = first; last < len(para); last++ {
		boxWidth := para[last].Width()
		if minLineWidth+minGlueWidth+boxWidth > opt.TextWidth {
			break
		}
		minLineWidth += minGlueWidth + boxWidth
		badness := float32(0)
		if numGlues > 0 && last != len(para)-1 {
			diff := opt.GlueWidth - (opt.TextWidth-boxesWidth)/float32(numGlues)
			badness = diff * diff
		}
		bps = append(bps, breakPoint{Index: last + 1, Badness: badness})
	}
	return bps
}

func bidi(line *Line, textDirection Direction) {
	var (
		result []Box
		stack  boxStack
		box    Box
	)
	result = make([]Box, 0, len(line.Boxes))
	for _, b := range line.Boxes {
		if b.Direction() != textDirection {
			stack = pushBox(stack, b)
		} else {
			for !stack.Empty() {
				stack, box = popBox(stack)
				result = append(result, box)
			}
			result = append(result, b)
		}
	}
	for !stack.Empty() {
		stack, box = popBox(stack)
		result = append(result, box)
	}
	line.Boxes = result
}

type boxStack []Box

func (s boxStack) Empty() bool {
	return len(s) == 0
}

func pushBox(s boxStack, b Box) boxStack {
	return append(s, b)
}

func popBox(s boxStack) (boxStack, Box) {
	return s[:len(s)-1], s[len(s)-1]
}
