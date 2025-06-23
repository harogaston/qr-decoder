package main

// QRMode represents the different QR code data modes.
type QRMode int

const (
	ECI QRMode = iota
	NumericMode
	AlphanumericMode
	ByteMode
	KanjiMode
	StructuredAppend
	UnknownMode // Default or error case
)

// String method for QRMode for better readability
func (m QRMode) String() string {
	switch m {
	case ECI:
		return "ECI"
	case NumericMode:
		return "Numeric"
	case AlphanumericMode:
		return "Alphanumeric"
	case ByteMode:
		return "Byte"
	case KanjiMode:
		return "Kanji"
	case StructuredAppend:
		return "StructuredAppend"
	default:
		return "Unknown"
	}
}
