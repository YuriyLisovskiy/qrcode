package qr

import (
	"math"
	"strings"
)

// Represents a character string to be encoded in a QR Code symbol. Each segment has
// a mode, and a sequence of characters that is already encoded as a sequence of bits.
// Instances of this class are immutable.
// This segment class imposes no length restrictions, but QR Codes have restrictions.
// Even in the most favorable conditions, a QR Code can only hold 7089 characters of data.
// Any segment longer than this is meaningless for the purpose of generating QR Codes.
type qrSegment struct {
	// The mode indicator for this segment.
	Mode modeType

	// The length of this segment's unencoded data, measured in characters. Always zero or positive.
	NumChars int

	// The data bits of this segment.
	Data []bool
}

// Returns a segment representing the given binary data encoded in byte mode.
func makeBytes(data *[]uint8) qrSegment {
	if len(*data) > math.MaxInt32 {
		panic(qrSegmentErr("makeBytes", "data too long"))
	}
	bitBuf := bitBuffer{}
	for _, bit := range *data {
		bitBuf = bitBuf.appendBits(uint32(bit), 8)
	}
	return qrSegment{isBYTE, len(*data), bitBuf}
}

// Returns a segment representing the given string of decimal digits encoded in numeric mode.
func makeNumeric(digits string) qrSegment {
	bitBuf := bitBuffer{}
	accumData, accumCount, charCount := 0, 0, 0
	for _, digit := range digits {
		if digit < '0' || digit > '9' {
			panic(qrSegmentErr("makeNumeric", "string contains non-numeric characters"))
		}
		accumData = accumData*10 + (int(digit) - '0')
		accumCount++
		if accumCount == 3 {
			bitBuf = bitBuf.appendBits(uint32(accumData), 10)
			accumData, accumCount = 0, 0
		}
		charCount++
	}
	if accumCount > 0 {
		bitBuf = bitBuf.appendBits(uint32(accumData), accumCount*3+1)
	}
	return qrSegment{isNUMERIC, charCount, bitBuf}
}

// Returns a segment representing the given text string encoded in alphanumeric mode.
//
// The characters allowed are: 0 to 9, A to Z (uppercase only), space,
// dollar, percent, asterisk, plus, hyphen, period, slash, colon.
func makeAlphanumeric(text string) qrSegment {
	bitBuf := bitBuffer{}
	accumData, accumCount, charCount := 0, 0, 0
	for _, char := range text {
		if !strings.ContainsRune(alphanumericCharset, char) {
			panic(qrSegmentErr("makeAlphanumeric", "string contains unencodable characters in alphanumeric mode"))
		}
		accumData = accumData*45 + (len(alphanumericCharset[strings.IndexRune(alphanumericCharset, char):])-len(alphanumericCharset))*-1
		accumCount++
		if accumCount == 2 {
			bitBuf = bitBuf.appendBits(uint32(accumData), 11)
			accumData, accumCount = 0, 0
		}
		charCount++
	}
	if accumCount > 0 {
		bitBuf = bitBuf.appendBits(uint32(accumData), 6)
	}
	return qrSegment{isALPHANUMERIC, charCount, bitBuf}
}

// Returns a list of zero or more segments to represent the given text string.
//
// The result may use various segment modes and switch modes to optimize the length of the bit stream.
func makeSegments(text string) []qrSegment {
	var result []qrSegment
	if text == "" {

	} else if isNumeric(text) {
		result = []qrSegment{makeNumeric(text)}
	} else if isAlphanumeric(text) {
		result = []qrSegment{makeAlphanumeric(text)}
	} else {
		bytes := []uint8(text)
		result = []qrSegment{makeBytes(&bytes)}
	}
	return result
}

// Returns a segment representing an Extended Channel Interpretation
// (ECI) designator with the given assignment value.
func (qrs qrSegment) makeEci(assignVal int64) qrSegment {
	bitBuf := bitBuffer{}
	if 0 <= assignVal && assignVal < (1<<7) {
		bitBuf = bitBuf.appendBits(uint32(assignVal), 8)
	} else if (1<<7) <= assignVal && assignVal < (1<<14) {
		bitBuf = bitBuf.appendBits(2, 2)
		bitBuf = bitBuf.appendBits(uint32(assignVal), 14)
	} else if (1<<14) <= assignVal && assignVal < 1000000 {
		bitBuf = bitBuf.appendBits(6, 3)
		bitBuf = bitBuf.appendBits(uint32(assignVal), 21)
	} else {
		panic(qrSegmentErr("makeEci", "ECI assignment value out of range"))
	}
	return qrSegment{isECI, 0, bitBuf}
}

// Helper function.
func getTotalBits(segs *[]qrSegment, version int) int {
	if version < 1 || version > 40 {
		panic(qrSegmentErr("getTotalBits", "version number out of range"))
	}
	result := 0
	for _, seg := range *segs {
		ccbits := seg.Mode.numCharCountBits(version)
		if uint(seg.NumChars) >= (uint(1) << uint(ccbits)) {
			return -1
		}
		if 4+ccbits > maxInt-result {
			return -1
		}
		result += 4 + ccbits
		if len(seg.Data) > (maxInt - result) {
			return -1
		}
		result += len(seg.Data)
	}
	return result
}

// Tests whether the given string can be encoded as a segment in alphanumeric mode.
func isAlphanumeric(text string) bool {
	for _, char := range text {
		if !strings.ContainsRune(alphanumericCharset, char) {
			return false
		}
	}
	return true
}

// Tests whether the given string can be encoded as a segment in numeric mode.
func isNumeric(text string) bool {
	for _, char := range text {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

// Returns qr segment's mode.
func (qrs qrSegment) getMode() modeType {
	return qrs.Mode
}

// Returns qr segment's num chars.
func (qrs qrSegment) getNumChars() int {
	return qrs.NumChars
}

// Returns an address of qr segment's data.
func (qrs qrSegment) getData() *[]bool {
	return &qrs.Data
}
