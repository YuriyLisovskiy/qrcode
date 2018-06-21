package qrcode

import (
	"math"
	"github.com/YuriyLisovskiy/qrcode/src/vars"
	"github.com/YuriyLisovskiy/qrcode/src/utils"
)

type QrCode struct {
	Version              int
	Size                 int
	ErrorCorrectionLevel utils.Ecc
	Mask                 int
	Modules              [][]bool
	IsFunction           [][]bool
}

func newQrCode(ver int, ecl utils.Ecc, dataCodewords []uint8, mask int) QrCode {
	if ver < vars.MinVersion || ver > vars.MaxVersion || mask < -1 || mask > 7 {
		panic("value out of range")
	}
	newQrCode := QrCode{Version: ver, ErrorCorrectionLevel: ecl}
	if vars.MinVersion <= ver && ver <= vars.MaxVersion {
		newQrCode.Size = ver*4 + 17
	} else {
		newQrCode.Size = -1
	}
	newQrCode.Modules = make([][]bool, newQrCode.Size)
	for i := range newQrCode.Modules {
		newQrCode.Modules[i] = make([]bool, newQrCode.Size)
	}
	newQrCode.IsFunction = make([][]bool, newQrCode.Size)
	for i := range newQrCode.IsFunction {
		newQrCode.IsFunction[i] = make([]bool, newQrCode.Size)
	}
	newQrCode.drawFunctionPatterns()
	allCodewords := newQrCode.appendErrorCorrection(&dataCodewords)
	newQrCode.drawCodewords(&allCodewords)
	newQrCode.Mask = newQrCode.handleConstructorMasking(mask)
	return newQrCode
}

func (qrc *QrCode) getFormatBits(ecl utils.Ecc) int {
	switch ecl {
	case utils.EccLOW:
		return 1
	case utils.EccMEDIUM:
		return 0
	case utils.EccQUARTILE:
		return 3
	case utils.EccHIGH:
		return 2
	default:
		panic("assertion error")
	}
}

func (qrc *QrCode) EncodeText(text string, ecl utils.Ecc) QrCode {
	segments := makeSegments(text)
	return qrc.encodeSegments(&segments, ecl, 1, 40, -1, true)
}

func (qrc *QrCode) EncodeBinary(data *[]uint8, ecl utils.Ecc) QrCode {
	return qrc.encodeSegments(&[]QrSegment{makeBytes(data)}, ecl, 1, 40, -1, true)
}

func (qrc *QrCode) encodeSegments(segs *[]QrSegment, ecl utils.Ecc, minVersion, maxVersion, mask int, boostEcl bool) QrCode {
	if !(vars.MinVersion <= minVersion && minVersion <= maxVersion && maxVersion <= vars.MaxVersion) || mask < -1 || mask > 7 {
		panic("invalid value")
	}
	var version, dataUsedBits int
	for version = minVersion; ; version++ {
		dataCapacityBits := qrc.getNumDataCodewords(version, ecl) * 8
		dataUsedBits = getTotalBits(segs, version)
		if dataUsedBits != -1 && dataUsedBits <= dataCapacityBits {
			break
		}
		if version >= maxVersion {
			panic("data too long")
		}
	}
	if dataUsedBits == -1 {
		panic("assertion error")
	}
	for _, newEcl := range []utils.Ecc{utils.EccMEDIUM, utils.EccQUARTILE, utils.EccHIGH} {
		if boostEcl && dataUsedBits <= qrc.getNumDataCodewords(version, newEcl)*8 {
			ecl = newEcl
		}
	}
	dataCapacityBits := uint(qrc.getNumDataCodewords(version, ecl) * 8)
	var bitBuf utils.BitBuffer
	for _, seg := range *segs {
		bitBuf = bitBuf.AppendBits(uint32(seg.getMode().GetModeBits()), 4)
		bitBuf = bitBuf.AppendBits(uint32(seg.getNumChars()), seg.getMode().NumCharCountBits(version))
		bitBuf = append(bitBuf, *seg.getData()...)
	}
	bitBuf = bitBuf.AppendBits(0, int(math.Min(float64(4), float64(dataCapacityBits-uint(len(bitBuf))))))
	bitBuf = bitBuf.AppendBits(0, (8-len(bitBuf)%8)%8)
	for padByte := 0xEC; uint(len(bitBuf)) < dataCapacityBits; padByte ^= 0xEC ^ 0x11 {
		bitBuf = bitBuf.AppendBits(uint32(padByte), 8)
	}
	if len(bitBuf)%8 != 0 {
		panic("assertion error")
	}
	return newQrCode(version, ecl, bitBuf.GetBytes(), mask)
}

func (qrc *QrCode) getVersion() int {
	return qrc.Version
}

func (qrc *QrCode) GetSize() int {
	return qrc.Size
}

func (qrc *QrCode) getErrorCorrectionLevel() utils.Ecc {
	return qrc.ErrorCorrectionLevel
}

func (qrc *QrCode) getMask() int {
	return qrc.Mask
}

func (qrc *QrCode) GetModule(x, y int) bool {
	return 0 <= x && x < qrc.Size && 0 <= y && y < qrc.Size && qrc.module(x, y)
}

func (qrc *QrCode) module(x, y int) bool {
	return qrc.Modules[y][x]
}

func (qrc *QrCode) getBit(x int, i uint) bool {
	return ((x >> i) & 1) != 0
}

func (qrc *QrCode) drawFunctionPatterns() {
	for i := 0; i < qrc.Size; i++ {
		qrc.setFunctionModule(6, i, i%2 == 0)
		qrc.setFunctionModule(i, 6, i%2 == 0)
	}
	qrc.drawFinderPattern(3, 3)
	qrc.drawFinderPattern(qrc.Size-4, 3)
	qrc.drawFinderPattern(3, qrc.Size-4)
	alignPatPos := []int(qrc.getAlignmentPatternPositions(qrc.Version))
	numAlign := len(alignPatPos)

	for i := 0; i < numAlign; i++ {
		for j := 0; j < numAlign; j++ {
			if (i == 0 && j == 0) || (i == 0 && j == numAlign-1) || (i == numAlign-1 && j == 0) {
				continue
			} else {
				qrc.drawAlignmentPattern(alignPatPos[i], alignPatPos[j])
			}
		}
	}
	qrc.drawFormatBits(0)
	qrc.drawVersion()
}

func (qrc *QrCode) drawFormatBits(mask int) {
	data := int(qrc.getFormatBits(qrc.ErrorCorrectionLevel)<<3 | mask)
	rem := int(data)
	for i := 0; i < 10; i++ {
		rem = (rem << 1) ^ ((rem >> 9) * 0x537)
	}
	data = data<<10 | rem
	data ^= 0x5412
	if data>>15 != 0 {
		panic("assertion error")
	}
	for i := 0; i <= 5; i++ {
		qrc.setFunctionModule(8, i, qrc.getBit(int(data), uint(i)))
	}
	qrc.setFunctionModule(8, 7, qrc.getBit(int(data), 6))
	qrc.setFunctionModule(8, 8, qrc.getBit(int(data), 7))
	qrc.setFunctionModule(7, 8, qrc.getBit(int(data), 8))
	for i := 9; i < 15; i++ {
		qrc.setFunctionModule(14-i, 8, qrc.getBit(int(data), uint(i)))
	}
	for i := 0; i <= 7; i++ {
		qrc.setFunctionModule(qrc.Size-1-i, 8, qrc.getBit(int(data), uint(i)))
	}
	for i := 8; i < 15; i++ {
		qrc.setFunctionModule(8, qrc.Size-15+i, qrc.getBit(int(data), uint(i)))
	}
	qrc.setFunctionModule(8, qrc.Size-8, true)
}

func (qrc *QrCode) drawVersion() {
	if qrc.Version < 7 {
		return
	}
	rem := int(qrc.Version)
	for i := 0; i < 12; i++ {
		rem = (rem << 1) ^ ((rem >> 11) * 0x1F25)
	}
	data := int(qrc.Version)<<12 | int(rem)
	if data>>18 != 0 {
		panic("assertion error")
	}
	for i := 0; i < 18; i++ {
		bit := qrc.getBit(int(data), uint(i))
		a, b := int(qrc.Size-11+i%3), int(i/3)
		qrc.setFunctionModule(a, b, bit)
		qrc.setFunctionModule(b, a, bit)
	}
}

func (qrc *QrCode) drawFinderPattern(x, y int) {
	for i := -4; i <= 4; i++ {
		for j := -4; j <= 4; j++ {
			dist := int(math.Max(math.Abs(float64(i)), math.Abs(float64(j))))
			xx, yy := int(x+j), int(y+i)
			if 0 <= xx && xx < qrc.Size && 0 <= yy && yy < qrc.Size {
				qrc.setFunctionModule(xx, yy, dist != 2 && dist != 4)
			}
		}
	}
}

func (qrc *QrCode) drawAlignmentPattern(x, y int) {
	for i := -2; i <= 2; i++ {
		for j := -2; j <= 2; j++ {
			qrc.setFunctionModule(x+j, y+i, math.Max(math.Abs(float64(i)), math.Abs(float64(j))) != 1)
		}
	}
}

func (qrc *QrCode) setFunctionModule(x, y int, isBlack bool) {
	qrc.Modules[y][x] = isBlack
	qrc.IsFunction[y][x] = true
}

func (qrc *QrCode) appendErrorCorrection(data *[]uint8) []uint8 {
	if len(*data) != qrc.getNumDataCodewords(qrc.Version, qrc.ErrorCorrectionLevel) {
		panic("invalid argument")
	}
	numBlocks := int(vars.NumErrorCorrectionBlocks[int(qrc.ErrorCorrectionLevel)][qrc.Version])
	blockEccLen := int(vars.EccCodewordsPerBlock[int(qrc.ErrorCorrectionLevel)][qrc.Version])
	rawCodewords := int(qrc.getNumRawDataModules(qrc.Version) / 8)
	numShortBlocks := int(numBlocks - rawCodewords%numBlocks)
	shortBlockLen := int(rawCodewords / numBlocks)
	var blocks [][]uint8
	rs := utils.NewReedSolomonGenerator(blockEccLen)
	i, k := 0, 0
	for ; i < numBlocks; i++ {
		c := 1
		if i < numShortBlocks {
			c = 0
		}
		dat := []uint8((*data)[k:(k + shortBlockLen - blockEccLen + c)])
		k += len(dat)
		ecc := []uint8(rs.GetRemainder(&dat))
		if i < numShortBlocks {
			dat = append(dat, 0)
		}
		dat = append(dat, ecc...)
		blocks = append(blocks, dat)
	}
	var result []uint8
	for i := 0; i < len(blocks[0]); i++ {
		for j := 0; j < len(blocks); j++ {
			if i != shortBlockLen-blockEccLen || j >= numShortBlocks {
				result = append(result, blocks[j][i])
			}
		}
	}
	if len(result) != rawCodewords {
		panic("assertion error")
	}
	return result
}

func (qrc *QrCode) drawCodewords(data *[]uint8) {
	if len(*data) != qrc.getNumRawDataModules(qrc.Version)/8 {
		panic("invalid argument")
	}
	i := uint(0)
	for right := qrc.Size - 1; right >= 1; right -= 2 {
		if right == 6 {
			right = 5
		}
		for vert := 0; vert < qrc.Size; vert++ {
			for j := 0; j < 2; j++ {
				x := int(right - j)
				y := vert
				if ((right + 1) & 2) == 0 {
					y = qrc.Size - 1 - vert
				}
				if !qrc.IsFunction[y][x] && i < uint(len(*data)*8) {
					qrc.Modules[y][x] = qrc.getBit(int((*data)[i>>3]), uint(7-(int(i&7))))
					i++
				}
			}
		}
	}
	if i != uint(len(*data)*8) {
		panic("assertion error")
	}
}

func (qrc *QrCode) applyMask(mask int) {
	if mask < 0 || mask > 7 {
		panic("mask value out of range")
	}
	for y := 0; y < qrc.Size; y++ {
		for x := 0; x < qrc.Size; x++ {
			var invert bool
			switch mask {
			case 0:
				invert = (x+y)%2 == 0
			case 1:
				invert = y%2 == 0
			case 2:
				invert = x%3 == 0
			case 3:
				invert = (x+y)%3 == 0
			case 4:
				invert = (x/3+y/2)%2 == 0
			case 5:
				invert = x*y%2+x*y%3 == 0
			case 6:
				invert = (x*y%2+x*y%3)%2 == 0
			case 7:
				invert = ((x+y)%2+x*y%3)%2 == 0
			default:
				panic("assertion error")
			}
			qrc.Modules[y][x] = utils.XOR(qrc.Modules[y][x], invert && !qrc.IsFunction[y][x])
		}
	}
}

func (qrc *QrCode) handleConstructorMasking(mask int) int {
	if mask == -1 {
		minPenalty := vars.MaxInt	// TODO: MaxInt32
		for i := 0; i < 8; i++ {
			qrc.drawFormatBits(i)
			qrc.applyMask(i)
			penalty := int(qrc.getPenaltyScore())
			if penalty < minPenalty {
				mask = i
				minPenalty = penalty
			}
			qrc.applyMask(i)
		}
	}
	if mask < 0 || mask > 7 {
		panic("assertion error")
	}
	qrc.drawFormatBits(mask)
	qrc.applyMask(mask)
	return mask
}

func (qrc *QrCode) getPenaltyScore() int64 {
	result := int64(0)
	for y := 0; y < qrc.Size; y++ {
		colorX := false
		runX := -1
		for x := 0; x < qrc.Size; x++ {
			if x == 0 || qrc.module(x, y) != colorX {
				colorX = qrc.module(x, y)
				runX = 1
			} else {
				runX++
				if runX == 5 {
					result += vars.PENALTY_N1
				} else if runX > 5 {
					result++
				}
			}
		}
	}
	for x := 0; x < qrc.Size; x++ {
		colorY := false
		runY := -1
		for y := 0; y < qrc.Size; y++ {
			if y == 0 || qrc.module(x, y) != colorY {
				colorY = qrc.module(x, y)
				runY = 1
			} else {
				runY++
				if runY == 5 {
					result += vars.PENALTY_N1
				} else if runY > 5 {
					result++
				}
			}
		}
	}
	for y := 0; y < qrc.Size-1; y++ {
		for x := 0; x < qrc.Size-1; x++ {
			color := qrc.module(x, y)
			if color == qrc.module(x+1, y) && color == qrc.module(x, y+1) && color == qrc.module(x+1, y+1) {
				result += vars.PENALTY_N2
			}
		}
	}
	for y := 0; y < qrc.Size; y++ {
		bits := 0
		for x := 0; x < qrc.Size; x++ {
			t := 0
			if qrc.module(x, y) {
				t = 1
			}
			bits = ((bits << 1) & 0x7FF) | t
			if x >= 10 && (bits == 0x05D || bits == 0x5D0) {
				result += vars.PENALTY_N3
			}
		}
	}
	for x := 0; x < qrc.Size; x++ {
		bits := 0
		for y := 0; y < qrc.Size; y++ {
			t := 0
			if qrc.module(x, y) {
				t = 1
			}
			bits = ((bits << 1) & 0x7FF) | t
			if y >= 10 && (bits == 0x05D || bits == 0x5D0) {
				result += vars.PENALTY_N3
			}
		}
	}
	black := 0
	for _, row := range qrc.Modules {
		for _, color := range row {
			if color {
				black++
			}
		}
	}
	total := qrc.Size * qrc.Size
	for k := 0; black*int(20) < (int(9)-k)*total || black*int(20) > (int(11)+k)*total; k++ {
		result += vars.PENALTY_N4
	}
	return result
}

func (qrc *QrCode) getAlignmentPatternPositions(ver int) []int {
	if ver < vars.MinVersion || ver > vars.MaxVersion {
		panic("Version number out of range")
	} else if ver == 1 {
		return []int{}
	} else {
		numAlign := int(ver / 7 + 2)
		var step int
		if ver != 32 {
			step = (ver * 4 + numAlign * 2 + 1) / (2 * numAlign - 2) * 2
		} else {
			step = 26
		}
		var result []int
		pos := int(ver * 4 + 10)
		for i := 0; i < numAlign - 1; i++ {
			result = append([]int{pos}, result...)
			pos -= step
		}
		result = append([]int{6}, result...)
		return result
	}
}

func (qrc *QrCode) getNumRawDataModules(ver int) int {
	if ver < vars.MinVersion || ver > vars.MaxVersion {
		panic("version number out of range")
	}
	result := int((16 * ver + 128) * ver + 64)
	if ver >= 2 {
		numAlign := int(ver / 7 + 2)
		result -= (25 * numAlign - 10) * numAlign - 55
		if ver >= 7 {
			result -= 18 * 2
		}
	}
	return result
}

func (qrc *QrCode) getNumDataCodewords(ver int, ecl utils.Ecc) int {
	if ver < vars.MinVersion || ver > vars.MaxVersion {
		panic("version number out of range")
	}

//	println("Ver:", ver)
//	println("Ecc:", ecl)
//	println(qrc.getNumRawDataModules(ver))
//	println(vars.EccCodewordsPerBlock[int(ecl)][ver])
//	println(vars.NumErrorCorrectionBlocks[int(ecl)][ver])
//	println()

	return qrc.getNumRawDataModules(ver) / 8 - vars.EccCodewordsPerBlock[int(ecl)][ver] * vars.NumErrorCorrectionBlocks[int(ecl)][ver]
}
