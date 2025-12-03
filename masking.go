package main

import (
	"math"
)

// Mask patterns
// 000: (i + j) mod 2 = 0
// 001: i mod 2 = 0
// 010: j mod 3 = 0
// 011: (i + j) mod 3 = 0
// 100: ((i div 2) + (j div 3)) mod 2 = 0
// 101: (i j) mod 2 + (i j) mod 3 = 0
// 110: ((i j) mod 2 + (i j) mod 3) mod 2 = 0
// 111: ((i+j) mod 2 + (i j) mod 3) mod 2 = 0

type maskFunc func(i, j int) bool

var maskPatterns = []maskFunc{
	func(i, j int) bool { return (i+j)%2 == 0 },
	func(i, j int) bool { return i%2 == 0 },
	func(i, j int) bool { return j%3 == 0 },
	func(i, j int) bool { return (i+j)%3 == 0 },
	func(i, j int) bool { return ((i/2)+(j/3))%2 == 0 },
	func(i, j int) bool { return (i*j)%2+(i*j)%3 == 0 },
	func(i, j int) bool { return ((i*j)%2+(i*j)%3)%2 == 0 },
	func(i, j int) bool { return ((i+j)%2+(i*j)%3)%2 == 0 },
}

// Apply mask to the matrix.
// Note: Masking is NOT applied to function patterns.
// We need a way to know which modules are function patterns.
// Usually, we mark them or keep a separate "is_function" map.
func (qr *qr) apply_mask(maskIndex int, matrix [][]module) [][]module {
	// Create a copy of the matrix
	size := len(matrix)
	maskedMatrix := make([][]module, size)
	for i := range size {
		maskedMatrix[i] = make([]module, size)
		copy(maskedMatrix[i], matrix[i])
	}

	maskFn := maskPatterns[maskIndex]

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			// Skip function patterns
			if qr.isFunctionPattern(i, j) {
				continue
			}

			if maskFn(i, j) {
				// Flip the bit
				if maskedMatrix[i][j].bit == Zero {
					maskedMatrix[i][j].bit = One
				} else if maskedMatrix[i][j].bit == One {
					maskedMatrix[i][j].bit = Zero
				}
			}
		}
	}
	return maskedMatrix
}

// Helper to identify function patterns
// This is tricky because we need to know exactly where they are.
// Ideally, we should have marked them during generation.
// For now, let's implement a check based on coordinates.
func (qr *qr) isFunctionPattern(row, col int) bool {
	return qr.is_function_pattern[row][col]
}

// Penalty calculation
func calculatePenalty(matrix [][]module) int {
	penalty := 0
	size := len(matrix)

	// Rule 1: 5+ same color modules in a row/col
	// Rows
	for i := 0; i < size; i++ {
		runLen := 1
		lastBit := matrix[i][0].bit
		for j := 1; j < size; j++ {
			if matrix[i][j].bit == lastBit {
				runLen++
			} else {
				if runLen >= 5 {
					penalty += 3 + (runLen - 5)
				}
				runLen = 1
				lastBit = matrix[i][j].bit
			}
		}
		if runLen >= 5 {
			penalty += 3 + (runLen - 5)
		}
	}
	// Cols
	for j := 0; j < size; j++ {
		runLen := 1
		lastBit := matrix[0][j].bit
		for i := 1; i < size; i++ {
			if matrix[i][j].bit == lastBit {
				runLen++
			} else {
				if runLen >= 5 {
					penalty += 3 + (runLen - 5)
				}
				runLen = 1
				lastBit = matrix[i][j].bit
			}
		}
		if runLen >= 5 {
			penalty += 3 + (runLen - 5)
		}
	}

	// Rule 2: 2x2 blocks of same color
	for i := 0; i < size-1; i++ {
		for j := 0; j < size-1; j++ {
			bit := matrix[i][j].bit
			if matrix[i+1][j].bit == bit && matrix[i][j+1].bit == bit && matrix[i+1][j+1].bit == bit {
				penalty += 3
			}
		}
	}

	// Rule 3: 1:1:3:1:1 pattern (dark:light:dark:light:dark)
	// Pattern: 1 0 1 1 1 0 1 (Dark Light Dark Dark Dark Light Dark)
	// Wait, the spec says 1:1:3:1:1 ratio.
	// Dark, Light, Dark, Light, Dark
	// 1, 1, 3, 1, 1
	// And 4 light modules on either side.
	// So: 0000 1011101 or 1011101 0000
	// Let's check for 1 0 1 1 1 0 1 sequence.
	// And check surrounding.

	// Rule 4: Dark module percentage
	darkCount := 0
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if matrix[i][j].bit == One {
				darkCount++
			}
		}
	}
	total := size * size
	percent := float64(darkCount) / float64(total) * 100

	// |percent - 50|
	// Spec: Calculate |prev5 - 50| / 5 and |next5 - 50| / 5. Take min?
	// Actually: 10 * k, where k is step of 5% deviation from 50%.
	// e.g. 45% or 55% -> k=1 -> 10 points.
	// 43% -> 40% (k=2) or 45% (k=1)?
	// Spec says: Count dark modules. Calculate percentage.
	// k = abs(percent - 50) / 5 (integer division? No, steps of 5)
	// "Determine the rating for each deviation from 50% in steps of 5%."
	// "N * 10" where N is the number of steps.
	// Example: 43%. Steps: 45-50 (1), 40-45 (2).
	// Usually implemented as:
	// deviation = abs(percent - 50)
	// score = floor(deviation / 5) * 10

	deviation := math.Abs(percent - 50)
	score4 := int(deviation/5) * 10
	penalty += score4

	return penalty
}
