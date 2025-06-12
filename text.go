package main

import (
	"bytes"
	"fmt"
)

type TextRequest struct {
	Size  int
	Chars [][]module
}

func WriteText(req TextRequest) {
	// FIXME: support more than version 20 or size 99 modules
	var b bytes.Buffer
	// Write header line
	b.WriteString("  ")
	for i := range req.Size {
		b.WriteString(fmt.Sprintf("%2d ", i))
	}
	b.WriteString("\n")

	// TODO: quiet zone

	// Write matrix with line number prefix
	for i, line := range req.Chars {
		b.WriteString(fmt.Sprintf("%2d ", i))
		for _, cell := range line {
			b.WriteString(cell.Char() + "  ")
		}
		b.WriteString("\n")
	}
	fmt.Println(b.String())
}
