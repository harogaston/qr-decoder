package version

import (
	"fmt"
)

type QRFormat string

const FORMAT_QR_MODEL_2 = QRFormat("model2") // included in QR (2024)
const FORMAT_QR = QRFormat("qr")             // 2024 specification NOT IMPLEMENTED yet
const FORMAT_MICRO_QR = QRFormat("micro")    // NOT IMPLEMENTED yet

type QRVersion struct {
	Format QRFormat
	Number int
}

func (v QRVersion) String() string {
	var format string
	if v.Format == FORMAT_MICRO_QR {
		format = "M"
	}
	return fmt.Sprintf("%s%d", format, v.Number)
}
