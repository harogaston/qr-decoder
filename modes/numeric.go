package modes

import (
	"github.com/harogaston/qr-decoder/bitseq"
	"strconv"
)

func EncodeNumeric(input string) bitseq.BitSeq {
	var output bitseq.BitSeq
	var prev int

	for prev < len(input) {
		if len(input) >= prev+3 { // Encode in 10 bits
			group := input[prev : prev+3]
			group_uint64, err := strconv.ParseUint(group, 10, 10)
			if err != nil {
				panic(err)
			}
			b := bitseq.FromInt(group_uint64, 10)
			output = output.Append(b)
			prev = prev + 3
		} else if len(input) == prev+2 { // Encode in 7 bits
			group := input[prev : prev+2]
			group_uint64, err := strconv.ParseUint(group, 10, 7)
			if err != nil {
				panic(err)
			}
			b := bitseq.FromInt(group_uint64, 7)
			output = output.Append(b)
			prev = prev + 2
		} else if len(input) == prev+1 { // Encode in 4 bits
			group := input[prev : prev+1]
			group_uint64, err := strconv.ParseUint(group, 10, 4)
			if err != nil {
				panic(err)
			}
			b := bitseq.FromInt(group_uint64, 4)
			output = output.Append(b)
			prev = prev + 1
		}
	}
	return output
}
