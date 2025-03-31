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
	qr.dummy_filler()
	qr.finder_patterns()
	qr.separator()
}

func (qr *qr) dummy_filler() {
	for i := range qr.size {
		for j := range qr.size {
			qr.matrix[i][j] = module{value: F}
		}
	}
}

func (qr *qr) finder_patterns() {
	// upper left corner
	for i := range 7 { // size 7
		for j := range 7 {
			qr.matrix[i][j] = module{value: Y}
		}
	}
	for i := 1; i < 6; i++ { // size 5
		for j := 1; j < 6; j++ {
			qr.matrix[i][j] = module{value: N}
		}
	}
	for i := 2; i < 5; i++ { // size 3
		for j := 2; j < 5; j++ {
			qr.matrix[i][j] = module{value: Y}
		}
	}

	// lower left corner
	for i := qr.size - 1; i > qr.size-7-1; i-- { // size 7
		for j := range 7 {
			qr.matrix[i][j] = module{value: Y}
		}
	}
	for i := qr.size - 1 - 1; i > qr.size-6-1; i-- { // size 5
		for j := 1; j < 6; j++ {
			qr.matrix[i][j] = module{value: N}
		}
	}
	for i := qr.size - 1 - 2; i > qr.size-5-1; i-- { // size 3
		for j := 2; j < 5; j++ {
			qr.matrix[i][j] = module{value: Y}
		}
	}

	// upper rigth corner
	for i := range 7 { // size 7
		for j := qr.size - 1; j > qr.size-7-1; j-- {
			qr.matrix[i][j] = module{value: Y}
		}
	}
	for i := 1; i < 6; i++ { // size 5
		for j := qr.size - 1 - 1; j > qr.size-6-1; j-- {
			qr.matrix[i][j] = module{value: N}
		}
	}
	for i := 2; i < 5; i++ { // size 3
		for j := qr.size - 1 - 2; j > qr.size-5-1; j-- {
			qr.matrix[i][j] = module{value: Y}
		}
	}
}

func (qr *qr) separator() {
	// upper left
	for i := range 8 {
		qr.matrix[i][7] = module{value: N}
	}
	for j := range 8 {
		qr.matrix[7][j] = module{value: N}
	}

	// lower left
	for i := qr.size - 1; i > qr.size-7-1; i-- {
		qr.matrix[i][7] = module{value: N}
	}
	for j := range 8 {
		qr.matrix[qr.size-7-1][j] = module{value: N}
	}

	// upper right
	for i := range 8 {
		qr.matrix[i][qr.size-7-1] = module{value: N}
	}
	for j := qr.size - 1; j > qr.size-7-1; j-- {
		qr.matrix[7][j] = module{value: N}
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
