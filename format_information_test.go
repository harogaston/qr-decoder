package main

import (
	"testing"
)

func TestGenerateFormatInformation(t *testing.T) {
	tests := []struct {
		name        string
		ecLevel     errcorr
		maskPattern int
		expected    uint16
		wantErr     bool
	}{
		{
			// Example from ISO/IEC 18004:2024 Section 7.9.1
			// Error correction level M (00)
			// Mask pattern 5 (101)
			// Data: 00101
			// BCH bits: 0011011100
			// Full sequence: 001010011011100
			// Mask: 101010000010010
			// Result: 100000011001110 (0x4066)
			name:        "ISO Spec Example (Level M, Mask 5)",
			ecLevel:     ERR_CORR_M,
			maskPattern: 5,
			expected:    0b100000011001110, // 0x4066
			wantErr:     false,
		},
		{
			// Level L (01), Mask 0 (000)
			// Data: 01000
			// BCH: 1011110111 (calculated manually or via trusted tool)
			// Full: 010001011110111
			// Mask: 101010000010010
			// Result: 0x77C4 (Calculated: Data 01000 -> BCH 0x3D6 -> Full 0x23D6 -> XOR Mask 0x5412 -> 0x77C4)
			name:        "Level L, Mask 0",
			ecLevel:     ERR_CORR_L,
			maskPattern: 0,
			expected:    0x77C4,
			wantErr:     false,
		},
		{
			// Invalid Mask Pattern
			name:        "Invalid Mask Pattern",
			ecLevel:     ERR_CORR_L,
			maskPattern: 8,
			expected:    0,
			wantErr:     true,
		},
		{
			// Invalid Error Correction Level
			name:        "Invalid EC Level",
			ecLevel:     "INVALID",
			maskPattern: 0,
			expected:    0,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateFormatInformation(tt.ecLevel, tt.maskPattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateFormatInformation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.expected {
				t.Errorf("GenerateFormatInformation() = 0x%X, want 0x%X", got, tt.expected)
			}
		})
	}
}

func TestCalculateBCH(t *testing.T) {
	// Test calculateBCH with the example from the spec
	// Data: 00101 (5)
	// Poly: 10100110111 (0x537)
	// Expected BCH: 0011011100 (0x0DC)

	data := uint(0b00101)
	poly := uint(0x537)
	expected := uint(0b0011011100)

	got := encodeBCH15_5(data, poly)
	if got != expected {
		t.Errorf("calculateBCH(0x%X, 0x%X) = 0x%X, want 0x%X", data, poly, got, expected)
	}
}
