package main

import (
	"bytes"
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
	Chars [][]module
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

func Char(m module) string {
	switch m.bit {
	case Zero:
		return White
	case One:
		return Black
	default:
		return undef
	}
}
