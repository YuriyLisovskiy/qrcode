package qr

import (
	"math"
	"strings"
)

type qrSegment struct {
	Mode     mode
	NumChars int
	Data     []bool
}

func makeBytes(data *[]uint8) qrSegment {
	if len(*data) > math.MaxInt32 {
		panic("data too long")
	}
	bitBuf := bitBuffer{}
	for _, bit := range *data {
		bitBuf = bitBuf.AppendBits(uint32(bit), 8)
	}
	return qrSegment{isBYTE, len(*data), bitBuf}
}

func makeNumeric(digits string) qrSegment {
	bitBuf := bitBuffer{}
	accumData, accumCount, charCount := 0, 0, 0
	for _, digit := range digits {
		if digit < '0' || digit > '9' {
			panic("string contains non-numeric characters")
		}
		accumData = accumData*10 + (int(digit) - '0')
		accumCount++
		if accumCount == 3 {
			bitBuf = bitBuf.AppendBits(uint32(accumData), 10)
			accumData, accumCount = 0, 0
		}
		charCount++
	}
	if accumCount > 0 {
		bitBuf = bitBuf.AppendBits(uint32(accumData), accumCount*3+1)
	}
	return qrSegment{isNUMERIC, charCount, bitBuf}
}

func makeAlphanumeric(text string) qrSegment {	// TODO: fix: accumData is incorrect, line 56
	bitBuf := bitBuffer{}
	accumData, accumCount, charCount := 0, 0, 0
	println("Text:", text)
	for _, char := range text {
		if !strings.ContainsRune(alphanumericCharset, char) {
			panic("string contains unencodable characters in alphanumeric mode")
		}
		accumData = accumData*45 + (int(char) - len(alphanumericCharset))
		println("accumData:", accumData)
		accumCount++
		if accumCount == 2 {
			bitBuf = bitBuf.AppendBits(uint32(accumData), 11)
			accumData, accumCount = 0, 0
		}
		charCount++
	}
	if accumCount > 0 {
		bitBuf = bitBuf.AppendBits(uint32(accumData), 6)
	}
	return qrSegment{isALPHANUMERIC, charCount, bitBuf}
}

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

func (qrs qrSegment) makeEci(assignVal int64) qrSegment {
	bitBuf := bitBuffer{}
	if 0 <= assignVal && assignVal < (1<<7) {
		bitBuf = bitBuf.AppendBits(uint32(assignVal), 8)
	} else if (1<<7) <= assignVal && assignVal < (1<<14) {
		bitBuf = bitBuf.AppendBits(2, 2)
		bitBuf = bitBuf.AppendBits(uint32(assignVal), 14)
	} else if (1<<14) <= assignVal && assignVal < 1000000 {
		bitBuf = bitBuf.AppendBits(6, 3)
		bitBuf = bitBuf.AppendBits(uint32(assignVal), 21)
	} else {
		panic("ECI assignment value out of range")
	}
	return qrSegment{isECI, 0, bitBuf}
}

func getTotalBits(segs *[]qrSegment, version int) int {
	if version < 1 || version > 40 {
		panic("version number out of range")
	}
	result := 0
	for _, seg := range *segs {
		ccbits := seg.Mode.NumCharCountBits(version)
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

func isAlphanumeric(text string) bool {
	for _, char := range text {
		if !strings.ContainsRune(alphanumericCharset, char) {
			return false
		}
	}
	return true
}

func isNumeric(text string) bool {
	for _, char := range text {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func (qrs qrSegment) getMode() mode {
	return qrs.Mode
}

func (qrs qrSegment) getNumChars() int {
	return qrs.NumChars
}

func (qrs qrSegment) getData() *[]bool {
	return &qrs.Data
}
