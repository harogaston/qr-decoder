package modes

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

// getMode follows a simple hierarchy. It checks input_data against
// the character sets of each mode in order of most to least "compressed."
// TODO: Add Kanji mode detection and mode switching
func GetMode(data string) QRMode {
	isNumeric := true
	isAlphanumeric := true

	for _, r := range data {
		if r < '0' || r > '9' {
			isNumeric = false
		}
		if _, ok := alphanumericValues[r]; !ok {
			isAlphanumeric = false
		}
	}

	if isNumeric {
		return NumericMode
	}
	if isAlphanumeric {
		return AlphanumericMode
	}
	return ByteMode
}
