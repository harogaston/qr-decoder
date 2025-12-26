package main

import (
	"errors"
)

const (
	// Format Info Mask Pattern: 0b101010000010010 (0x5412)
	format_information_mask_pattern = 0x5412
	// BCH(15, 5) Generator Polynomial: x^10 + x^8 + x^5 + x^4 + x^2 + x + 1
	// 10100110111 (0x537)
	format_information_generator_poly = 0x537
)

// GenerateFormatInformation calculates the 15-bit format information sequence
// containing 5 data bits (2 for error correction level, 3 for mask pattern)
// and 10 error correction bits.
// The result is XORed with the standard mask pattern.
func GenerateFormatInformation(ecLevel errcorr, maskPattern int) (uint16, error) {
	if maskPattern < 0 || maskPattern > 7 {
		return 0, errors.New("invalid mask pattern reference")
	}

	ecBits := get_err_corr_bits_for_level(ecLevel)

	// 2. Combine with Mask Pattern (3 bits)
	// Format: [EC Level (2 bits)] [Mask Pattern (3 bits)]
	data := (ecBits << 3) | uint(maskPattern)

	// 3. Calculate BCH(15, 5) Error Correction Bits (10 bits)
	bchBits := encodeBCH15_5(data, format_information_generator_poly)

	// 4. Assemble 15-bit sequence: [Data (5 bits)] [BCH (10 bits)]
	fullSequence := (data << 10) | bchBits

	// 5. XOR with Mask Pattern
	maskedSequence := fullSequence ^ format_information_mask_pattern

	return uint16(maskedSequence), nil
}

// encodeBCH15_5 calculates the BCH error correction bits.
// data: The data bits (5 bits for format info).
// poly: The generator polynomial.
func encodeBCH15_5(data uint, poly uint) uint {
	// We are calculating 10 parity bits for a (15, 5) code.
	// The input data is 5 bits.
	// We align the data to the MSB of the 15-bit field.
	// data << 10

	d := data << 10

	// Polynomial division
	// We need to find the remainder of d / poly in GF(2).
	// Since the generator polynomial is degree 10, and we want 10 parity bits,
	// we are essentially doing the division.

	// Iterate from the MSB of the data (bit 14 down to bit 10).
	// The data is in the upper 5 bits of the 15-bit register 'd'.
	// The generator poly is 11 bits (degree 10).

	// Let's look at how many shifts we need.
	// The data is 5 bits. We shifted it left by 10.
	// So the significant bits are at 14..10.
	// The generator poly's MSB is at bit 10 (if we treat it as 11 bits 10..0).
	// Wait, 0x537 is 10100110111. That's 11 bits. MSB is bit 10.

	// We want to eliminate the 1s in the data part.
	for i := range 5 {
		// Check if the MSB of the current remainder is 1.
		// The MSB we are interested in is bit (14 - i).
		if (d>>(14-i))&1 == 1 {
			// XOR with poly shifted to align with the current MSB.
			// Poly is degree 10. Its MSB is at bit 10.
			// We want to align it with bit (14 - i).
			// So shift amount = (14 - i) - 10 = 4 - i.
			d ^= poly << (4 - i)
		}
	}

	// The remainder is the lower 10 bits of d.
	return d & 0x3FF // Mask to 10 bits
}
