package main

import (
	"fmt"
	"image/color"
	"math"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

const (
	Black = `■` // \u25A0
	White = `□` // \u25A1
	undef = `▿` // \u25BF or \u25AA or \u25AB
)

// qr definition and temporary data structures
type qr struct {
	matrix              [][]module
	version             QRVersion
	error_corr_level    errcorr
	size                int
	data                []byte
	encoded_data        bit_seq
	mode                QRMode
	mask                int
	is_function_pattern [][]bool
	debug_no_mask       bool
}

func (qr *qr) generate() {
	// Functions patterns. This sections DO NOT encode data.
	qr.finder_patterns()
	qr.separators()
	qr.timing_patterns()
	qr.alignment_patterns()

	// Encoding region
	// Encoding region
	// qr.format_information() // Removed: handled in masking loop
	qr.version_information()
	qr.reserve_format_information_area()
	qr.data_and_error_correction()

	// Masking
	// We need to try all 8 masks and pick the best one.
	// But format information depends on the mask!
	// So we have to:
	// 1. Generate matrix with data (already done by data_and_error_correction)
	// 2. For each mask pattern:
	//    a. Apply mask to a copy of matrix
	//    b. Calculate penalty
	// 3. Select best mask
	// 4. Apply best mask to original matrix (or keep the best copy)
	// 5. Re-generate format information with the selected mask

	minPenalty := math.MaxInt32
	var bestMatrix [][]module
	var bestMatrixMask int

	// Save original matrix state (before masking)
	// Actually, data_and_error_correction places data.
	// Masking flips bits.
	// So we can just clone the matrix for each try.

	originalMatrix := make([][]module, qr.size)
	for i := range qr.size {
		originalMatrix[i] = make([]module, qr.size)
		copy(originalMatrix[i], qr.matrix[i])
	}

	if qr.debug_no_mask {
		fmt.Println("DEBUG: Masking disabled (--debug-no-mask)")
		// Use original matrix, but we still need to place format info.
		// We'll use mask 0 for format info generation, but won't apply the mask XOR to data.
		qr.matrix = originalMatrix
		qr.mask = 0
		qr.place_format_information(0)
		bestMatrix = originalMatrix
		bestMatrixMask = 0
	} else {
		for mask := 0; mask < 8; mask++ {
			// Apply mask
			masked := qr.apply_mask(mask, originalMatrix)

			// Add format info (it also has mask bits!)
			// But format info is in the reserved area.
			// The spec says:
			// "The format information is added to the symbol after the masking process."
			// Wait, format info contains the mask pattern.
			// So we need to calculate penalty including format info?
			// "The penalty score is calculated based on the entire symbol."
			// So we should place format info (with current mask) before calculating penalty.

			qr.matrix = masked // Temporarily set to masked to call format_information
			qr.place_format_information(mask)

			penalty := calculatePenalty(qr.matrix)
			if penalty < minPenalty {
				minPenalty = penalty
				bestMatrix = masked // This already has format info for this mask
				bestMatrixMask = mask
			}
		}
	}

	// Set final matrix
	qr.matrix = bestMatrix
	qr.mask = bestMatrixMask // We need to capture the mask index
	// Format info is already there.

	qr.add_quiet_zone()
}

func (qr *qr) add_quiet_zone() {
	// Quiet zone is 4 modules wide on all sides for standard QR codes
	quietZoneWidth := 4
	if qr.version.format == "Micro" {
		// Quiet zone is 2 modules wide on all sides for Micro QR codes
		quietZoneWidth = 2
	}
	newSize := qr.size + 2*quietZoneWidth
	newMatrix := make([][]module, newSize)

	for i := 0; i < newSize; i++ {
		newMatrix[i] = make([]module, newSize)
		for j := 0; j < newSize; j++ {
			// Default is white (Zero)
			newMatrix[i][j] = module{bit: Zero}
		}
	}

	// Copy original matrix into the center
	for i := 0; i < qr.size; i++ {
		for j := 0; j < qr.size; j++ {
			newMatrix[i+quietZoneWidth][j+quietZoneWidth] = qr.matrix[i][j]
		}
	}

	qr.matrix = newMatrix
}

// Helper to place format info with specific mask
func (qr *qr) place_format_information(mask int) {
	// 2 bits

	// 5 bits
	// Calculate format information
	// 15 bits: 5 data bits + 10 BCH bits, masked with 101010000010010
	formatInfo, err := GenerateFormatInformation(qr.error_corr_level, mask)
	if err != nil {
		panic(err)
	}

	// Convert to modules (bit 14 is MSB, bit 0 is LSB)
	format_modules := make([]module, 15)
	for i := 0; i < 15; i++ {
		if (formatInfo>>(14-i))&1 == 1 {
			format_modules[i] = module{bit: One}
		} else {
			format_modules[i] = module{bit: Zero}
		}
	}

	// Copy 1: Top-Left (around finder pattern)
	// Bits 14-9 at (8, 0-5)
	for i := 0; i < 6; i++ {
		qr.matrix[8][i] = format_modules[i]
	}
	// Bit 8 at (8, 7)
	qr.matrix[8][7] = format_modules[6]
	// Bit 7 at (8, 8)
	qr.matrix[8][8] = format_modules[7]
	// Bit 6 at (7, 8)
	qr.matrix[7][8] = format_modules[8]
	// Bits 5-0 at (5-0, 8)
	for i := 0; i < 6; i++ {
		qr.matrix[5-i][8] = format_modules[9+i]
	}

	// Copy 2: Split (Top-Right and Bottom-Left)
	// Top-Right: Bits 7-0 at (8, size-8 to size-1)
	// Bit 7 at (8, size-8) ... Bit 0 at (8, size-1)
	for i := 0; i < 8; i++ {
		qr.matrix[8][qr.size-8+i] = format_modules[7+i] // format_modules[7] is Bit 7, [14] is Bit 0
	}
	// Bottom-Left: Bits 14-8 at (size-7 to size-1, 8)
	// Bit 14 at (size-7, 8) ... Bit 8 at (size-1, 8)
	for i := 0; i < 7; i++ {
		qr.matrix[qr.size-7+i][8] = format_modules[i] // format_modules[0] is Bit 14, [6] is Bit 8
	}

	// set always dark module 4V + 9, 8
	qr.matrix[4*qr.version.number+9][8] = module{bit: One}
}

// marks the format information area as reserved
func (qr *qr) reserve_format_information_area() {
	// row 8
	for j := range qr.size {
		if j < 6 || j == 7 || j > qr.size-8-1 {
			qr.is_function_pattern[8][j] = true
		}
	}

	// column 8
	for i := qr.size - 1; i >= 0; i-- {
		if i < 6 || i > 6 && i < 9 || i > qr.size-8 {
			qr.is_function_pattern[i][8] = true
		}
	}

	// set always dark module 4V + 9, 8
	qr.is_function_pattern[4*qr.version.number+9][8] = true
}

// places finder pattern modules
func (qr *qr) finder_patterns() {
	// upper left corner
	for i := range 7 { // size 7
		for j := range 7 {
			qr.matrix[i][j] = module{bit: One}
			qr.is_function_pattern[i][j] = true
		}
	}
	for i := 1; i < 6; i++ { // size 5
		for j := 1; j < 6; j++ {
			qr.matrix[i][j] = module{bit: Zero}
			qr.is_function_pattern[i][j] = true
		}
	}
	for i := 2; i < 5; i++ { // size 3
		for j := 2; j < 5; j++ {
			qr.matrix[i][j] = module{bit: One}
			qr.is_function_pattern[i][j] = true
		}
	}

	// lower left corner
	for i := qr.size - 1; i > qr.size-7-1; i-- { // size 7
		for j := range 7 {
			qr.matrix[i][j] = module{bit: One}
			qr.is_function_pattern[i][j] = true
		}
	}
	for i := qr.size - 1 - 1; i > qr.size-6-1; i-- { // size 5
		for j := 1; j < 6; j++ {
			qr.matrix[i][j] = module{bit: Zero}
			qr.is_function_pattern[i][j] = true
		}
	}
	for i := qr.size - 1 - 2; i > qr.size-5-1; i-- { // size 3
		for j := 2; j < 5; j++ {
			qr.matrix[i][j] = module{bit: One}
			qr.is_function_pattern[i][j] = true
		}
	}

	// upper rigth corner
	for i := range 7 { // size 7
		for j := qr.size - 1; j > qr.size-7-1; j-- {
			qr.matrix[i][j] = module{bit: One}
			qr.is_function_pattern[i][j] = true
		}
	}
	for i := 1; i < 6; i++ { // size 5
		for j := qr.size - 1 - 1; j > qr.size-6-1; j-- {
			qr.matrix[i][j] = module{bit: Zero}
			qr.is_function_pattern[i][j] = true
		}
	}
	for i := 2; i < 5; i++ { // size 3
		for j := qr.size - 1 - 2; j > qr.size-5-1; j-- {
			qr.matrix[i][j] = module{bit: One}
			qr.is_function_pattern[i][j] = true
		}
	}
}

// places separator modules
func (qr *qr) separators() {
	// upper left
	for i := range 8 {
		qr.matrix[i][7] = module{bit: Zero}
		qr.is_function_pattern[i][7] = true
	}
	for j := range 8 {
		qr.matrix[7][j] = module{bit: Zero}
		qr.is_function_pattern[7][j] = true
	}

	// lower left
	for i := qr.size - 1; i > qr.size-7-1; i-- {
		qr.matrix[i][7] = module{bit: Zero}
		qr.is_function_pattern[i][7] = true
	}
	for j := range 8 {
		qr.matrix[qr.size-7-1][j] = module{bit: Zero}
		qr.is_function_pattern[qr.size-7-1][j] = true
	}

	// upper right
	for i := range 8 {
		qr.matrix[i][qr.size-7-1] = module{bit: Zero}
		qr.is_function_pattern[i][qr.size-7-1] = true
	}
	for j := qr.size - 1; j > qr.size-7-1; j-- {
		qr.matrix[7][j] = module{bit: Zero}
		qr.is_function_pattern[7][j] = true
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
		qr.is_function_pattern[6][j] = true
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
		qr.is_function_pattern[i][6] = true
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
			qr.is_function_pattern[i][j] = true
		}
	}
	// 3 by 3 light square
	for i := row - 1; i <= row+1; i++ {
		for j := col - 1; j <= col+1; j++ {
			qr.matrix[i][j] = module{bit: Zero}
			qr.is_function_pattern[i][j] = true
		}
	}

	// single central dark module
	qr.matrix[row][col] = module{bit: One}
	qr.is_function_pattern[row][col] = true
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

// encodes `data` to BCH(15,5). Thank you Gemini 😉

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
			qr.is_function_pattern[i][j] = true
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
			qr.is_function_pattern[i][j] = true
			pos++
		}
	}
}

func (qr *qr) data_and_error_correction() {
	// 1. Get data codewords and block info
	data := capacityData[qr.version.number]
	ecInfo := data.ecInfo[qr.error_corr_level]

	// Convert bit_seq to bytes
	dataBytes := qr.encoded_data.ToBytes()

	// 2. Split into blocks and calculate EC
	var dataBlocks [][]byte
	var ecBlocks [][]byte

	offset := 0
	for _, group := range ecInfo.BlockGroups {
		for i := 0; i < group.NumBlocks; i++ {
			// Extract data block
			end := min(offset+group.DataCodewords, len(dataBytes))
			block := dataBytes[offset:end]
			// Ensure block is copied to avoid slicing issues if we modify it (we don't)
			// But we need to make sure it has the correct length.
			// If it's shorter, we should probably pad it?
			// The spec says data codewords are filled.
			// If we are short, it means we didn't generate enough padding bits.
			// But we should have.

			dataBlocks = append(dataBlocks, block)
			offset = end

			// Calculate EC block
			ecBlock := reedSolomonEncode(block, ecInfo.ECCodewordsPerBlock)

			fmt.Printf("Data Codewords (%d): %v\n", len(block), block)
			fmt.Printf("EC Codewords (%d): %v\n", len(ecBlock), ecBlock)

			ecBlocks = append(ecBlocks, ecBlock)
		}
	}

	// 3. Interleave Data
	var finalMessage []byte

	// Max data length
	maxDataLen := 0
	for _, b := range dataBlocks {
		if len(b) > maxDataLen {
			maxDataLen = len(b)
		}
	}

	for i := 0; i < maxDataLen; i++ {
		for _, block := range dataBlocks {
			if i < len(block) {
				finalMessage = append(finalMessage, block[i])
			}
		}
	}

	// 4. Interleave EC
	// All EC blocks have same length
	ecLen := ecInfo.ECCodewordsPerBlock
	for i := 0; i < ecLen; i++ {
		for _, block := range ecBlocks {
			finalMessage = append(finalMessage, block[i])
		}
	}

	// 5. Place modules
	qr.placeCodewords(finalMessage)
}

func (qr *qr) placeCodewords(data []byte) {
	// Zig-zag scan
	// Start at bottom right
	row := qr.size - 1
	col := qr.size - 1
	direction := -1 // -1 for up, 1 for down

	bitIndex := 0
	byteIndex := 0

	for col > 0 {
		if col == 6 { // Skip timing pattern column
			col--
		}

		for row >= 0 && row < qr.size {
			for c := 0; c < 2; c++ {
				x := col - c
				y := row

				// Skip function patterns
				if !qr.isFunctionPattern(y, x) {
					// Place bit
					var bit int
					if byteIndex < len(data) {
						if (data[byteIndex]>>(7-bitIndex))&1 == 1 {
							bit = 1
						} else {
							bit = 0
						}
						bitIndex++
						if bitIndex == 8 {
							bitIndex = 0
							byteIndex++
						}
					} else {
						// Remainder bits (should be 0)
						bit = 0
					}

					if bit == 1 {
						qr.matrix[y][x] = module{bit: One}
					} else {
						qr.matrix[y][x] = module{bit: Zero}
					}
				}
			}
			row += direction
		}
		row -= direction       // Step back to valid row
		direction = -direction // Change direction
		col -= 2
	}
}

// Append or concatenate two uints (effectively left shifts)
func uint_append(first, second uint) uint {
	return first<<bits.Len(uint(second)) + second
}

// Calculates character count of given input data in the
// corresponding data mode
func character_count(mode QRMode, version QRVersion, input string) bit_seq {
	bs, _ := NewBitSeqWithSize(uint(len(input)), GetCharCountLength(version, mode))
	return bs
}

var alphanumericValues = map[rune]int{
	'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
	'A': 10, 'B': 11, 'C': 12, 'D': 13, 'E': 14, 'F': 15, 'G': 16, 'H': 17, 'I': 18, 'J': 19,
	'K': 20, 'L': 21, 'M': 22, 'N': 23, 'O': 24, 'P': 25, 'Q': 26, 'R': 27, 'S': 28, 'T': 29,
	'U': 30, 'V': 31, 'W': 32, 'X': 33, 'Y': 34, 'Z': 35,
	' ': 36, '$': 37, '%': 38, '*': 39, '+': 40, '-': 41, '.': 42, '/': 43, ':': 44,
}

func encode(mode QRMode, input string) bit_seq {
	switch mode {
	case NumericMode:
		return encode_numeric(input)
	case AlphanumericMode:
		return encode_alphanumeric(input)
	case ByteMode:
		return encode_byte(input)
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

func encode_alphanumeric(input string) bit_seq {
	var output bit_seq
	for i := 0; i < len(input); i += 2 {
		if i+1 < len(input) {
			// Pair
			val1 := alphanumericValues[rune(input[i])]
			val2 := alphanumericValues[rune(input[i+1])]
			val := 45*val1 + val2
			b, _ := NewBitSeqWithSize(uint(val), 11)
			output = Concat(output, b)
		} else {
			// Single
			val := alphanumericValues[rune(input[i])]
			b, _ := NewBitSeqWithSize(uint(val), 6)
			output = Concat(output, b)
		}
	}
	return output
}

func encode_byte(input string) bit_seq {
	var output bit_seq
	for i := 0; i < len(input); i++ {
		b, _ := NewBitSeqWithSize(uint(input[i]), 8)
		output = Concat(output, b)
	}
	return output
}

func getBestMode(data string) QRMode {
	isNumeric := true
	isAlphanumeric := true

	for _, r := range data {
		if r < '0' || r > '9' {
			isNumeric = false
		}
		if _, ok := alphanumericValues[r]; !ok {
			isAlphanumeric = false
		}
	}

	if isNumeric {
		return NumericMode
	}
	if isAlphanumeric {
		return AlphanumericMode
	}
	return ByteMode
}

func NewQRCode(r QRRequest) *qr {
	// Structure
	// a mode indicator, a character count indicator, data bits (input data + error correction) + padding

	// TODO: Determine best mode from input data
	mode := getBestMode(r.input)

	// Format
	format := QR_FORMAT_QR
	if r.is_micro {
		format = QR_FORMAT_MICRO_QR
	}

	// Encode input_data_bits
	input_data_bits := encode(mode, r.input)

	var version_num int
	if r.version != 0 {
		version_num = r.version
	} else {
		version_num = GetVersionNumber(mode, format, input_data_bits, errcorr(r.err_corr_level))
	}

	version := QRVersion{
		format: format,
		number: version_num,
	}

	character_count := character_count(mode, version, r.input)

	output := ConcatMany(GetModeIndicatorBits(version, mode), character_count, input_data_bits)

	// Add Padding
	// Calculate total data capacity in bytes
	dataInfo := capacityData[version.number]
	ecInfo := dataInfo.ecInfo[errcorr(r.err_corr_level)]
	totalECCodewords := 0
	for _, group := range ecInfo.BlockGroups {
		totalECCodewords += ecInfo.ECCodewordsPerBlock * group.NumBlocks
	}
	dataCapacityBytes := dataInfo.totalCodewords - totalECCodewords

	output = add_padding(output, dataCapacityBytes)

	print_bit_seq(output)

	// NOTE: hasta acá va bien el ejemplo ======

	// Calculate error correction codewords
	size := 21 + (version.number-1)*4
	matrix := make([][]module, size)
	isFunctionPattern := make([][]bool, size)
	for i := range size {
		matrix[i] = make([]module, size)
		isFunctionPattern[i] = make([]bool, size)
	}
	qr := &qr{
		matrix:              matrix,
		is_function_pattern: isFunctionPattern,
		version:             version,
		error_corr_level:    errcorr(r.err_corr_level),
		size:                size,
		data:                []byte(r.input),
		encoded_data:        output,
		mode:                mode,
		debug_no_mask:       r.debug_no_mask,
	}
	qr.generate()
	return qr
}

type errcorr string

const (
	ERR_CORR_L = "L"
	ERR_CORR_M = "M"
	ERR_CORR_Q = "Q"
	ERR_CORR_H = "H"
)

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

func (qr *qr) Draw(shape Shape) {
	pixs := make([][]color.Color, len(qr.matrix[0]))
	for y, row := range qr.matrix {
		imgRow := make([]color.Color, len(row))
		for x, c := range row {
			imgRow[x] = c.Color()
		}
		pixs[y] = imgRow
	}
	req := ImageRequest{
		Scale:  16,
		Pixels: pixs,
		Shape:  shape,
	}
	WriteImage(req)
}

type QRRequest struct {
	input          string
	is_micro       bool
	err_corr_level string
	version        int
	debug_no_mask  bool
}

func main() {
	args := os.Args[1:]

	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help") {
		fmt.Println("Usage: qr-decoder [ErrorCorrectionLevel] [Data] [Version] [IsMicro] [Shape]")
		fmt.Println("")
		fmt.Println("Arguments:")
		fmt.Println("  ErrorCorrectionLevel: L, M, Q, H (default: L)")
		fmt.Println("  Data: String to encode (default: \"01234567\")")
		fmt.Println("  Version: QR Code version 1-40 (optional, auto-detected if 0 or omitted)")
		fmt.Println("  IsMicro: true/false (default: false)")
		fmt.Println("  IsMicro: true/false (default: false)")
		fmt.Println("  Shape: square, circle, rounded, diamond (default: square)")
		fmt.Println("  --debug-no-mask: Disable masking for debugging (optional)")
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("  qr-decoder L \"Hello World\"")
		fmt.Println("  qr-decoder M \"1234567890\" 5 false circle")
		return
	}

	// Error correction level parsing
	var err_corr_level string
	if len(args) > 0 {
		if args[0] == ERR_CORR_L || args[0] == ERR_CORR_M || args[0] == ERR_CORR_Q || args[0] == ERR_CORR_H {
			err_corr_level = args[0]
		} else {
			panic("could not parse 'err_corr_level'")
		}
	} else {
		err_corr_level = "L"
	}

	// Data
	data := "01234567"
	if len(args) > 1 {
		data = args[1]
	}

	// Version
	var version int64
	if len(args) > 2 {
		version, _ = strconv.ParseInt(args[2], 10, 64)
	}

	var is_micro bool
	if len(args) > 3 {
		if v, err := strconv.ParseBool(args[3]); err == nil {
			is_micro = v
		} else {
			panic("could not parse 'is_micro'")
		}
	}

	var shape Shape = ShapeSquare
	if len(args) > 4 {
		shape = Shape(args[4])
	}

	var debug_no_mask bool
	for _, arg := range args {
		if arg == "--debug-no-mask" {
			debug_no_mask = true
			break
		}
	}

	req := QRRequest{
		input:          data,
		is_micro:       is_micro,
		err_corr_level: err_corr_level,
		version:        int(version),
		debug_no_mask:  debug_no_mask,
	}

	// qr := NewQRCode(version, err_corr_level, data)
	qr := NewQRCode(req)
	fmt.Println(qr.String())
	fmt.Printf("Mode: %s\n", qr.mode)
	formatInfo, _ := GenerateFormatInformation(qr.error_corr_level, qr.mask)
	fmt.Printf("Mask Pattern: %d (%015b)\n", qr.mask, formatInfo)
	qr.Draw(shape)

	// Dump matrix for comparison
	f, _ := os.Create("my_matrix.txt")
	defer f.Close()
	// Skip quiet zone (4 modules)
	start := 4
	end := qr.size - 4
	for i := start; i < end; i++ {
		for j := start; j < end; j++ {
			if qr.matrix[i][j].bit == One {
				f.WriteString("1")
			} else {
				f.WriteString("0")
			}
		}
		f.WriteString("\n")
	}
}
