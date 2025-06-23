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

// qr definition and temporary data structures
type qr struct {
	matrix           [][]module
	version          QRVersion
	error_corr_level errcorr
	size             int
	data             []byte
}

func (qr *qr) generate() {
	// Functions patterns. This sections DO NOT encode data.
	qr.finder_patterns()
	qr.separators()
	qr.timing_patterns()
	qr.alignment_patterns()

	// Encoding region
	qr.format_information()
	qr.version_information()
	qr.data_and_error_correction()

	// TODO: Add quiet zone
}

// places finder pattern modules
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

// places separator modules
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

// places timing pattern modules
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

// places alignment patter modules
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

// places an alignment pattern module (5x5) in the given position
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

// Converts a unit to a character sequence of ones and zeros
// `targetSize` controls left padding (cannot be smaller than the
// number of bits needed to represet `num`), if zero no padding is added
func uint_to_string(num uint, targetSize int) string {
	res := strings.Builder{}
	// The number of bits needed to represent 'num'
	numBits := bits.Len(uint(num))

	// Adjust for desired size
	var padding string
	if targetSize > 0 {
		if targetSize >= numBits {
			for range targetSize - numBits {
				padding = fmt.Sprintf("0%s", padding)
			}
		} else {
			panic("cannot write uint to desired size")
		}
	}

	res.WriteString(padding)

	for range numBits {
		if num&(1<<(numBits-1)) == 1<<(numBits-1) {
			res.WriteString("1")
		} else {
			res.WriteString("0")
		}
		// Shift left by 1 to check the next bit
		num <<= 1
	}

	return res.String()
}

// Converts a unit to a slice od modules. `targetSize` controls
// left padding (cannot be smaller than the number of bits needed to represet `num`)
// if zero, no padding is added
func uint_to_modules(num uint, targetSize int) []module {
	var res []module

	// The number of bits needed to represent 'num'
	numBits := bits.Len(uint(num))

	for range numBits {
		if num&(1<<(numBits-1)) == 1<<(numBits-1) {
			res = append(res, module{bit: One})
		} else {
			res = append(res, module{bit: Zero})
		}
		// Shift left by 1 to check the next bit
		num <<= 1
	}

	// Adjust for desired size
	var padding []module
	if targetSize > 0 {
		if targetSize >= numBits {
			padding = make([]module, targetSize-numBits)
			for i := range padding {
				padding[i] = module{bit: Zero}
			}
		} else {
			panic("cannot write uint to desired module length")
		}
	}
	return append(padding, res...)
}

// place format information modules
func (qr *qr) format_information() {
	// 2 bits
	err_corr_level := get_error_correction_for_level(qr.error_corr_level)

	// 3 bits
	// FIXME: Try all and select the correct mask pattern (hardcoding 101 for now)
	mask_pattern := get_mask_pattern_for_mask(3)

	// 5 bits
	format_data := err_corr_level<<3 + uint(mask_pattern.bits)

	// 10 bits
	bhc_code := encodeBCH15_5(format_data)

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

// encodes `data` to BCH(15,5). Thank you Gemini 😉
func encodeBCH15_5(data uint) uint {
	const n = 15
	const k = 5
	const numParityBits = n - k // 10
	const generatorPoly uint = 0b10100110111

	// Pad the data: data << numParityBits
	dividend := data << numParityBits

	// Initialize remainder
	remainder := dividend

	// Perform polynomial long division
	for i := range k {
		// Check the leading bit of the current remainder segment
		// The leading bit is at the (n-1 - i) position relative to the original dividend's MSB
		// For example, for 15 bits, if i=0, we check bit 14. If i=1, check bit 13, etc.
		// We need to check the bit that aligns with the MSB of the generator polynomial.

		// This is where it gets tricky with fixed uints:
		// How do you know the "current leading bit" without converting to array or complex masking?
		// You would need to check the bit at (numParityBits + k - 1 - i) position.

		currentMSBPos := numParityBits + k - 1 - i // This is the current bit position to check

		if (remainder>>currentMSBPos)&0x1 == 1 { // Check if the leading bit is 1
			// XOR with the generator polynomial, shifted to align with the current leading bit
			// The generator needs to be shifted right to align its MSB with the remainder's current MSB
			shiftAmount := currentMSBPos - 10
			remainder ^= (generatorPoly << shiftAmount)
		}
	}

	// The remainder now holds the parity bits
	// Extract the last 'numParityBits' bits (LSBs)
	return remainder & ((1 << numParityBits) - 1)
}

// encodes `data` to Golay(18,6) equivalent to BCH(18,6). Thank you Gemini 😉
func encodeGolay18_6(data uint) uint {
	const n = 18
	const k = 6
	const numParityBits = n - k // 12
	const generatorPoly uint = 0b1111100100101

	// Pad the data: data << numParityBits
	dividend := data << numParityBits

	// Initialize remainder
	remainder := dividend

	// Perform polynomial long division
	for i := range k {
		// Check the leading bit of the current remainder segment
		// The leading bit is at the (n-1 - i) position relative to the original dividend's MSB
		// For example, for 15 bits, if i=0, we check bit 14. If i=1, check bit 13, etc.
		// We need to check the bit that aligns with the MSB of the generator polynomial.

		// This is where it gets tricky with fixed uints:
		// How do you know the "current leading bit" without converting to array or complex masking?
		// You would need to check the bit at (numParityBits + k - 1 - i) position.

		currentMSBPos := numParityBits + k - 1 - i // This is the current bit position to check

		if (remainder>>currentMSBPos)&0x1 == 1 { // Check if the leading bit is 1
			// XOR with the generator polynomial, shifted to align with the current leading bit
			// The generator needs to be shifted right to align its MSB with the remainder's current MSB
			shiftAmount := currentMSBPos - 12
			remainder ^= (generatorPoly << shiftAmount)
		}
	}

	// The remainder now holds the parity bits
	// Extract the last 'numParityBits' bits (LSBs)
	return remainder & ((1 << numParityBits) - 1)
}

// places version information modules
func (qr *qr) version_information() {
	// Version information is only included for version 7 and up
	if qr.version.number < 7 {
		return
	}

	version_data := uint(qr.version.number)

	// 12 bits
	golay_code := encodeGolay18_6(version_data)

	// 18 bits
	data := version_data<<12 + golay_code

	version_modules := uint_to_modules(data, 18)

	// 3 x 6 top right module block
	// With 0 representing the least significant bit the placement must be as shown
	//  0  1  2
	//  3  4  5
	//  6  7  8
	//  9 10 11
	// 12 13 14
	// 15 16 17
	var pos int
	for i := 6 - 1; i >= 0; i-- {
		for j := qr.size - 8 - 1; j >= qr.size-8-3; j-- {
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
	for j := 6 - 1; j >= 0; j-- {
		for i := qr.size - 8 - 1; i >= qr.size-8-3; i-- {
			qr.matrix[i][j] = version_modules[pos]
			pos++
		}
	}
}

func (qr *qr) data_and_error_correction() {

}

// Append or concatenate two uints (effectively left shifts)
func uint_append(first, second uint) uint {
	return first<<bits.Len(uint(second)) + second
}

// Calculates character count of given input data in the
// corresponding data mode
func character_count(mode QRMode, version QRVersion, input string) bit_seq {
	switch mode {
	case NumericMode:
		return character_count_numeric(version, mode, input)
	default:
		panic("character_count: data mode not implemented!")
	}
}

func character_count_numeric(version QRVersion, mode QRMode, input string) bit_seq {
	bs, _ := NewBitSeqWithSize(uint(len(input)), GetCharCountLength(version, mode))
	return bs
}

func encode(mode QRMode, input string) bit_seq {
	switch mode {
	case NumericMode:
		return encode_numeric(input)
	default:
		panic("encode: data mode not implemented!")
	}
}

func encode_numeric(input string) bit_seq {
	var output bit_seq
	var prev int

	for prev < len(input) {
		if len(input) == prev+1 { // Encode in 4 bits
			group := input[prev : prev+1]
			group_uint64, err := strconv.ParseUint(group, 10, 4)
			if err != nil {
				panic(err)
			}
			b, _ := NewBitSeqWithSize(uint(group_uint64), 4)
			output = Concat(output, b)
			prev = prev + 1
		} else if len(input) == prev+2 { // Encode in 7 bits
			group := input[prev : prev+2]
			group_uint64, err := strconv.ParseUint(group, 10, 7)
			if err != nil {
				panic(err)
			}
			b, _ := NewBitSeqWithSize(uint(group_uint64), 7)
			output = Concat(output, b)
			prev = prev + 2
		} else if len(input) >= prev+3 { // Encode in 10 bits
			group := input[prev : prev+3]
			group_uint64, err := strconv.ParseUint(group, 10, 10)
			if err != nil {
				panic(err)
			}
			b, _ := NewBitSeqWithSize(uint(group_uint64), 10)
			output = Concat(output, b)
			prev = prev + 3
		}

	}
	return output
}

func NewQRCode(want_version int, error_correction_level string, input string) *qr {
	// TODO: Determine best mode from input data
	mode := NumericMode

	// Encode data
	data := encode(mode, input)

	// TODO: Determine version from encoded data length
	version := QRVersion{
		format: QR_FORMAT_QR,
		number: want_version,
	}

	character_count := character_count(mode, version, input)
	output := ConcatMany(GetModeIndicatorBits(version, mode), character_count, data)
	print_bit_seq(output)

	// NOTE: hasta acá va bien el ejemplo ======

	// Calculate error correction codewords
	size := 21 + (want_version-1)*4
	matrix := make([][]module, size)
	for i := range size {
		matrix[i] = make([]module, size)
	}
	qr := &qr{
		matrix: matrix,
		version: QRVersion{
			format: QR_FORMAT_QR,
			number: want_version,
		},
		error_corr_level: errcorr(error_correction_level),
		size:             size,
		data:             []byte(input),
	}
	qr.generate()
	return qr
}

type errcorr string

const ERR_CORR_L = "L"
const ERR_CORR_M = "M"
const ERR_CORR_Q = "Q"
const ERR_CORR_H = "H"

func (qr *qr) FullVersion() string {
	return fmt.Sprintf("%s-%s", qr.version.String(), qr.error_corr_level)
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
	b.WriteString(fmt.Sprintf("QR version %s (size %d)", qr.FullVersion(), qr.size))

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

	// Data
	data := "01234567"
	if len(args) > 2 {
		data = args[2]
	}

	qr := NewQRCode(version, err_corr_level, data)
	fmt.Println(qr.String())
	qr.Draw()
	// qr.Print()
}
