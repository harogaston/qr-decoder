package bitseq

import (
	"fmt"
	"log/slog"
	"math/bits"
	"strings"
)

type BitOrder int

const (
	MSBFirst BitOrder = iota // Most Significant Bit first
	LSBFirst                 // Least Significant Bit first
)

type BitSeq struct {
	data  []byte // Internal storage (normalized to MSB)
	nbits int    // Total valid bits
}

// FromInt creates a sequence from a uint64 with a fixed bit size.
func FromInt(val uint64, size int) BitSeq {
	if size <= 0 {
		return BitSeq{}
	}
	if bits.Len64(val) > size {
		slog.Warn("FromInt: value exceeds specified size", "value", val, "size", size)
	}

	byteSize := (size + 7) / 8
	data := make([]byte, byteSize)

	for i := range size {
		// Check if the i-th bit (from MSB) of the value is 1.
		// The bits of `val` are read from most-significant to least-significant.
		if (val>>uint(size-1-i))&1 == 1 {
			// Calculate the position in the byte slice and set the bit.
			byteIdx := i / 8
			bitIdx := 7 - (i % 8)
			data[byteIdx] |= 1 << uint(bitIdx)
		}
	}

	return BitSeq{data: data, nbits: size}
}

// Append concatenates another BitSequence to the current one.
// It handles bit-shifting if the first sequence is not byte-aligned.
func (bs BitSeq) Append(other BitSeq) BitSeq {
	if other.nbits == 0 {
		return bs
	}

	newNBits := bs.nbits + other.nbits
	newSize := (newNBits + 7) / 8
	newData := make([]byte, newSize)
	copy(newData, bs.data)

	bitOffset := bs.nbits % 8
	if bitOffset == 0 {
		// Simple case: existing sequence ends exactly at a byte boundary
		copy(newData[len(bs.data):], other.data)
	} else {
		// Complex case: must shift the new bits to fit into the partial byte
		lastIdx := len(bs.data) - 1
		for _, b := range other.data {
			// Fill the remaining gap in the current last byte
			newData[lastIdx] |= b >> bitOffset
			lastIdx++
			if lastIdx < newSize {
				// Put the overflow into the next byte
				newData[lastIdx] = b << (8 - bitOffset)
			}
		}
	}

	return BitSeq{data: newData, nbits: newNBits}
}

// ConcatMany joins multiple BitSeq
func ConcatMany(seqs ...BitSeq) BitSeq {
	result := BitSeq{}
	for _, s := range seqs {
		result = result.Append(s)
	}
	return result
}

// Bytes exports the sequence as a byte slice with the requested bit ordering.
func (bs BitSeq) Bytes(order BitOrder) []byte {
	if len(bs.data) == 0 {
		return nil
	}

	if order == MSBFirst {
		// Internal format is already MSB; return a defensive copy
		res := make([]byte, len(bs.data))
		copy(res, bs.data)
		return res
	}

	// For LSBFirst, we reverse the bit order within each byte and the byte order
	res := make([]byte, len(bs.data))
	lastIdx := len(bs.data) - 1
	for i, b := range bs.data {
		res[lastIdx-i] = bits.Reverse8(b)
	}

	// Adjust padding for LSB if not a multiple of 8.
	// In LSB, valid bits are expected to be right-aligned (least significant).
	if bs.nbits%8 != 0 {
		shift := uint(8 - (bs.nbits % 8))
		var acc uint64
		for _, b := range res {
			acc = (acc << 8) | uint64(b)
		}
		acc >>= shift
		for i := range res {
			res[len(res)-1-i] = byte(acc >> (uint(i) * 8))
		}
	}

	return res
}

// Len returns the number of valid bits in the sequence.
func (bs BitSeq) Len() int {
	return bs.nbits
}

// Index 0 is the most significant bit (leftmost).
// It panics if the index is out of range.
func (bs BitSeq) Bit(index int) bool {
	if index < 0 || index >= bs.nbits {
		panic(fmt.Sprintf("bit index %d out of range [0-%d]", index, bs.nbits-1))
	}

	// 1. Identify which byte contains the bit
	byteIdx := index / 8
	// 2. Identify the position within that byte (MSB ordering)
	// For index 0, bitPos is 7 (1 << 7)
	bitPos := uint(7 - (index % 8))

	// 3. Extract the bit using a mask
	return (bs.data[byteIdx] & (1 << bitPos)) != 0
}

// ZeroSequence returns a sequence of 'n' bits set to 0.
func ZeroSequence(n int) BitSeq {
	return BitSeq{
		data:  make([]byte, (n+7)/8),
		nbits: n,
	}
}

// AlignToByte returns the number of bits needed to reach the next byte boundary.
func (bs BitSeq) AlignToByte() int {
	if bs.nbits%8 == 0 {
		return 0
	}
	return 8 - (bs.nbits % 8)
}

// String implements the fmt.Stringer interface.
// It returns a visual representation of the bits (e.g., "10110").
func (bs BitSeq) String() string {
	var sb strings.Builder
	sb.Grow(bs.nbits) // Pre-allocate memory for efficiency

	for i := 0; i < bs.nbits; i++ {
		if bs.Bit(i) {
			sb.WriteByte('1')
		} else {
			sb.WriteByte('0')
		}
	}

	return sb.String()
}
