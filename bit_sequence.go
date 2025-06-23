package main

import (
	"errors"
	"math/bits"
)

type bit_seq struct {
	value uint
	size  int
}

func NewBitSeq(i uint) bit_seq {
	return bit_seq{
		value: i,
		size:  bits.Len(i),
	}
}

func NewBitSeqWithSize(i uint, size int) (bit_seq, error) {
	if size < bits.Len(i) {
		return bit_seq{}, errors.New("invalid desired size")
	}
	return bit_seq{
		value: i,
		size:  size,
	}, nil
}

func Concat(first bit_seq, second bit_seq) bit_seq {
	if first.size+second.size > bits.UintSize {
		panic("cannot concat uints, exceeding maximum size")
	}
	return bit_seq{
		value: first.value<<second.size | second.value,
		size:  first.size + second.size,
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
