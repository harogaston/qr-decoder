package main

import (
	"testing"
)

func TestInterleaving(t *testing.T) {
	// Version 5-Q has 2 blocks.
	// Capacity: 62 data codewords.
	// Input: 70 chars (Byte mode) -> V5 should be enough?
	// V5-Q Total: 134, EC: 72 (2 blocks of 36). Data: 62.
	// Wait, 70 chars > 62. So it will go to V6?
	// V6-Q: Total 172, EC 96. Data: 76.
	// V6-Q has 4 blocks.

	input := "This is a long string to trigger multiple blocks for Reed-Solomon interleaving verification. It needs to be long enough."
	// Length: ~110 chars.

	// Create QR Code
	// This should not panic
	qr := NewQRCode(QRRequest{
		input_data:     input,
		err_corr_level: ERR_CORR_Q,
	})

	if qr.version.Number < 5 {
		t.Logf("Version selected: %d. Expected >= 5 for multi-block test.", qr.version.Number)
	} else {
		t.Logf("Version selected: %d. Multi-block test active.", qr.version.Number)
	}

	// We can't easily verify the matrix content without a decoder,
	// but successful execution implies the interleaving logic didn't crash
	// and produced a matrix.

	if len(qr.matrix) == 0 {
		t.Error("Matrix is empty")
	}
}

func TestAlphanumeric(t *testing.T) {
	input := "AC-42"
	// Should select Alphanumeric Mode

	qr := NewQRCode(QRRequest{
		input_data:     input,
		err_corr_level: ERR_CORR_L,
	})

	// We can't easily check the mode directly as it's internal to generate()
	// But we can check if it runs without panic.
	// And we can check if the version is small (V1).
	// "AC-42" is 5 chars.
	// Numeric: No.
	// Alphanumeric: Yes. 5 chars -> 9 bits length indicator (V1-9).
	// Data: 5 chars. 2 pairs + 1 single. 11+11+6 = 28 bits.
	// Mode indicator: 4 bits.
	// Total: 4 + 9 + 28 = 41 bits.
	// V1-L capacity: 19 bytes = 152 bits.
	// Fits easily in V1.

	if qr.version.Number != 1 {
		t.Errorf("Expected Version 1, got %d", qr.version.Number)
	}
}
