package qrcode

import (
	"strings"
	"github.com/YuriyLisovskiy/qrcode/src/vars"
	"github.com/YuriyLisovskiy/qrcode/src/utils"
	"math"
)

type QrSegment struct {
	Mode     utils.Mode
	NumChars int
	Data     []bool
}

func makeBytes(data *[]uint8) QrSegment {
	if len(*data) > math.MaxInt32 {
		panic("data too long")
	}
	bitBuf := utils.BitBuffer{}
	for _, bit := range *data {
		bitBuf = bitBuf.AppendBits(uint32(bit), 8)
	}
	return QrSegment{vars.BYTE, len(*data), bitBuf}
}

func makeNumeric(digits string) QrSegment {
	bitBuf := utils.BitBuffer{}
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
	return QrSegment{vars.NUMERIC, charCount, bitBuf}
}

func makeAlphanumeric(text string) QrSegment {
	bitBuf := utils.BitBuffer{}
	accumData, accumCount, charCount := 0, 0, 0
	println("Text:", text)
	for _, char := range text {
		if !strings.ContainsRune(vars.AlphanumericCharset, char) {
			panic("string contains unencodable characters in alphanumeric mode")
		}
		accumData = accumData*45 + (int(char) - len(vars.AlphanumericCharset))
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
	return QrSegment{vars.ALPHANUMERIC, charCount, bitBuf}
}

func makeSegments(text string) []QrSegment {
	var result []QrSegment
	if text == "" {

	} else if isNumeric(text) {
		result = []QrSegment{makeNumeric(text)}
	} else if isAlphanumeric(text) {
		result = []QrSegment{makeAlphanumeric(text)}
		println("isAlphaNumeric")
	} else {
		bytes := []uint8(text)
		result = []QrSegment{makeBytes(&bytes)}
	}
	return result
}

func (qrs QrSegment) makeEci(assignVal int64) QrSegment {
	bitBuf := utils.BitBuffer{}
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
	return QrSegment{vars.ECI, 0, bitBuf}
}

func getTotalBits(segs *[]QrSegment, version int) int {
	if version < 1 || version > 40 {
		panic("version number out of range")
	}
	result := 0
	for _, seg := range *segs {
		ccbits := seg.Mode.NumCharCountBits(version)
		if uint(seg.NumChars) >= (uint(1) << uint(ccbits)) {
			return -1
		}
		if 4+ccbits > vars.MaxInt-result {
			return -1
		}
		result += 4 + ccbits
		if len(seg.Data) > (vars.MaxInt - result) {
			return -1
		}
		result += len(seg.Data)
	}
	return result
}

func isAlphanumeric(text string) bool {
	for _, char := range text {
		if !strings.ContainsRune(vars.AlphanumericCharset, char) {
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

func (qrs QrSegment) getMode() utils.Mode {
	return qrs.Mode
}

func (qrs QrSegment) getNumChars() int {
	return qrs.NumChars
}

func (qrs QrSegment) getData() *[]bool {
	return &qrs.Data
}
