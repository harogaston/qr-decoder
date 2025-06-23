package main

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
var charCountData = map[QRMode]CharCountDataLength{
	NumericMode: {
		M1Version:      3,
		M2Version:      4,
		M3Version:      5,
		M4Version:      6,
		V1To9Version:   10,
		V10To26Version: 12,
		V27To40Version: 14,
	},
	AlphanumericMode: {
		M1Version:      0,
		M2Version:      3,
		M3Version:      4,
		M4Version:      5,
		V1To9Version:   9,
		V10To26Version: 11,
		V27To40Version: 13,
	},
	ByteMode: {
		M1Version:      0,
		M2Version:      0,
		M3Version:      4,
		M4Version:      5,
		V1To9Version:   8,
		V10To26Version: 16,
		V27To40Version: 16,
	},
	KanjiMode: {
		M1Version:      0,
		M2Version:      0,
		M3Version:      3,
		M4Version:      4,
		V1To9Version:   8,
		V10To26Version: 10,
		V27To40Version: 12,
	},
}

// GetCharCountLength retrieves the character count for a given QR version and mode.
// Returns 0 for N/A cases.
func GetCharCountLength(version QRVersion, mode QRMode) int {
	if version.format == QR_FORMAT_MICRO_QR {
		v := version.number
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
	if version.format == QR_FORMAT_QR {
		v := version.number
		switch {
		case v <= 1 && v <= 9:
			return charCountData[mode][V1To9Version]
		case v >= 10 && v <= 26:
			return charCountData[mode][V10To26Version]
		case v >= 27 && v <= 40:
			return charCountData[mode][V27To40Version]
		}
	}

	return 0
}
