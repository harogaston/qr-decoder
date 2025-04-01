package main

import (
	"math/bits"
)

// Returns a slice of modules representing 'num' starting from
// the most significant bit
func modules_from_int(num int, size int, sort bitsorting) []module {
	var res []module

	// The number of bits needed to represent 'num'
	numBits := bits.Len(uint(num))

	if sort == big_endian {
		for i := numBits - 1; i >= 0; i-- {
			// Check if the bit at position i is 1 or 0
			if num&(1<<i) != 0 {
				res = append(res, module{bit: One})
			} else {
				res = append(res, module{bit: Zero})
			}
		}
	} else {
		for range numBits {
			if num&1 == 1 {
				res = append(res, module{bit: One})
			} else {
				res = append(res, module{bit: Zero})
			}
			// Shift right by 1 to check the next bit
			num >>= 1
		}
	}

	// Adjust for desired size
	padding := make([]module, size-numBits)
	for i := range padding {
		padding[i] = module{bit: Zero}
	}

	if sort == big_endian {
		return append(res, padding...)
	} else {
		return append(padding, res...)
	}
}
