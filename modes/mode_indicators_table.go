package modes

import (
	"github.com/harogaston/qr-decoder/bitseq"
	"github.com/harogaston/qr-decoder/version"
)

// Constants for specific versions
type modeIndicatorVersionClass string

const (
	M1         modeIndicatorVersionClass = "M1"
	M2         modeIndicatorVersionClass = "M2"
	M3         modeIndicatorVersionClass = "M3"
	M4         modeIndicatorVersionClass = "M4"
	AllQRCodes modeIndicatorVersionClass = "all"
)

type modeIndicatorMap map[modeIndicatorVersionClass]bitseq.BitSeq

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
		M1:         bitseq.BitSeq{},
		M2:         bitseq.BitSeq{},
		M3:         bitseq.BitSeq{},
		M4:         bitseq.BitSeq{},
		AllQRCodes: bitseq.FromInt(7, 4),
	},
	NumericMode: {
		M1:         bitseq.BitSeq{},
		M2:         bitseq.FromInt(0, 1),
		M3:         bitseq.FromInt(0, 2),
		M4:         bitseq.FromInt(0, 3),
		AllQRCodes: bitseq.FromInt(1, 4),
	},
	AlphanumericMode: {
		M1:         bitseq.BitSeq{},
		M2:         bitseq.FromInt(1, 1),
		M3:         bitseq.FromInt(1, 2),
		M4:         bitseq.FromInt(1, 3),
		AllQRCodes: bitseq.FromInt(2, 4),
	},
	ByteMode: {
		M1:         bitseq.BitSeq{},
		M2:         bitseq.BitSeq{},
		M3:         bitseq.FromInt(2, 2),
		M4:         bitseq.FromInt(2, 3),
		AllQRCodes: bitseq.FromInt(4, 4),
	},
	KanjiMode: {
		M1:         bitseq.BitSeq{},
		M2:         bitseq.BitSeq{},
		M3:         bitseq.FromInt(3, 2),
		M4:         bitseq.FromInt(3, 3),
		AllQRCodes: bitseq.FromInt(8, 4),
	},
	StructuredAppend: {
		M1:         bitseq.BitSeq{},
		M2:         bitseq.BitSeq{},
		M3:         bitseq.BitSeq{},
		M4:         bitseq.BitSeq{},
		AllQRCodes: bitseq.FromInt(3, 4),
	},
}

// GetModeIndicatorBits returns the mode indicator bits for a given QR version and mode.
func GetModeIndicatorBits(qrversion version.QRVersion, mode QRMode) bitseq.BitSeq {
	if qrversion.Format == version.FORMAT_MICRO_QR {
		v := qrversion.Number
		switch v {
		case 1:
			return modeIndicatorData[mode][M1]
		case 2:
			return modeIndicatorData[mode][M2]
		case 3:
			return modeIndicatorData[mode][M3]
		case 4:
			return modeIndicatorData[mode][M4]
		}
	}
	if qrversion.Format == version.FORMAT_QR || qrversion.Format == version.FORMAT_QR_MODEL_2 {
		return modeIndicatorData[mode][AllQRCodes]
	}

	return bitseq.BitSeq{}
}
