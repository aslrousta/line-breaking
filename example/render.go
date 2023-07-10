package main

import (
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"strings"

	lb "github.com/aslrousta/line-breaking"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type Word struct {
	text  string
	width float32
}

func (w Word) Direction() lb.Direction { return lb.LeftToRight }
func (w Word) Text() string            { return w.text }
func (w Word) Width() float32          { return w.width }

func main() {
	text, err := os.ReadFile("sample.txt")
	if err != nil {
		log.Fatal(err)
	}

	regular, err := opentype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	face, err := opentype.NewFace(regular, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	canvas := image.NewRGBA(image.Rect(0, 0, 800, 1400))
	draw.Draw(canvas, canvas.Bounds(), image.White, image.Pt(0, 0), draw.Over)

	drawer := font.Drawer{
		Dst:  canvas,
		Src:  image.Black,
		Face: face,
		Dot:  fixed.P(0, 50),
	}

	spaceWidth := fixedToFloat(drawer.MeasureString(" "))

	for _, paragraph := range strings.Split(string(text), "\n") {
		words := strings.Split(paragraph, " ")
		boxes := make([]lb.Box, 0, len(words))
		for _, w := range words {
			boxes = append(boxes, &Word{
				text:  w,
				width: fixedToFloat(drawer.MeasureString(w)),
			})
		}

		// Do the line-breaking using Knuth-Plass algorithm.
		lines := lb.KnuthPlass(boxes, &lb.Options{
			TextWidth:     700,
			TextDirection: lb.LeftToRight,
			GlueWidth:     spaceWidth,
			GlueShrink:    spaceWidth / 5, /* 20% shrink */
			GlueExpand:    spaceWidth / 3, /* 33% expand */
		})

		for _, line := range lines {
			drawer.Dot.X = fixed.I(50)
			drawer.Dot.Y += drawer.Face.Metrics().Height.Mul(floatToFixed(1.2))
			for _, w := range line.Boxes {
				drawer.DrawString(w.(*Word).Text())
				drawer.Dot.X += floatToFixed(line.GlueWidth)
			}
		}
	}

	file, err := os.Create("sample.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err := png.Encode(file, canvas); err != nil {
		log.Fatal(err)
	}
}

var fixedOne = fixed.I(1)

func fixedToFloat(v fixed.Int26_6) float32 {
	return float32(v) / float32(fixedOne)
}

func floatToFixed(v float32) fixed.Int26_6 {
	return fixed.Int26_6(v * float32(fixedOne))
}
