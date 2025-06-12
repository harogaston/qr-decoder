package main

import (
	"fmt"
	"image/color"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

const Black = `■` // \u25A0
const White = `□` // \u25A1
const undef = `▿` // \u25BF or \u25AA or \u25AB

type qr struct {
	matrix  [][]module
	version qrversion
	format  qrformat
	size    int
}

func (qr *qr) init() {
	// Functions patterns
	qr.finder_patterns()
	qr.separators()
	qr.timing_patterns()
	qr.alignment_patterns()

	// Encoding region
	qr.format_information()
	qr.version_information()
	qr.data()
	// TODO: Add quiet zone
}

func (qr *qr) finder_patterns() {
	// upper left corner
	for i := range 7 { // size 7
		for j := range 7 {
			qr.matrix[i][j] = module{bit: One}
		}
	}
	for i := 1; i < 6; i++ { // size 5
		for j := 1; j < 6; j++ {
			qr.matrix[i][j] = module{bit: Zero}
		}
	}
	for i := 2; i < 5; i++ { // size 3
		for j := 2; j < 5; j++ {
			qr.matrix[i][j] = module{bit: One}
		}
	}

	// lower left corner
	for i := qr.size - 1; i > qr.size-7-1; i-- { // size 7
		for j := range 7 {
			qr.matrix[i][j] = module{bit: One}
		}
	}
	for i := qr.size - 1 - 1; i > qr.size-6-1; i-- { // size 5
		for j := 1; j < 6; j++ {
			qr.matrix[i][j] = module{bit: Zero}
		}
	}
	for i := qr.size - 1 - 2; i > qr.size-5-1; i-- { // size 3
		for j := 2; j < 5; j++ {
			qr.matrix[i][j] = module{bit: One}
		}
	}

	// upper rigth corner
	for i := range 7 { // size 7
		for j := qr.size - 1; j > qr.size-7-1; j-- {
			qr.matrix[i][j] = module{bit: One}
		}
	}
	for i := 1; i < 6; i++ { // size 5
		for j := qr.size - 1 - 1; j > qr.size-6-1; j-- {
			qr.matrix[i][j] = module{bit: Zero}
		}
	}
	for i := 2; i < 5; i++ { // size 3
		for j := qr.size - 1 - 2; j > qr.size-5-1; j-- {
			qr.matrix[i][j] = module{bit: One}
		}
	}
}

func (qr *qr) separators() {
	// upper left
	for i := range 8 {
		qr.matrix[i][7] = module{bit: Zero}
	}
	for j := range 8 {
		qr.matrix[7][j] = module{bit: Zero}
	}

	// lower left
	for i := qr.size - 1; i > qr.size-7-1; i-- {
		qr.matrix[i][7] = module{bit: Zero}
	}
	for j := range 8 {
		qr.matrix[qr.size-7-1][j] = module{bit: Zero}
	}

	// upper right
	for i := range 8 {
		qr.matrix[i][qr.size-7-1] = module{bit: Zero}
	}
	for j := qr.size - 1; j > qr.size-7-1; j-- {
		qr.matrix[7][j] = module{bit: Zero}
	}
}

func (qr *qr) timing_patterns() {
	// row 6
	alternating_flag := false
	for j := 8; j < qr.size-8; j++ {
		if alternating_flag {
			qr.matrix[6][j] = module{bit: Zero}
		} else {
			qr.matrix[6][j] = module{bit: One}
		}
		alternating_flag = !alternating_flag
	}

	// column 6
	alternating_flag = false
	for i := 8; i < qr.size-8; i++ {
		if alternating_flag {
			qr.matrix[i][6] = module{bit: Zero}
		} else {
			qr.matrix[i][6] = module{bit: One}
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
		qr.add_alignment_pattern_module(alignment_pos[0], alignment_pos[1])
	}
}

func (qr *qr) add_alignment_pattern_module(row int, col int) {
	// 5 by 5 dark square
	for i := row - 2; i <= row+2; i++ {
		for j := col - 2; j <= col+2; j++ {
			qr.matrix[i][j] = module{bit: One}
		}
	}
	// 3 by 3 light square
	for i := row - 1; i <= row+1; i++ {
		for j := col - 1; j <= col+1; j++ {
			qr.matrix[i][j] = module{bit: Zero}
		}
	}

	// single central dark module
	qr.matrix[row][col] = module{bit: One}
}

func uint_to_modules(num uint, targetSize int) []module {
	var res []module

	// The number of bits needed to represent 'num'
	numBits := bits.Len(uint(num))

	for range numBits {
		if num&1 == 1 {
			res = append(res, module{bit: One})
		} else {
			res = append(res, module{bit: Zero})
		}
		// Shift right by 1 to check the next bit
		num >>= 1
	}

	// Adjust for desired size
	var padding []module
	if targetSize > 0 && targetSize > numBits {
		padding = make([]module, targetSize-numBits)
		for i := range padding {
			padding[i] = module{bit: Zero}
		}
	}
	return append(padding, res...)
}

func print_module_slice(mods []module) {
	b := strings.Builder{}
	for _, m := range mods {
		b.WriteString(m.Char())
	}
	fmt.Println(b.String())
}

func (qr *qr) format_information() {
	// 2 bits
	err_corr_level := get_error_correction_for_level(qr.version.error_corr_level)

	// 3 bits
	// FIXME: Try all and select the correct mask pattern (hardcoding 0 for now)
	mask_pattern := get_mask_pattern_for_mask(0)

	// 5 bits
	format_data := err_corr_level<<3 + uint(mask_pattern.bits)

	// 10 bits
	// TODO: Implement BHC algorithm
	bhc_code := get_bhc_codes(format_data)

	// 15 bits
	unmasked_data := format_data<<10 + bhc_code
	xor_mask_pattern := uint(0b101010000010010)

	data := unmasked_data ^ xor_mask_pattern
	format_modules := uint_to_modules(data, 15)

	var pos int
	// row 8
	for j := range qr.size {
		if j < 6 || j == 7 || j > qr.size-8-1 {
			qr.matrix[8][j] = format_modules[pos]
			pos++
		}
	}

	pos = 0
	// column 8
	for i := qr.size - 1; i >= 0; i-- {
		if i < 6 || i > 6 && i < 9 || i > qr.size-8 {
			qr.matrix[i][8] = format_modules[pos]
			pos++
		}
	}

	// set always dark module 4V + 9, 8
	qr.matrix[4*qr.version.number+9][8] = module{bit: One}
}

func get_bhc_codes(_ uint) uint {
	return 0
}

func (qr *qr) version_information() {
	// Version information is only included for version 7 and up
	if qr.version.number < 7 {
		return
	}

	version_pattern := uint_to_modules(uint(qr.version.number), 18)
	// FIXME: implement error correction
	version_modules := version_pattern

	// 3 x 6 top right module block
	// With 0 representing the least significant bit the placement must be as shown
	//  0  1  2
	//  3  4  5
	//  6  7  8
	//  9 10 11
	// 12 13 14
	// 15 16 17
	var pos int
	for i := range 6 {
		for j := qr.size - 8 - 3; j < qr.size-8; j++ {
			qr.matrix[i][j] = version_modules[pos]
			pos++
		}
	}

	// 6 x 3 lower left module block
	// With 0 representing the least significant bit the placement must be as shown
	// 0  3  6  9 12 15
	// 1  4  7 10 13 16
	// 2  5  8 11 14 17
	pos = 0
	for j := range 6 {
		for i := qr.size - 8 - 3; i < qr.size-8; i++ {
			qr.matrix[i][j] = version_modules[pos]
			pos++
		}
	}
}

func (qr *qr) data() {

}

func NewQRCode(version int, is_micro bool, error_correction_level string) *qr {
	size := 21 + (version-1)*4
	matrix := make([][]module, size)
	for i := range size {
		matrix[i] = make([]module, size)
	}
	qr := &qr{
		matrix: matrix,
		version: qrversion{
			number:           version,
			error_corr_level: errcorr(error_correction_level),
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

type errcorr string

const ERR_CORR_L = "L"
const ERR_CORR_M = "M"
const ERR_CORR_Q = "Q"
const ERR_CORR_H = "H"

type qrversion struct {
	error_corr_level errcorr
	number           int
}

func (qr *qr) Version() string {
	var format string
	if qr.format == QR_FORMAT_MICRO_QR {
		format = "M"
	}
	return fmt.Sprintf("%s%d-%s", format, qr.version.number, qr.version.error_corr_level)
}

type Bit uint

const (
	Undef Bit = iota
	Zero
	One
)

type module struct {
	bit Bit
}

func (m *module) Color() color.Color {
	if m.bit == Zero {
		return color.White
	} else if m.bit == One {
		return color.Black
	} else {
		return color.Transparent
	}
}

func (m *module) Char() string {
	if m.bit == Zero {
		return White
	} else if m.bit == One {
		return Black
	} else {
		return undef
	}
}

func (qr *qr) String() string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("QR version %s (size %d)", qr.Version(), qr.size))

	return b.String()
}

func (qr *qr) Print() {
	req := TextRequest{
		Size:  qr.size,
		Chars: qr.matrix,
	}
	WriteText(req)
}

func (qr *qr) Draw() {
	pixs := make([][]color.Color, len(qr.matrix[0]))
	for y, row := range qr.matrix {
		imgRow := make([]color.Color, len(row))
		for x, c := range row {
			imgRow[x] = c.Color()
		}
		pixs[y] = imgRow
	}
	req := ImageRequest{
		Scale:  4,
		Pixels: pixs,
	}
	WriteImage(req)
}

func main() {
	args := os.Args[1:]

	// QR Code version parsing
	version := 1
	if len(args) > 0 {
		if v, err := strconv.Atoi(args[0]); err == nil {
			version = v
		}
	}

	// Error correction level parsing
	err_corr_level := "L"
	if len(args) > 1 {
		if args[1] == ERR_CORR_L || args[1] == ERR_CORR_M || args[1] == ERR_CORR_Q || args[1] == ERR_CORR_H {
			err_corr_level = args[1]
		}
	}

	qr := NewQRCode(version, false, err_corr_level)
	fmt.Println(qr.String())
	qr.Draw()
	// qr.Print()
}
