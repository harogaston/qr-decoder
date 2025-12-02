package main

import "fmt"

type QRFormat string

const QR_FORMAT_MODEL_1 = "model1"
const QR_FORMAT_MODEL_2 = "model2"  // included in QR (2024)
const QR_FORMAT_QR = QRFormat("qr") // 2024 specification
const QR_FORMAT_MICRO_QR = QRFormat("micro")

type QRVersion struct {
	format QRFormat
	number int
}

func (v QRVersion) String() string {
	var format string
	if v.format == QR_FORMAT_MICRO_QR {
		format = "M"
	}
	return fmt.Sprintf("%s%d", format, v.number)
}

func (v QRVersion) TerminatorBits() bit_seq {
	panic("not implemented")
}
