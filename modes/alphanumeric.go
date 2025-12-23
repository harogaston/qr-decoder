package modes

import "github.com/harogaston/qr-decoder/bitseq"

var alphanumericValues = map[rune]int{
	'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
	'A': 10, 'B': 11, 'C': 12, 'D': 13, 'E': 14, 'F': 15, 'G': 16, 'H': 17, 'I': 18, 'J': 19,
	'K': 20, 'L': 21, 'M': 22, 'N': 23, 'O': 24, 'P': 25, 'Q': 26, 'R': 27, 'S': 28, 'T': 29,
	'U': 30, 'V': 31, 'W': 32, 'X': 33, 'Y': 34, 'Z': 35,
	' ': 36, '$': 37, '%': 38, '*': 39, '+': 40, '-': 41, '.': 42, '/': 43, ':': 44,
}

func EncodeAlphanumeric(input string) bitseq.BitSeq {
	var output bitseq.BitSeq
	for i := 0; i < len(input); i += 2 {
		if i+1 < len(input) {
			// Pair
			val1 := alphanumericValues[rune(input[i])]
			val2 := alphanumericValues[rune(input[i+1])]
			val := 45*val1 + val2
			b := bitseq.FromInt(uint64(val), 11)
			output = output.Append(b)
		} else {
			// Single
			val := alphanumericValues[rune(input[i])]
			b := bitseq.FromInt(uint64(val), 6)
			output = output.Append(b)
		}
	}
	return output
}
