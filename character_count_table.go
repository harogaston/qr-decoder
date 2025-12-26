package main

import (
	"github.com/harogaston/qr-decoder/bitseq"
	"github.com/harogaston/qr-decoder/modes"
	"github.com/harogaston/qr-decoder/version"
)

// Constants for specific versions
type charCountVersionClass string

const (
	M1Version      charCountVersionClass = "M1"
	M2Version      charCountVersionClass = "M2"
	M3Version      charCountVersionClass = "M3"
	M4Version      charCountVersionClass = "M4"
	V1To9Version   charCountVersionClass = "1 to 9"
	V10To26Version charCountVersionClass = "10 to 26"
	V27To40Version charCountVersionClass = "27 to 40"
)

type CharCountDataLength map[charCountVersionClass]int

// | Version   | Numeric mode | Alphanumeric mode | Byte mode | Kanji mode |
// | :-------- | :----------- | :---------------- | :-------- | :--------- |
// | M1        | 3            | NA                | NA        | NA         |
// | M2        | 4            | 3                 | NA        | NA         |
// | M3        | 5            | 4                 | 4         | 3          |
// | M4        | 6            | 5                 | 5         | 4          |
// | 1 to 9    | 10           | 9                 | 8         | 8          |
// | 10 to 26  | 12           | 11                | 16        | 10         |
// | 27 to 40  | 14           | 13                | 16        | 12         |
var charCountData = map[modes.QRMode]CharCountDataLength{
	modes.NumericMode: {
		M1Version:      3,
		M2Version:      4,
		M3Version:      5,
		M4Version:      6,
		V1To9Version:   10,
		V10To26Version: 12,
		V27To40Version: 14,
	},
	modes.AlphanumericMode: {
		M1Version:      0,
		M2Version:      3,
		M3Version:      4,
		M4Version:      5,
		V1To9Version:   9,
		V10To26Version: 11,
		V27To40Version: 13,
	},
	modes.ByteMode: {
		M1Version:      0,
		M2Version:      0,
		M3Version:      4,
		M4Version:      5,
		V1To9Version:   8,
		V10To26Version: 16,
		V27To40Version: 16,
	},
	modes.KanjiMode: {
		M1Version:      0,
		M2Version:      0,
		M3Version:      3,
		M4Version:      4,
		V1To9Version:   8,
		V10To26Version: 10,
		V27To40Version: 12,
	},
}

func GetVersionNumber(mode modes.QRMode, format version.QRFormat, data bitseq.BitSeq, ecLevel errcorr) int {
	switch format {
	case version.FORMAT_QR, version.FORMAT_QR_MODEL_2:
		for num := 1; num <= 40; num++ {
			v := version.QRVersion{Format: format, Number: num}
			// 1. Mode indicator length
			modeBits := 4 // For QR Code (except Micro)

			// 2. Char count indicator length
			charCountBits := GetCharCountLength(v, mode)

			// 3. Total bits
			totalBits := modeBits + charCountBits + data.Len()

			// 4. Data capacity
			dataCapacityBits := (getTotalDataCodewords(v, ecLevel)) * 8

			if totalBits <= dataCapacityBits {
				return num
			}
		}
	case version.FORMAT_MICRO_QR:
		// TODO: Implement Micro QR capacity check
		panic("Micro QR not implemented yet")
	}
	return 0
}

// GetCharCountLength retrieves the character count for a given QR version and mode.
// Returns 0 for N/A cases.
func GetCharCountLength(qrversion version.QRVersion, mode modes.QRMode) int {
	if qrversion.Format == version.FORMAT_MICRO_QR {
		v := qrversion.Number
		switch {
		case v == 1:
			return charCountData[mode][M1Version]
		case v == 2:
			return charCountData[mode][M2Version]
		case v == 3:
			return charCountData[mode][M3Version]
		case v == 4:
			return charCountData[mode][M4Version]
		}
	}
	if qrversion.Format == version.FORMAT_QR || qrversion.Format == version.FORMAT_QR_MODEL_2 {
		v := qrversion.Number
		switch {
		case v >= 1 && v <= 9:
			return charCountData[mode][V1To9Version]
		case v >= 10 && v <= 26:
			return charCountData[mode][V10To26Version]
		case v >= 27 && v <= 40:
			return charCountData[mode][V27To40Version]
		}
	}

	return 0
}
