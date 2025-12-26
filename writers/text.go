package writers

import (
	"bytes"
	"image/color"
	"os"
)

const txtpath = "qr.txt"

const (
	Black = `■` // \u25A0
	White = `□` // \u25A1
	undef = `▿` // \u25BF or \u25AA or \u25AB
)

type TextRequest struct {
	Size  int
	Chars [][]color.Color
}

func WriteText(req TextRequest) {
	f, _ := os.Create(txtpath)
	defer f.Close()

	var b bytes.Buffer
	for _, line := range req.Chars {
		for _, cell := range line {
			b.WriteString(Char(cell))
		}
		b.WriteString("\n")
	}
	f.WriteString(b.String())
}

func Char(c color.Color) string {
	switch c {
	case color.Black:
		return Black
	case color.White:
		return White
	default:
		return undef
	}
}
