package qrcode

import (
	"strings"
	"github.com/YuriyLisovskiy/qrcode/src/vars"
	"github.com/YuriyLisovskiy/qrcode/src/utils"
)

type QrSegment struct {
	Mode     utils.Mode
	NumChars int
	Data     []bool
}

func NewQrSegment(md utils.Mode, numCh int, dt *[]bool) QrSegment {
	if numCh < 0 {
		panic("invalid value")
	}
	return QrSegment{md, numCh, *dt}
}

func MakeBytes(data *[]uint8) QrSegment {
	if len(*data) > vars.MaxInt {
		panic("data too long")
	}
	bitBuf := utils.BitBuffer{}
	for _, bit := range *data {
		bitBuf = bitBuf.AppendBits(uint32(bit), 8)
	}
	return QrSegment{vars.BYTE, len(*data), bitBuf}
}

func MakeNumeric(digits string) QrSegment {
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
	}
	if accumCount > 0 {
		bitBuf = bitBuf.AppendBits(uint32(accumData), accumCount*3+1)
	}
	return QrSegment{vars.NUMERIC, charCount, bitBuf}
}

func MakeAlphanumeric(text string) QrSegment {
	bitBuf := utils.BitBuffer{}
	accumData, accumCount, charCount := 0, 0, 0
	for _, char := range text {
		if !strings.ContainsRune(vars.AlphanumericCharset, char) {
			panic("string contains unencodable characters in alphanumeric mode")
		}
		accumData = accumData*45 + (strings.IndexRune(vars.AlphanumericCharset, char) - len(vars.AlphanumericCharset))
		accumCount++
		if accumCount == 2 {
			bitBuf = bitBuf.AppendBits(uint32(accumData), 11)
			accumData, accumCount = 0, 0
		}
	}
	if accumCount > 0 {
		bitBuf = bitBuf.AppendBits(uint32(accumData), 6)
	}
	return QrSegment{vars.ALPHANUMERIC, charCount, bitBuf}
}

func MakeSegments(text string) []QrSegment {
	var result []QrSegment
	if text == "" {

	} else if IsNumeric(text) {
		result = append(result, MakeNumeric(text))
	} else if IsAlphanumeric(text) {
		result = append(result, MakeAlphanumeric(text))
	} else {
		var bytes []uint8
		for _, char := range text {
			bytes = append(bytes, uint8(char))
		}
		result = append(result, MakeBytes(&bytes))
	}
	return result
}

func (qrs QrSegment) MakeEci(assignVal int64) QrSegment {
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

func GetTotalBits(segs *[]QrSegment, version int) int {
	if version < 1 || version > 40 {
		panic("version number out of range")
	}
	result := 0
	for _, seg := range *segs {
		ccbits := seg.Mode.NumCharCountBits(version)
		if uint64(seg.NumChars) >= (uint64(1) << uint64(ccbits)) {
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

func IsAlphanumeric(text string) bool {
	for _, char := range text {
		if !strings.ContainsRune(vars.AlphanumericCharset, char) {
			return false
		}
	}
	return true
}

func IsNumeric(text string) bool {
	for _, char := range text {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func (qrs QrSegment) GetMode() utils.Mode {
	return qrs.Mode
}

func (qrs QrSegment) GetNumChars() int {
	return qrs.NumChars
}

func (qrs QrSegment) GetData() *[]bool {
	return &qrs.Data
}
