package modes

import (
	"github.com/harogaston/qr-decoder/bitseq"
	"github.com/harogaston/qr-decoder/version"
)

func GetTerminatorBits(qrversion version.QRVersion, mode QRMode) bitseq.BitSeq {
	// TODO: Implement different terminator bits for Micro QR Codes
	return bitseq.FromInt(0, 4)
}
