package main

import (
	"bytes"
	"fmt"
)

const Y = `■` // \u25A0
const N = `□` // \u25A1
const F = `▿` // \u25BF or \u25AA or \u25AB
var W = module{value: N}
var B = module{value: Y}

type qr struct {
	matrix  [][]module
	version qrversion
	format  qrformat
	size    int
}

func (qr *qr) init() {
	// Filler
	for i := range qr.size {
		for j := range qr.size {
			qr.matrix[i][j] = module{value: F}
		}
	}

	// finder patterns
	// upper right corner
	for i := range 7 {
		for j := range 7 {
			qr.matrix[i][j] = module{value: Y}
		}
	}
	for i := 1; i < 6; i++ {
		for j := 1; j < 6; j++ {
			qr.matrix[i][j] = module{value: N}
		}
	}
	for i := 2; i < 5; i++ {
		for j := 2; j < 5; j++ {
			qr.matrix[i][j] = module{value: Y}
		}
	}
}

func NewQRCode(version int, is_micro bool) *qr {
	size := 21 + (version-1)*4
	matrix := make([][]module, size)
	for i := range size {
		matrix[i] = make([]module, size)
	}
	qr := &qr{
		matrix: matrix,
		version: qrversion{
			number: version,
		},
		format: QR_FORMAT_QR,
		size:   size,
	}
	if is_micro {
		qr.format = QR_FORMAT_MICRO_QR
	}
	qr.init()
	return qr
}

type qrformat string

const QR_FORMAT_QR = "full"
const QR_FORMAT_MICRO_QR = "micro"

type qrversion struct {
	error_corr_level string
	number           int
}

func (qr *qr) Version() string {
	var format string
	if qr.format == QR_FORMAT_MICRO_QR {
		format = "M"
	}
	return fmt.Sprintf("%s%d-%s", format, qr.version.number, qr.version.error_corr_level)
}

type module struct {
	value string
}

func (qr *qr) String() string {
	var b bytes.Buffer
	// Write header line
	b.WriteString("  ")
	for i := range qr.size {
		b.WriteString(fmt.Sprintf("%2d ", i))
	}
	b.WriteString("\n")

	// TODO: quiet zone

	// Write matrix with line number prefix
	for i, line := range qr.matrix {
		b.WriteString(fmt.Sprintf("%2d ", i))
		for _, cell := range line {
			b.WriteString(cell.value + "  ")
		}
		b.WriteString("\n")
	}
	return b.String()
}

func main() {
	qr := NewQRCode(1, false)
	fmt.Println(qr.String())
}
