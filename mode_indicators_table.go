package main

// Constants for specific versions
type modeIndicatorVersionClass string

const (
	M1         modeIndicatorVersionClass = "M1"
	M2         modeIndicatorVersionClass = "M2"
	M3         modeIndicatorVersionClass = "M3"
	M4         modeIndicatorVersionClass = "M4"
	AllQRCodes modeIndicatorVersionClass = "all"
)

type modeIndicatorMap map[modeIndicatorVersionClass]bit_seq

// | Version   | Numeric mode | Alphanumeric mode | Byte mode | Kanji mode |
// | :-------- | :----------- | :---------------- | :-------- | :--------- |
// | M1        | 3            | NA                | NA        | NA         |
// | M2        | 4            | 3                 | NA        | NA         |
// | M3        | 5            | 4                 | 4         | 3          |
// | M4        | 6            | 5                 | 5         | 4          |
// | 1 to 9    | 10           | 9                 | 8         | 8          |
// | 10 to 26  | 12           | 11                | 16        | 10         |
// | 27 to 40  | 14           | 13                | 16        | 12         |
var modeIndicatorData = map[QRMode]modeIndicatorMap{
	ECI: {
		M1:         bit_seq{},
		M2:         bit_seq{},
		M3:         bit_seq{},
		M4:         bit_seq{},
		AllQRCodes: MustNewBitSeqWithSize(7, 4),
	},
	NumericMode: {
		M1:         bit_seq{},
		M2:         MustNewBitSeqWithSize(0, 1),
		M3:         MustNewBitSeqWithSize(0, 2),
		M4:         MustNewBitSeqWithSize(0, 3),
		AllQRCodes: MustNewBitSeqWithSize(1, 4),
	},
	AlphanumericMode: {
		M1:         bit_seq{},
		M2:         MustNewBitSeqWithSize(1, 1),
		M3:         MustNewBitSeqWithSize(1, 2),
		M4:         MustNewBitSeqWithSize(1, 3),
		AllQRCodes: MustNewBitSeqWithSize(2, 4),
	},
	ByteMode: {
		M1:         bit_seq{},
		M2:         bit_seq{},
		M3:         MustNewBitSeqWithSize(2, 2),
		M4:         MustNewBitSeqWithSize(2, 3),
		AllQRCodes: MustNewBitSeqWithSize(4, 4),
	},
	KanjiMode: {
		M1:         bit_seq{},
		M2:         bit_seq{},
		M3:         MustNewBitSeqWithSize(3, 2),
		M4:         MustNewBitSeqWithSize(3, 3),
		AllQRCodes: MustNewBitSeqWithSize(8, 4),
	},
	StructuredAppend: {
		M1:         bit_seq{},
		M2:         bit_seq{},
		M3:         bit_seq{},
		M4:         bit_seq{},
		AllQRCodes: MustNewBitSeqWithSize(3, 4),
	},
}

// GetModeIndicatorBits returns the mode indicator bits for a given QR version and mode.
func GetModeIndicatorBits(version QRVersion, mode QRMode) bit_seq {
	if version.format == QR_FORMAT_MICRO_QR {
		v := version.number
		switch {
		case v == 1:
			return modeIndicatorData[mode][M1]
		case v == 2:
			return modeIndicatorData[mode][M2]
		case v == 3:
			return modeIndicatorData[mode][M3]
		case v == 4:
			return modeIndicatorData[mode][M4]
		}
	}
	if version.format == QR_FORMAT_QR {
		return modeIndicatorData[mode][AllQRCodes]
	}

	return bit_seq{}
}
