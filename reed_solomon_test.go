package main

import (
	"testing"
)

func TestGFArithmetic(t *testing.T) {
	// Test primitive polynomial: 0x11D (285)
	// 2^8 = 256. In GF(2^8), 256 = 0x100.
	// 0x100 ^ 0x11D = 0x1D = 29.
	// So 2^8 should be 29 in our table if we generated it correctly?
	// Wait, expTable[8] should be 29? No.
	// expTable[8] is 2^8 mod primitive.
	// Let's check a few values.
	// 2^0 = 1
	// 2^1 = 2
	// ...
	// 2^7 = 128
	// 2^8 = 29 (0x1D)

	if expTable[0] != 1 {
		t.Errorf("expTable[0] = %d; want 1", expTable[0])
	}
	if expTable[1] != 2 {
		t.Errorf("expTable[1] = %d; want 2", expTable[1])
	}
	if expTable[8] != 29 {
		t.Errorf("expTable[8] = %d; want 29", expTable[8])
	}

	// Test multiplication
	// 2 * 3 = 6 (simple)
	if res := gfMul(2, 3); res != 6 {
		t.Errorf("gfMul(2, 3) = %d; want 6", res)
	}
	// 128 * 2 = 256 -> 29 (overflow)
	if res := gfMul(128, 2); res != 29 {
		t.Errorf("gfMul(128, 2) = %d; want 29", res)
	}

	// Test division
	// 6 / 3 = 2
	if res := gfDiv(6, 3); res != 2 {
		t.Errorf("gfDiv(6, 3) = %d; want 2", res)
	}
	// 29 / 2 = 128
	if res := gfDiv(29, 2); res != 128 {
		t.Errorf("gfDiv(29, 2) = %d; want 128", res)
	}
}

func TestRSEncode(t *testing.T) {
	// Test vector from a reliable source or manual calculation.
	// Example: QR Code Version 1-M.
	// Data codewords: 16
	// EC codewords: 10
	// Total: 26
	// Let's try a smaller example if possible, or use a known one.
	//
	// Example from Thonky's QR Code Tutorial:
	// Message: "hello world" (Byte mode)
	// We can just test the polynomial division logic with a simple case.
	//
	// Let's use the example from Wikipedia or similar for RS.
	// Or just trust the math if the GF arithmetic is correct.
	//
	// Let's try to encode "01234567" as in the main.go example.
	// We don't have the exact expected output without running a reference implementation.
	// However, we can check properties.
	//
	// Property: The resulting code (message + parity) should be divisible by the generator polynomial.
	// i.e. polyMod(code, generator) should be 0.

	data := []byte{0x10, 0x20, 0x0F} // Random data
	numEC := 7
	ec := reedSolomonEncode(data, numEC)

	if len(ec) != numEC {
		t.Errorf("Encoded length = %d; want %d", len(ec), numEC)
	}

	// Verify divisibility
	fullMsg := make([]int, len(data)+len(ec))
	for i, b := range data {
		fullMsg[i] = int(b)
	}
	for i, b := range ec {
		fullMsg[len(data)+i] = int(b)
	}

	msgPoly := newPolynomial(fullMsg)
	genPoly := generateGeneratorPolynomial(numEC)
	remPoly := polyMod(msgPoly, genPoly)

	// Remainder should be zero (or empty coeffs)
	if len(remPoly.coeffs) > 1 || (len(remPoly.coeffs) == 1 && remPoly.coeffs[0] != 0) {
		t.Errorf("RS check failed: remainder is not zero: %v", remPoly.coeffs)
	}
}
