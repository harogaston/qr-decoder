package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

const Y = `■`     // \u25A0
const N = `□`     // \u25A1
const undef = `▿` // \u25BF or \u25AA or \u25AB

type qr struct {
	matrix  [][]module
	version qrversion
	format  qrformat
	size    int
}

func (qr *qr) init() {
	qr.dummy_filler()
	qr.finder_patterns()
	qr.separators()
	qr.timing_patterns()
	qr.alignment_patterns()
	qr.format_information()
	qr.version_information()
}

func (qr *qr) dummy_filler() {
	for i := range qr.size {
		for j := range qr.size {
			qr.matrix[i][j] = module{value: undef}
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

func (qr *qr) separators() {
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

func (qr *qr) timing_patterns() {
	// TODO: fix for MQ
	// row 6
	alternating_flag := false
	for j := 8; j < qr.size-8; j++ {
		if alternating_flag {
			qr.matrix[6][j] = module{value: N}
		} else {
			qr.matrix[6][j] = module{value: Y}
		}
		alternating_flag = !alternating_flag
	}

	alternating_flag = false
	// column 6
	for i := 8; i < qr.size-8; i++ {
		if alternating_flag {
			qr.matrix[i][6] = module{value: N}
		} else {
			qr.matrix[i][6] = module{value: Y}
		}
		alternating_flag = !alternating_flag
	}
}

func (qr *qr) alignment_patterns() {
	coords := get_alignment_patterns_for_version(qr.version.number)
	for _, alignment_pos := range coords {
		alignment_pattern_upper_left := []int{alignment_pos[0] - 2, alignment_pos[1] - 2}
		alignment_pattern_lower_left := []int{alignment_pos[0] + 2, alignment_pos[1] - 2}
		alignment_pattern_upper_right := []int{alignment_pos[0] - 2, alignment_pos[1] + 2}

		// check upper left colission
		finder_lower_right := []int{6, 6}
		if finder_lower_right[0] >= alignment_pattern_upper_left[0] &&
			finder_lower_right[1] >= alignment_pattern_upper_left[1] {
			continue
		}

		// check lower left colission
		finder_upper_right := []int{qr.size - 7, 6}
		if finder_upper_right[0] <= alignment_pattern_lower_left[0] &&
			finder_upper_right[1] >= alignment_pattern_lower_left[1] {
			continue
		}

		// check upper right colission
		finder_lower_left := []int{6, qr.size - 7}
		if finder_lower_left[0] >= alignment_pattern_upper_right[0] &&
			finder_lower_left[1] <= alignment_pattern_upper_right[1] {
			continue
		}

		// no colissions good to go
		qr.add_alignment_pattern_modules(alignment_pos[0], alignment_pos[1])
	}
}

func (qr *qr) add_alignment_pattern_modules(row int, col int) {
	// 5 by 5 dark square
	for i := row - 2; i <= row+2; i++ {
		for j := col - 2; j <= col+2; j++ {
			qr.matrix[i][j] = module{value: Y}
		}
	}
	// 3 by 3 light square
	for i := row - 1; i <= row+1; i++ {
		for j := col - 1; j <= col+1; j++ {
			qr.matrix[i][j] = module{value: N}
		}
	}

	// single central dark module
	qr.matrix[row][col] = module{value: Y}
}

func (qr *qr) format_information() {
	// FIXME: Fill dummy data for now
	// row 8
	for j := range qr.size {
		if j < 6 || j > 6 && j < 8 || j > qr.size-8-1 {
			qr.matrix[8][j] = module{value: "f"}
		}
	}
	// column 8
	for i := range qr.size {
		if i < 6 || i > 6 && i <= 8 || i > qr.size-8 {
			qr.matrix[i][8] = module{value: "f"}
		}
	}

	// set always dark module 4V + 9, 8
	qr.matrix[4*qr.version.number+9][8] = module{value: Y}
}

func (qr *qr) version_information() {
	// Version information is only included for version 7 and up
	if qr.version.number < 7 {
		return
	}

	// FIXME: Fill dummy data for now
	// 6 x 3 top right module block
	for i := range 6 {
		for j := qr.size - 8 - 1; j > qr.size-8-1-3; j-- {
			qr.matrix[i][j] = module{value: "v"}
		}
	}

	// 3 x 6 lower left module block
	for i := qr.size - 8 - 1; i > qr.size-8-1-3; i-- {
		for j := range 6 {
			qr.matrix[i][j] = module{value: "v"}
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
	// FIXME: support more than version 20 or size 99 modules
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
	args := os.Args[1:]
	version := 1
	if len(args) > 0 {
		if v, err := strconv.Atoi(args[0]); err == nil {
			version = v
		}
	}
	fmt.Println(version)
	qr := NewQRCode(version, false)
	fmt.Println(qr.String())
}
