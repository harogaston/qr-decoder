package main

import (
	"errors"
	"math/bits"
)

type bit_seq struct {
	data []byte // Bits packed into bytes. data[0] contains MSB.
	len  int    // Number of valid bits
}

func NewBitSeq(i uint) bit_seq {
	size := bits.Len(i)
	return createBitSeq(uint64(i), size)
}

func NewBitSeqWithSize(i uint, size int) (bit_seq, error) {
	if size < bits.Len(i) {
		return bit_seq{}, errors.New("invalid desired size")
	}
	return createBitSeq(uint64(i), size), nil
}

func MustNewBitSeqWithSize(i uint, size int) bit_seq {
	bs, err := NewBitSeqWithSize(i, size)
	if err != nil {
		panic(err)
	}
	return bs
}

func createBitSeq(val uint64, size int) bit_seq {
	numBytes := (size + 7) / 8
	data := make([]byte, numBytes)

	// Fill data
	// We want the bits to be placed such that concatenation works easily.
	// Let's adopt a convention:
	// data[0] contains the first 8 bits (MSB first).
	// Example: size=4, val=0b1010. data=[0b10100000] (padded with 0s at end? or start?)
	// Usually for streams, we append bits.
	// Let's say we fill from MSB to LSB.
	// If size=4, val=10 (1010).
	// We want to append 1, 0, 1, 0.

	for i := 0; i < size; i++ {
		if (val>>(size-1-i))&1 == 1 {
			byteIdx := i / 8
			bitIdx := 7 - (i % 8)
			data[byteIdx] |= 1 << bitIdx
		}
	}

	return bit_seq{
		data: data,
		len:  size,
	}
}

func Concat(first bit_seq, second bit_seq) bit_seq {
	newLen := first.len + second.len
	numBytes := (newLen + 7) / 8
	newData := make([]byte, numBytes)

	// Copy first
	// We can't just copy bytes because of alignment.
	// We have to copy bits.
	// Optimization: copy bytes if aligned?
	// For now, bit by bit copy is safer and easier to implement correctly.

	setBit := func(idx int, val int) {
		if val == 1 {
			byteIdx := idx / 8
			bitIdx := 7 - (idx % 8)
			newData[byteIdx] |= 1 << bitIdx
		}
	}

	getBit := func(bs bit_seq, idx int) int {
		byteIdx := idx / 8
		bitIdx := 7 - (idx % 8)
		if byteIdx >= len(bs.data) {
			return 0
		}
		if (bs.data[byteIdx]>>bitIdx)&1 == 1 {
			return 1
		}
		return 0
	}

	for i := 0; i < first.len; i++ {
		setBit(i, getBit(first, i))
	}
	for i := 0; i < second.len; i++ {
		setBit(first.len+i, getBit(second, i))
	}

	return bit_seq{
		data: newData,
		len:  newLen,
	}
}

func ConcatMany(seqs ...bit_seq) bit_seq {
	if len(seqs) == 0 {
		panic("call to Concat with empty collection")
	}
	acc := seqs[0]
	for _, bs := range seqs[1:] {
		acc = Concat(acc, bs)
	}
	return acc
}

func (bs *bit_seq) Get(index int) int {
	// Cast index to int if it's uint in caller?
	// The caller uses uint. Let's change signature to match usage or cast.
	// But `main.go` calls it with `uint` loop variable `i`.
	// Wait, `main.go` loop: `for i := uint(0); i < bs.len; i++`
	// bs.len is int.
	// So i is uint.
	// I should accept int.

	byteIdx := index / 8
	bitIdx := 7 - (index % 8)
	if byteIdx >= len(bs.data) {
		return 0
	}
	if (bs.data[byteIdx]>>bitIdx)&1 == 1 {
		return 1
	}
	return 0
}

// ToBytes returns the byte slice, assuming the sequence is byte-aligned and full.
// If it's not multiple of 8, the last byte is padded with zeros at the end.
func (bs *bit_seq) ToBytes() []byte {
	// Return a copy to avoid modification
	res := make([]byte, len(bs.data))
	copy(res, bs.data)
	return res
}
