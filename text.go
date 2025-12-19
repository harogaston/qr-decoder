package main

import (
	"bytes"
	"os"
)

const txtpath = "qr.txt"

type TextRequest struct {
	Size  int
	Chars [][]module
}

func WriteText(req TextRequest) {
	f, _ := os.Create(txtpath)
	defer f.Close()

	var b bytes.Buffer
	for _, line := range req.Chars {
		for _, cell := range line {
			b.WriteString(cell.Char())
		}
		b.WriteString("\n")
	}
	f.WriteString(b.String())
}
