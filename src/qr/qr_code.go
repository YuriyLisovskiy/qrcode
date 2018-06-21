package qr

import "math"

type Generator struct {
	version              int
	size                 int
	errorCorrectionLevel Ecc
	mask                 int
	modules              [][]bool
	isFunction           [][]bool
}

func (qrc *Generator) EncodeText(text string) Generator {
	ecl := EccLOW
	segments := makeSegments(text)
	return qrc.encodeSegments(&segments, ecl, 1, 40, -1, true)
}

func (qrc *Generator) EncodeBinary(data *[]uint8) Generator {
	ecl := EccLOW
	return qrc.encodeSegments(&[]qrSegment{makeBytes(data)}, ecl, 1, 40, -1, true)
}

func (qrc *Generator) Draw(border, pixelSize int) {
	for y := -border; y < qrc.GetSize() + border; y++ {
		for x := -border; x < qrc.GetSize() + border; x++ {
			if qrc.GetModule(x, y) {
				print(make([]rune, pixelSize))
			} else {
				print("\u2588\u2588")
			}
		}
		println()
	}
	println()
}

func (qrc *Generator) DrawFile(path string) {

}

func newQrCode(ver int, ecl Ecc, dataCodewords []uint8, mask int) Generator {
	if ver < minVersion || ver > maxVersion || mask < -1 || mask > 7 {
		panic("value out of range")
	}
	newQrCode := Generator{version: ver, errorCorrectionLevel: ecl}
	if minVersion <= ver && ver <= maxVersion {
		newQrCode.size = ver*4 + 17
	} else {
		newQrCode.size = -1
	}
	newQrCode.modules = make([][]bool, newQrCode.size)
	for i := range newQrCode.modules {
		newQrCode.modules[i] = make([]bool, newQrCode.size)
	}
	newQrCode.isFunction = make([][]bool, newQrCode.size)
	for i := range newQrCode.isFunction {
		newQrCode.isFunction[i] = make([]bool, newQrCode.size)
	}
	newQrCode.drawFunctionPatterns()
	allCodewords := newQrCode.appendErrorCorrection(&dataCodewords)
	newQrCode.drawCodewords(&allCodewords)
	newQrCode.mask = newQrCode.handleConstructorMasking(mask)
	return newQrCode
}

func (qrc *Generator) getFormatBits(ecl Ecc) int {
	switch ecl {
	case EccLOW:
		return 1
	case EccMEDIUM:
		return 0
	case EccQUARTILE:
		return 3
	case EccHIGH:
		return 2
	default:
		panic("assertion error")
	}
}

func (qrc *Generator) encodeSegments(segs *[]qrSegment, ecl Ecc, minVersion, maxVersion, mask int, boostEcl bool) Generator {
	if !(minVersion <= minVersion && minVersion <= maxVersion && maxVersion <= maxVersion) || mask < -1 || mask > 7 {
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
	for _, newEcl := range []Ecc{EccMEDIUM, EccQUARTILE, EccHIGH} {
		if boostEcl && dataUsedBits <= qrc.getNumDataCodewords(version, newEcl)*8 {
			ecl = newEcl
		}
	}
	dataCapacityBits := uint(qrc.getNumDataCodewords(version, ecl) * 8)
	var bitBuf bitBuffer
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
	return newQrCode(version, ecl, bitBuf.getBytes(), mask)
}

func (qrc *Generator) getVersion() int {
	return qrc.version
}

func (qrc *Generator) GetSize() int {
	return qrc.size
}

func (qrc *Generator) getErrorCorrectionLevel() Ecc {
	return qrc.errorCorrectionLevel
}

func (qrc *Generator) getMask() int {
	return qrc.mask
}

func (qrc *Generator) GetModule(x, y int) bool {
	return 0 <= x && x < qrc.size && 0 <= y && y < qrc.size && qrc.module(x, y)
}

func (qrc *Generator) module(x, y int) bool {
	return qrc.modules[y][x]
}

func (qrc *Generator) getBit(x int, i uint) bool {
	return ((x >> i) & 1) != 0
}

func (qrc *Generator) drawFunctionPatterns() {
	for i := 0; i < qrc.size; i++ {
		qrc.setFunctionModule(6, i, i%2 == 0)
		qrc.setFunctionModule(i, 6, i%2 == 0)
	}
	qrc.drawFinderPattern(3, 3)
	qrc.drawFinderPattern(qrc.size-4, 3)
	qrc.drawFinderPattern(3, qrc.size-4)
	alignPatPos := []int(qrc.getAlignmentPatternPositions(qrc.version))
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

func (qrc *Generator) drawFormatBits(mask int) {
	data := int(qrc.getFormatBits(qrc.errorCorrectionLevel)<<3 | mask)
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
		qrc.setFunctionModule(qrc.size-1-i, 8, qrc.getBit(int(data), uint(i)))
	}
	for i := 8; i < 15; i++ {
		qrc.setFunctionModule(8, qrc.size-15+i, qrc.getBit(int(data), uint(i)))
	}
	qrc.setFunctionModule(8, qrc.size-8, true)
}

func (qrc *Generator) drawVersion() {
	if qrc.version < 7 {
		return
	}
	rem := int(qrc.version)
	for i := 0; i < 12; i++ {
		rem = (rem << 1) ^ ((rem >> 11) * 0x1F25)
	}
	data := int(qrc.version)<<12 | int(rem)
	if data>>18 != 0 {
		panic("assertion error")
	}
	for i := 0; i < 18; i++ {
		bit := qrc.getBit(int(data), uint(i))
		a, b := int(qrc.size-11+i%3), int(i/3)
		qrc.setFunctionModule(a, b, bit)
		qrc.setFunctionModule(b, a, bit)
	}
}

func (qrc *Generator) drawFinderPattern(x, y int) {
	for i := -4; i <= 4; i++ {
		for j := -4; j <= 4; j++ {
			dist := int(math.Max(math.Abs(float64(i)), math.Abs(float64(j))))
			xx, yy := int(x+j), int(y+i)
			if 0 <= xx && xx < qrc.size && 0 <= yy && yy < qrc.size {
				qrc.setFunctionModule(xx, yy, dist != 2 && dist != 4)
			}
		}
	}
}

func (qrc *Generator) drawAlignmentPattern(x, y int) {
	for i := -2; i <= 2; i++ {
		for j := -2; j <= 2; j++ {
			qrc.setFunctionModule(x+j, y+i, math.Max(math.Abs(float64(i)), math.Abs(float64(j))) != 1)
		}
	}
}

func (qrc *Generator) setFunctionModule(x, y int, isBlack bool) {
	qrc.modules[y][x], qrc.isFunction[y][x] = isBlack, true
}

func (qrc *Generator) appendErrorCorrection(data *[]uint8) []uint8 {
	if len(*data) != qrc.getNumDataCodewords(qrc.version, qrc.errorCorrectionLevel) {
		panic("invalid argument")
	}
	numBlocks := int(numErrorCorrectionBlocks[int(qrc.errorCorrectionLevel)][qrc.version])
	blockEccLen := int(eccCodewordsPerBlock[int(qrc.errorCorrectionLevel)][qrc.version])
	rawCodewords := int(qrc.getNumRawDataModules(qrc.version) / 8)
	numShortBlocks := int(numBlocks - rawCodewords%numBlocks)
	shortBlockLen := int(rawCodewords / numBlocks)
	var blocks [][]uint8
	rs := NewReedSolomonGenerator(blockEccLen)
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

func (qrc *Generator) drawCodewords(data *[]uint8) {
	if len(*data) != qrc.getNumRawDataModules(qrc.version)/8 {
		panic("invalid argument")
	}
	i := uint(0)
	for right := qrc.size - 1; right >= 1; right -= 2 {
		if right == 6 {
			right = 5
		}
		for vert := 0; vert < qrc.size; vert++ {
			for j := 0; j < 2; j++ {
				x := int(right - j)
				y := vert
				if ((right + 1) & 2) == 0 {
					y = qrc.size - 1 - vert
				}
				if !qrc.isFunction[y][x] && i < uint(len(*data)*8) {
					qrc.modules[y][x] = qrc.getBit(int((*data)[i>>3]), uint(7-(int(i&7))))
					i++
				}
			}
		}
	}
	if i != uint(len(*data)*8) {
		panic("assertion error")
	}
}

func (qrc *Generator) applyMask(mask int) {
	if mask < 0 || mask > 7 {
		panic("mask value out of range")
	}
	for y := 0; y < qrc.size; y++ {
		for x := 0; x < qrc.size; x++ {
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
			qrc.modules[y][x] = XOR(qrc.modules[y][x], invert && !qrc.isFunction[y][x])
		}
	}
}

func (qrc *Generator) handleConstructorMasking(mask int) int {
	if mask == -1 {
		minPenalty := maxInt	// TODO: MaxInt32
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

func (qrc *Generator) getPenaltyScore() int64 {
	result := int64(0)
	for y := 0; y < qrc.size; y++ {
		colorX := false
		runX := -1
		for x := 0; x < qrc.size; x++ {
			if x == 0 || qrc.module(x, y) != colorX {
				colorX = qrc.module(x, y)
				runX = 1
			} else {
				runX++
				if runX == 5 {
					result += penaltyN1
				} else if runX > 5 {
					result++
				}
			}
		}
	}
	for x := 0; x < qrc.size; x++ {
		colorY := false
		runY := -1
		for y := 0; y < qrc.size; y++ {
			if y == 0 || qrc.module(x, y) != colorY {
				colorY = qrc.module(x, y)
				runY = 1
			} else {
				runY++
				if runY == 5 {
					result += penaltyN1
				} else if runY > 5 {
					result++
				}
			}
		}
	}
	for y := 0; y < qrc.size-1; y++ {
		for x := 0; x < qrc.size-1; x++ {
			color := qrc.module(x, y)
			if color == qrc.module(x+1, y) && color == qrc.module(x, y+1) && color == qrc.module(x+1, y+1) {
				result += penaltyN2
			}
		}
	}
	for y := 0; y < qrc.size; y++ {
		bits := 0
		for x := 0; x < qrc.size; x++ {
			t := 0
			if qrc.module(x, y) {
				t = 1
			}
			bits = ((bits << 1) & 0x7FF) | t
			if x >= 10 && (bits == 0x05D || bits == 0x5D0) {
				result += penaltyN3
			}
		}
	}
	for x := 0; x < qrc.size; x++ {
		bits := 0
		for y := 0; y < qrc.size; y++ {
			t := 0
			if qrc.module(x, y) {
				t = 1
			}
			bits = ((bits << 1) & 0x7FF) | t
			if y >= 10 && (bits == 0x05D || bits == 0x5D0) {
				result += penaltyN3
			}
		}
	}
	black := 0
	for _, row := range qrc.modules {
		for _, color := range row {
			if color {
				black++
			}
		}
	}
	total := qrc.size * qrc.size
	for k := 0; black*int(20) < (int(9)-k)*total || black*int(20) > (int(11)+k)*total; k++ {
		result += penaltyN4
	}
	return result
}

func (qrc *Generator) getAlignmentPatternPositions(ver int) []int {
	if ver < minVersion || ver > maxVersion {
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

func (qrc *Generator) getNumRawDataModules(ver int) int {
	if ver < minVersion || ver > maxVersion {
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

func (qrc *Generator) getNumDataCodewords(ver int, ecl Ecc) int {
	if ver < minVersion || ver > maxVersion {
		panic("version number out of range")
	}
	return qrc.getNumRawDataModules(ver) / 8 - eccCodewordsPerBlock[int(ecl)][ver] * numErrorCorrectionBlocks[int(ecl)][ver]
}
